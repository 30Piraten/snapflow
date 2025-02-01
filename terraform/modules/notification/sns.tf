resource "aws_sns_topic" "topic" {
  name = "snapflow"
}

resource "aws_sns_topic_policy" "sns" {
  arn = aws_sns_topic.topic.arn

  policy = data.aws_iam_policy_document.sns_policy.json
}
  
data "aws_iam_policy_document" "sns_policy" {
  statement {
    actions = [
      "SNS:Publish",
    ]
    principals {
      type        = "AWS"
      identifiers = [var.lambda_execution_role_arn]
    }
    resources = [aws_sns_topic.topic.arn]
  }
}