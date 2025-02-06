# S3 bucket config 
resource "aws_s3_bucket" "processed_image_bucket" {
  bucket = var.bucket_name
  force_destroy = var.force_destroy
  tags = {
    Name        = var.bucket_name
    Environment = var.environment
  }
}

# Added a expiration of 7 days, since Company X 
# does not move 'Sprint-print' photos to glacier
resource "aws_s3_bucket_lifecycle_configuration" "processedS3_bucket_lifecycle" {
  bucket = aws_s3_bucket.processed_image_bucket.id 
  rule {
    id = "expired-processed-photos"
    expiration {
      days = 7
    }
    status = "Enabled"
  }
}

# S3 bucket public access block for security 
# and preventing accidental public exposure
resource "aws_s3_bucket_public_access_block" "processed_bucket_block" {
  bucket = aws_s3_bucket.processed_image_bucket.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# S3 bucket versioning
resource "aws_s3_bucket_versioning" "processed_bucket_version" {
  bucket = aws_s3_bucket.processed_image_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

# S3 server-side encryption config
# Used default KMS config here
resource "aws_s3_bucket_server_side_encryption_configuration" "processedS3_bucket_sse" {
  bucket = aws_s3_bucket.processed_image_bucket.id
  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.processed_kms_sse.id
      sse_algorithm     = "aws:kms"
    }
  }
}

