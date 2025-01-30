resource "aws_iam_role" "lambda_role" {
    name = ""
    assume_role_policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect = "Allow"
                Action = "sts:AssumeRole"
                Principal = {
                    Service = "lambda.amazonaws.com"
                }
                Resource = "*" // lambda arn 
            }
        ]
    })
}

resource "aws_iam_policy" "lambda_policy" {
  name = ""

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
    {
        Action = [
            "dynamodb:UpdateItem",
            "sqs:ReceiveMessage",
        ]
        Effect = "Allow"
        Resource = ""
    },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy_attachment" {
    policy_arn = aws_iam_policy.lambda_policy.arn 
    role = aws_iam_role.lambda_role.name 
}