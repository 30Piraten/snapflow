# S3 variable declaration
variable "bucket_name" {
  type    = string
}

variable "force_destroy" {
  type    = bool
}

variable "enable_key_rotation" {
  type = bool
}

variable "deletion_window_in_days" {
  type = number
}

variable "environment" {
  type = string 
}

variable "s3_bucket_iam_role" {
  type = string   
}

variable "s3_ksm_description" {
  type = string 
}