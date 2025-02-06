# SES variable declaration
variable "region" {
  type = string 
}

variable "ses_email" {
  type = string 
}

variable "ses_policy_name" {
  type = string 
}

variable "ses_policy_description" {
  type = string 
}

variable "lambda_exec_role_name" {
  type = string 
}