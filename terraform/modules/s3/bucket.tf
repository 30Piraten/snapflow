resource "aws_s3_bucket" "originalS3_bucket" {
    bucket = "originalS3_bucket"
    
}

resource "aws_s3_bucket" "processedS3_bucket" {
    bucket = "processedS3_bucket"
}