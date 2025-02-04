variable "max_message_size" {
  type = number
}

variable "delay_seconds" {
  type = number 
}

variable "message_retention_seconds" {
  type = number
}

variable "visibility_timeout_seconds" {
  type = number
}

variable "queue_name" {
  type = string
}

variable "lambda_exec_role_name" {
  type = string 
}

variable "sqs_policy_description" {
  type = string 
}