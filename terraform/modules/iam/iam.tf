resource "aws_iam_role" "lambda_role" {
    name = "lambda-role-service"
    assume_role_policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect = "Allow"
                Action = "sts:AssumeRole"
                Principal = {
                    Service = "lambda.amazonaws.com"
                }
            }
        ]
    })
}

resource "aws_iam_policy" "lambda_policy" {
  name = "lambda-iam-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
    {
        Action = [
            "dynamodb:UpdateItem",
            "sqs:ReceiveMessage",
        ]
        Effect = "Allow"
        Resource = "${aws_lambda_function.dummy_print_service.arn}"
    },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy_attachment" {
    policy_arn = aws_iam_policy.lambda_policy.arn 
    role = aws_iam_role.lambda_role.name 
}