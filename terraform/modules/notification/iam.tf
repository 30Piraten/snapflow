resource "aws_iam_policy" "lambda_sns_publish_policy" {
    name        = "lambda_sns_publish_policy"
    description = "Allow lambda to publish to sns"
    policy      = data.aws_iam_policy_document.lambda_sns_publish_policy.json
}

resource "aws_iam_role" "lambda_execution_role" {
    name = "lambda_execution_role"
    assume_role_policy = jsonencode({
        Version = "2012-10-17",
        Statement = [
            {
                Effect = "Allow",
                Principal = {
                    Service = "lambda.amazonaws.com"
                }
                Action = "sts:AssumeRole"
            }
        ]
    })
}

# Define the permissions for the Lambda execution role to allow
# publishing SNS messages
data "aws_iam_policy_document" "lambda_sns_publish_policy" {
    statement {
      actions = [
        "SNS:Publish",
      ]
      resources = [ aws_sns_topic.topic.arn ]
    }
}

resource "aws_iam_role_policy_attachment" "lambda_sns_publish_attachment" {
  policy_arn = aws_iam_policy.lambda_sns_publish_policy.arn
  role = aws_iam_role.lambda_execution_role.name
}
