variable "region" {
  type    = string
  default = "us-east-1"
}

variable "environment" {
  type = string 
  default = "Production"
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

# SQS QUEUE DEFINITION
variable "max_message_size" {
  type    = number
  default = 262144
}

variable "delay_seconds" {
  type    = number
  default = 0
}

variable "message_retention_seconds" {
  type    = number
  default = 86400
}

variable "visibility_timeout_seconds" {
  type    = number
  default = 60
}

variable "queue_name" {
  type    = string
  default = "snapflow-photo-print-queue"
}

variable "sqs_policy_description" {
  type = string 
  default = "SQS-Lambda policy"
}

# SNS CONFIG DEFINITION
variable "sns_topic_name" {
  type = string
  default = "snapflowSNSTopic"
}

variable "sns_lambda_policy_name" {
  type = string 
  default = "snsLambdaPolicy"
}

variable "sns_policy_description" {
  type = string 
  default = "SNS-Lambda policy"
}

# SES CONFIG DEFINITION
variable "ses_email" {
  type = string 
  default = "kioyarautenberg@gmail.com"
}

variable "ses_policy_name" {
  type = string 
  default = "sesSnapflowPolicy"
}

variable "ses_policy_description" {
  type = string 
  default = "SES Policy for sending emails"
}