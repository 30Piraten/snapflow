# Lambda variable declaration
variable "event_source_arn" {
  type = string 
}

variable "region" {
  type = string 
}

variable "sqs_queue_arn" {
  type = string 
}

variable "lambda_exec_role" {
  type = string 
}

variable "dynamodb_arn" {
  type = string 
}

variable "dynamodb_table_name" {
  type = string 
}

variable "sqs_queue_url" {
  type = string 
}

variable "sns_topic_arn" {
  type = string 
}

variable "ses_email" {
  type = string 
}

variable "s3_processed_image_bucket_id" {
  type = string 
}

variable "lambda_polic_name" {
  type = string 
}

variable "lambda_policy_description" {
  type = string 
}