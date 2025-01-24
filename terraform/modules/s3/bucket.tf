# I am not using s3 glacier even though most of the photos
# that we edit and process at Company X might be saved long term
# Company X usually saves photos not more than 60 days. Here i am
# using the basic s3 bucket to store processed only photos, that's it.
# at Compnay X it is called sprint print (where no photos are backed up
# for a longer duration, just edit and send for printing)


resource "aws_s3_bucket" "processedS3_bucket" {
    bucket = var.bucket_name

    force_destroy = var.force_destroy 

    tags = {
      Name = var.bucket_name
      Environment = "Production"
    }
}

resource "aws_s3_bucket_public_access_block" "processedS3_bucket_block" {
  bucket = aws_s3_bucket.processedS3_bucket.id 

  block_public_acls = true 
  block_public_policy = true
  ignore_public_acls = true 
  restrict_public_buckets = true 
}

resource "aws_s3_bucket_acl" "processedS3_bucket_acl" {
  bucket = aws_s3_bucket.processedS3_bucket.id 
  acl = "private"
}

resource "aws_s3_bucket_versioning" "processedS3_bucket_version" {
  bucket = aws_s3_bucket.processedS3_bucket.id 

  versioning_configuration {
    status = "Enabled"
  }
}

