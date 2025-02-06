# DynamoDB IAM policy for Snapflow
resource "aws_iam_role" "dynamodb_role" {
  name = var.dynamodb_iam_role_name
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

# IAM policy definition for DynamoDB
resource "aws_iam_policy" "dynamodb_policy" {
  name = var.dynamodb_iam_policy_name
  description = var.dynamodb_iam_policy_description

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

# IAM policy attachment for DynamoDB
resource "aws_iam_role_policy_attachment" "dynamodb_role_attachment" {
  role = aws_iam_role.dynamodb_role.name
  policy_arn = aws_iam_policy.dynamodb_policy.arn
}