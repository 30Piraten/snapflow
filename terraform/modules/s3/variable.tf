variable "bucket_name" {
  type    = string
  default = "processedS3_bucket"
}

variable "force_destroy" {
  type    = bool
  default = false
}

variable "enable_key_rotation" {
  type = true
}

variable "deletion_window_in_days" {
  type = number
}