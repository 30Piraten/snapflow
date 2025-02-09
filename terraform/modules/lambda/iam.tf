resource "aws_iam_role" "lambda_exec_role" {
  name = var.lambda_exec_role
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
        {
            Action = "sts:AssumeRole"
            Effect = "Allow"
            Principal = {
                Service = "lambda.amazonaws.com"
            }
        }
    ]
  })
}

resource "aws_iam_policy" "lambda_policy" {
  name = var.lambda_polic_name
  description = var.lambda_policy_description
  policy = data.aws_iam_policy_document.lambda_policy_document.json 

}

data "aws_iam_policy_document" "lambda_policy_document" {
  // SQS PERMISSIONS
  statement {
    actions = [ 
      "sqs:SendMessage",
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes"
     ]
     resources = [ var.sqs_queue_arn ]
  }

// SNS PERMISSIONS
  statement {
    actions = [ "sns:Publish" ]
    resources = [var.sns_topic_arn]
  }

  // SES PERMISSIONS
  statement {
    actions = ["ses:SendEmail"]
    resources = [ "arn:aws:ses:${var.region}:${data.aws_caller_identity.account.account_id}:identity:/${var.ses_email}", ]
  }

// DYNAMODB PERMISSIONS
  statement {
    actions = [
      "dynamodb:UpdateItem",
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:Query"
      ]
    resources = [var.dynamodb_arn]
  }

// S3 BUCKET PERMISSIONS
  statement {
    actions = [ 
      "s3:PutObject",
      "s3:GetObject",
      "s3:ListBucket"
     ]
    #  resources = ["arn:aws:s3:::${aws_s3_bucket.processed_image_bucket.id}/*"]
     resources = ["arn:aws:s3:::${var.s3_processed_image_bucket_id}/*"]
  }

  statement {
    actions = [ "s3:ListObject" ]
    resources = ["arn:aws:s3:::${var.s3_processed_image_bucket_id}"]
  }

//CLOUDWATCH PERMISSIONS
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${var.region}:${data.aws_caller_identity.account.account_id}:log-group:/aws/lambda/${aws_lambda_function.dummy_print_service.function_name}:*"
      ]
  }
}


resource "aws_iam_role_policy_attachment" "lambda_policy_attachment" {
    role = aws_iam_role.lambda_exec_role.name  
    policy_arn = aws_iam_policy.lambda_policy.arn 
}