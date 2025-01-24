module "s3" {
  source = "./modules/s3"
  bucket_name = var.bucket_name
  force_destroy = var.force_destroy
  enable_key_rotation = var.enable_key_rotation
  deletion_window_in_days = var.deletion_window_in_days
}

# module "dynamodb" {
#   source = "./modules/dynamodb"
#   dynamodb_name = var.dynamodb_name
#   billing_mode = var.billing_mode
# }