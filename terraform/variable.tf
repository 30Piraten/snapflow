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

# SQS QUEUE DEFINITION
variable "max_message_size" {
  type    = number
  default = 262144
}

variable "delay_seconds" {
  type    = number
  default = 5
}

variable "message_retention_seconds" {
  type    = number
  default = 86400
}

variable "visibility_timeout_seconds" {
  type    = number
  default = 30
}

variable "queue_name" {
  type    = string
  default = "snapflow-photo-print-queue"
}
