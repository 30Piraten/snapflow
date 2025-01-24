variable "bucket_name" {
  type = string 
  default = "processedS3_bucket"
}

variable "force_destroy" {
  type = bool 
  default = false
}