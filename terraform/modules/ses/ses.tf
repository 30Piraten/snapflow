resource "aws_ses_email_identity" "ses_email" {
  email = var.ses_email
}