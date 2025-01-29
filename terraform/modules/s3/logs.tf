# Dedicated logging bucket for CloudFront logs
resource "aws_s3_bucket" "logging_bucket" {
  bucket = "snapflow-cloudfront-logs"
}

resource "aws_s3_bucket_lifecycle_configuration" "logging_bucket_lifecycle" {
  bucket = aws_s3_bucket.logging_bucket.id

  rule {
    id     = "log_expiration"
    status = "Enabled"

    expiration {
      days = 90
    }
  }
}