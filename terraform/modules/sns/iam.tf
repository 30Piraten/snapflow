# SNS IAM policy document
data "aws_iam_policy_document" "sns_policy_document" {
    statement {
      actions = [ 
        "sns:Publish"
       ]
       resources = [aws_sns_topic.snapflow_sns_topic.arn]
    }
}

# SNS IAM policy 
resource "aws_iam_policy" "sns_lambda_policy" {
  name = var.sns_lambda_policy_name
  description = var.sns_policy_description
  policy = data.aws_iam_policy_document.sns_policy_document.json 
}

# SNS-Lambda IAM policy attachment
resource "aws_iam_role_policy_attachment" "sns_policy_attachment" {
  role = var.lambda_exec_role_name
  policy_arn = aws_iam_policy.sns_lambda_policy.arn 
}