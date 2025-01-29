variable "domain_name" {
  description = "The domain name for the CloudFront distribution"
  type        = string
}

variable "origin_id" {
  type = string
}

variable "logging_bucket" {
  type = string 
}
