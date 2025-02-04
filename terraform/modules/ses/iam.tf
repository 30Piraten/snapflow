data "aws_caller_identity" "account" {}

data "aws_iam_policy_document" "ses_policy_document" {
    statement {
      actions = [
        "ses:SendEmail",
      ]
      resources = [ 
        "arn:aws:ses:${var.region}:${data.aws_caller_identity.account.account_id}:identity:/${var.ses_email}",
       ]
    }
}

resource "aws_iam_policy" "ses_policy" {
  name = var.ses_policy_name 
  description = var.ses_policy_description 
  policy = data.aws_iam_policy_document.ses_policy_document.json
}

resource "aws_iam_role_policy_attachment" "ses_policy_attachment" {
  role = var.lambda_exec_role_name 
  policy_arn = aws_iam_policy.ses_policy.arn  
}