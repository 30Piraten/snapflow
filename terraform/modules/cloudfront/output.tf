output "cloudfront_distribution_arn" {
  value = aws_cloudfront_distribution.snapflow_cloudfront.arn
}

output "cloudfront_domain_name" {
  value = aws_cloudfront_distribution.snapflow_cloudfront.domain_name
}

output "origin_access_identity" {
  value = aws_cloudfront_origin_access_identity.snapflow_oai.cloudfront_access_identity_path
}

output "origin_access_identity_arn" {
  value = aws_cloudfront_origin_access_identity.snapflow_oai.iam_arn
}

