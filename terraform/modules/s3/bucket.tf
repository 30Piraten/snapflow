# I am not using s3 glacier even though most of the photos
# that we edit and process at Company X might be saved long term
# Company X usually saves photos not more than 60 days. Here i am
# using the basic s3 bucket to store processed only photos, that's it.
# at Compnay X it is called sprint print (where no photos are backed up
# for a longer duration, just edit and send for printing)

# So, the use of lifecyle rule is required to delete or remove photos
# once a confirmation via SES is sent to the user

resource "random_id" "id" {
  byte_length = 8
}

resource "aws_s3_bucket" "processedS3_bucket" {
  bucket = "${var.bucket_name}-${random_id.id.hex}"

  force_destroy = var.force_destroy

  tags = {
    Name        = var.bucket_name
    Environment = "Production"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "processedS3_bucket_lifecycle" {
  
  bucket = aws_s3_bucket.processedS3_bucket.id 

  rule {
    id = "expired-processed-photos"
    expiration {
      days = 2
    }
    status = "Enabled"
  }
}



resource "aws_s3_bucket_public_access_block" "processedS3_bucket_block" {
  bucket = aws_s3_bucket.processedS3_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_acl" "processedS3_bucket_acl" {
  bucket = aws_s3_bucket.processedS3_bucket.id
  acl    = "private"
}

resource "aws_s3_bucket_versioning" "processedS3_bucket_version" {
  bucket = aws_s3_bucket.processedS3_bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "processedS3_bucket_sse" {
  bucket = aws_s3_bucket.processedS3_bucket.id


  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.processeds3_kms_sse.id
      sse_algorithm     = "aws:kms"
    }
  }
}


# KMS
resource "aws_kms_key" "processeds3_kms_sse" {
  description             = "This key is used to encrypt the processed photos"
  deletion_window_in_days = var.deletion_window_in_days
  enable_key_rotation     = var.enable_key_rotation
}


resource "aws_iam_role" "iam_role" {
  name = "s3_iam_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
        {
            Action = "sts:AssumeRole"
            Effect = "Allow"
            Principal = {
                Service = "lambda.amazonaws.com"
            }
        }
    ]
  })

  tags = {
    Name = "s3_iam_role"
    Environment = "Production"
  }
}

resource "aws_iam_role_policy" "iam_role_policy" {
    role = aws_iam_role.iam_role.id 
    
    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Action = [
                    "s3:GetObject", 
                    "s3:PutObject",
                    "s3:ListBucket"
                ]
                Effect = "Allow"
                Resource = [
                    "arn:aws:s3:::${aws_s3_bucket.processeds3_bucket.id}/*"
                ]
            }
        ]
    })
  
}

# Target for CloudFront Signed URL
resource "aws_s3_bucket_policy" "processedS3_bucket_policy" {
  bucket = aws_s3_bucket.processedS3_bucket.id 
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
        {
            Effect = "Allow",
            Principal = {

            }
            Action = ["S3:GetObject", "s3:PutObject"]
            Resource = "arn:aws:s3:::${aws_s3_bucket.processeds3_bucket.id}/*"
        }
    ]
  })
}