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
  name = "lambda-iam-policy"
  description = "Lambda permissions to interact with SQS, DynamoDB and SNS"
  policy = data.aws_iam_policy_document.lambda_policy_document.json 

}

data "aws_iam_policy_document" "lambda_policy_document" {
  statement {
    actions = [ 
      "sqs:SendMessage",
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes"
     ]
     resources = [ var.sqs_queue_arn ]
  }

  statement {
    actions = [ "sns:Publish" ]
    resources = [var.sns_topic_arn]
  }

  statement {
    actions = ["dynamodb:UpdateItem"]
    resources = [var.dynamodb_arn]
  }

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