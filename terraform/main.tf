module "s3" {
  source = "./modules/s3"
  bucket_name = var.bucket_name
  force_destroy = var.force_destroy
  enable_key_rotation = var.enable_key_rotation
  deletion_window_in_days = var.deletion_window_in_days
  origin_access_identity_arn = module.cloudfront.origin_access_identity_arn
  cloudfront_distribution_arn = module.cloudfront.cloudfront_distribution_arn
}

module "dynamodb" {
  source = "./modules/dynamodb"
  dynamodb_name = var.dynamodb_name
  billing_mode = var.billing_mode
}

module "cloudfront" {
  source = "./modules/cloudfront"
  domain_name = module.s3.domain_name
  origin_id = module.s3.origin_id
  logging_bucket = module.s3.logging_bucket
}

module "notification" {
  source = "./modules/notification"
}