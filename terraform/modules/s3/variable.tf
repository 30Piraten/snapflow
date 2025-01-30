variable "bucket_name" {
  type    = string
}

variable "logging_bucket_name" {
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

variable "cloudfront_distribution_arn" {
  type = string
}

variable "origin_access_identity_arn" {
  type = string
}