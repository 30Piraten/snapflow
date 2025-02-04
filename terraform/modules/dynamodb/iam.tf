# DynamoDB IAM policy for Snapflow
resource "aws_iam_role" "dynamodb_role" {
  name = "snapflow-dynamodb-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
  
}

resource "aws_iam_policy" "dynamodb_policy" {
  name = "snapflow-dynamodb-policy"
  description = "IAM policy for DynamoDB permissions"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:Query",
          "dynamodb:UpdateItem"
        ]
        Resource = "${aws_dynamodb_table.customer_data_table.arn}"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "dynamodb_role_attachment" {
  role = aws_iam_role.dynamodb_role.name
  policy_arn = aws_iam_policy.dynamodb_policy.arn
}