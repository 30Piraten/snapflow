# I am not using s3 glacier even though most of the photos
# that we edit and process at Company X might be saved long term
# Company X usually saves photos not more than 60 days. Here i am
# using the basic s3 bucket to store processed only photos, that's it.
# at Compnay X it is called sprint print (where no photos are backed up
# for a longer duration, just edit and send for printing)

# So, the use of lifecyle rule is required to delete or remove photos
# once a confirmation via SES is sent to the user

# resource "random_id" "id" {
#   byte_length = 8
# }

resource "aws_s3_bucket" "processed_image_bucket" {
  bucket = var.bucket_name

  force_destroy = var.force_destroy

  tags = {
    Name        = var.bucket_name
    Environment = "Production"
  }
}

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


resource "aws_s3_bucket_public_access_block" "processed_bucket_block" {
  bucket = aws_s3_bucket.processed_image_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_versioning" "processed_bucket_version" {
  bucket = aws_s3_bucket.processed_image_bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "processedS3_bucket_sse" {
  bucket = aws_s3_bucket.processed_image_bucket.id


  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.processed_kms_sse.id
      sse_algorithm     = "aws:kms"
    }
  }
}

# Target for CloudFront Signed URL
resource "aws_s3_bucket_policy" "processed_bucket_policy" {
  bucket = aws_s3_bucket.processed_image_bucket.id 
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
        {
            Effect = "Allow",
            Principal = {
                Service = "lambda.amazonaws.com"
            }
            Action = ["S3:GetObject", "s3:PutObject"]
            Resource = "arn:aws:s3:::${aws_s3_bucket.processed_image_bucket.id}/*"
        }
    ]
  })
}

