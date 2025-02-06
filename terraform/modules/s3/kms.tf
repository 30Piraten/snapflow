# KMS config
resource "aws_kms_key" "processed_kms_sse" {
  description             = var.s3_ksm_description
  deletion_window_in_days = var.deletion_window_in_days
  enable_key_rotation     = var.enable_key_rotation
}