# S3 bucket IAM role
resource "aws_iam_role" "iam_role" {
  name = var.s3_bucket_iam_role
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
    Name = var.s3_bucket_iam_role
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
                    "s3:ListBucket",
                    "s3:ListObject"
                ]
                Effect = "Allow"
                Resource = [
                    "arn:aws:s3:::${aws_s3_bucket.processed_image_bucket.id}/*"
                ]
            }
        ]
    })
}