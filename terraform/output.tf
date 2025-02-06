# Added outouts to confirm environment variables  
# declared in the Terraform configuration
output "s3_bucket_name" {
  value = module.s3.s3_bucket_name
}

output "dynamodb_table_name" {
  value = module.dynamodb.dynamodb_table_name
}

output "sqs_queue_url" {
  value = module.sqs.sqs_queue_url
}

output "sqs_queue_url_id" {
  value = module.sqs.sqs_queue_url_id
}

output "sns_topic_arn" {
  value = module.sns.sns_topic_arn
}