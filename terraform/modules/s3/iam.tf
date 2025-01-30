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
                    "arn:aws:s3:::${aws_s3_bucket.processed_image_bucket.id}/*"
                ]
            }
        ]
    })
  
}

# IAM policy for CloudFront to access S3 bucket
resource "aws_s3_bucket_policy" "snapflow_s3_policy" {
    bucket = aws_s3_bucket.processed_image_bucket.id

    policy = jsonencode({
        Version = "2012-10-17",
        Statement = [
            {
                Sid = "AllowCloudFrontReadAccess"
                Effect = "Allow"
                Principal = {
                    AWS = var.origin_access_identity_arn #aws_cloudfront_origin_access_identity.snapflow_oai.iam_arn
                }
                Action = "s3:GetObject"
                Resource = "${aws_s3_bucket.processed_image_bucket.arn}/*"
            }
        ]
    })
}

// Logging policy for CloudFront logs
resource "aws_s3_bucket_policy" "logging_bucket_policy" {
    bucket = aws_s3_bucket.logging_bucket.id

    policy = jsonencode({
        Version = "2012-10-17",
        Statement = [
            {
                Sid = "AllowCloudFrontLogging"
                Effect = "Allow"
                Principal = {
                    Service = "cloudfront.amazonaws.com"
                }
                Action = "s3:PutObject"
                Resource = "${aws_s3_bucket.logging_bucket.arn}/*"
                Condition = {
                    StringEquals = {
                        "AWS:SourceArn" = var.cloudfront_distribution_arn #aws_cloudfront_distribution.snapflow_cloudfront.arn # 
                    }
                }
            }
        ]
    })
}