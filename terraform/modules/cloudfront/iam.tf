# Origin Access Identity (OAI) for Cloudfront
resource "aws_cloudfront_origin_access_identity" "snapflow_oai" {
  comment = "Snapflow Cloudfront OAI"
}