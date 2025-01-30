output "s3_bucket_name" {
 value = aws_s3_bucket.processed_image_bucket.bucket 
}

output "processed_image_bucket_id" {
 value = aws_s3_bucket.processed_image_bucket.id
}

output "domain_name" {
 value = aws_s3_bucket.processed_image_bucket.bucket_regional_domain_name
}

output "origin_id" {
    value = aws_s3_bucket.processed_image_bucket.bucket_regional_domain_name
}

# output "logging_bucket" {
#   value = aws_s3_bucket.logging_bucket.bucket 
# }