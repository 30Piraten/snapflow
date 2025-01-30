output "s3_bucket_name" {
  value = module.s3.s3_bucket_name
}

output "dynamodb_table_name" {
  value = module.dynamodb.dynamodb_table_name
}

output "sns_topic_arn" {
  value = module.notification.sns_topic_arn
}

output "sqs_queue_url" {
  value = module.sqs.sqs_queue_url
}