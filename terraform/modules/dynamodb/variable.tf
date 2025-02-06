# DynamoDB variable declaration
variable "dynamodb_name" {
  type = string
}

variable "billing_mode" {
  type = string
}

variable "dynamodb_iam_role_name" {
  type = string 
}

variable "dynamodb_iam_policy_name" {
  type = string 
}

variable "dynamodb_iam_policy_description" {
  type = string 
}