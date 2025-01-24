# KMS
resource "aws_kms_key" "processeds3_kms_sse" {
  description             = "This key is used to encrypt the processed photos"
  deletion_window_in_days = var.deletion_window_in_days
  enable_key_rotation     = var.enable_key_rotation
}