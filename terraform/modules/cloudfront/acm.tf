resource "aws_acm_certificate" "cdn_cert" {
  domain_name       = var.domain_name
  validation_method = "DNS"

  tags = {
    Environment = "production"
    Service     = "snapflow"
    Managed_by  = "terraform"
  }
}