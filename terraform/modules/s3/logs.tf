# Dedicated logging bucket for CloudFront logs
resource "aws_s3_bucket" "logging_bucket" {
  bucket = var.logging_bucket_name
}

resource "aws_s3_bucket_acl" "logging_bucket_acl" {
  bucket = aws_s3_bucket.logging_bucket.id
  # acl = "log-delivery-write"
  acl = "private"
  depends_on = [ aws_s3_bucket_ownership_controls.logging_bucket_ownership ]
}

resource "aws_s3_bucket_lifecycle_configuration" "logging_bucket_lifecycle" {
  bucket = aws_s3_bucket.logging_bucket.id

  rule {
    id     = "log_expiration"
    status = "Enabled"

    expiration {
      days = 30
    }
  }
}

resource "aws_s3_bucket_ownership_controls" "logging_bucket_ownership" {
  bucket = aws_s3_bucket.logging_bucket.id

  rule {
    object_ownership = "ObjectWriter"
  }
}

