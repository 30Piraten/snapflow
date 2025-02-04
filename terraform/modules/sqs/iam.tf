data "aws_iam_policy_document" "sqs_role_policy" {
  statement {
    actions = [ 
        "sqs:SendMessage",
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes"
     ]

     resources = [ aws_sqs_queue.print_queue.arn ]
  }
}

resource "aws_iam_policy" "lambda_sqs_queue" {
  name = "lambda-sqs-policy"
   description = var.sqs_policy_description
   policy = data.aws_iam_policy_document.sqs_role_policy.json 
}

resource "aws_iam_role_policy_attachment" "lambda_sqs_attachment" {
  policy_arn = aws_iam_policy.lambda_sqs_queue.arn 
  role = var.lambda_exec_role_name
}