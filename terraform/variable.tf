variable "region" {
  type    = string
  default = "us-east-1"
}


# S3 BUCKET DEFINITION
variable "force_destroy" {
  type    = bool
  default = true
}

variable "bucket_name" {
  type    = string
  default = "snaps3flowbucket02025"
}

variable "enable_key_rotation" {
  type    = bool
  default = true
}

variable "deletion_window_in_days" {
  type    = number
  default = 7
}

# DYNAMODB BUCKET DEFINITION
variable "dynamodb_name" {
  type    = string
  default = "processedCustomerTable2025"
}

variable "billing_mode" {
  type    = string
  default = "PAY_PER_REQUEST"
}

# CLOUDFRONT DISTRIBUTION DEFINITION
variable "logging_bucket" {
  type    = string
  default = "snapflow-cloudfront-logs"  
}