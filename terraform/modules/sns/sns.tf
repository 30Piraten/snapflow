resource "aws_sns_topic" "snapflow_sns_topic" {
  name = var.sns_topic_name 
}

resource "aws_sns_topic_subscription" "sns_ses_subscription" {
  topic_arn = aws_sns_topic.snapflow_sns_topic.arn 
  protocol = "email-json"
  endpoint = var.ses_email_identity
}