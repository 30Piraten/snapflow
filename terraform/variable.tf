variable "region" {
  type    = string
  default = "us-west-2"
}


# S3 BUCKET DEFINITION
variable "force_destroy" {
  type    = bool
  default = true
}

variable "bucket_name" {
  type    = string
  default = "snapflowbucket"
}

variable "enable_key_rotation" {
  type    = bool
  default = true
}

variable "deletion_window_in_days" {
  type    = number
  default = 5
}

# DYNAMODB BUCKET DEFINITION
variable "dynamodb_name" {
  type    = string
  default = "SnapFlowTable"
}

variable "billing_mode" {
  type    = string
  default = "PAY_PER_REQUEST"
}