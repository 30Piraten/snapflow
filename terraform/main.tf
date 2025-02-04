module "s3" {
  source = "./modules/s3"
  bucket_name = var.bucket_name
  force_destroy = var.force_destroy
  enable_key_rotation = var.enable_key_rotation
  deletion_window_in_days = var.deletion_window_in_days
}

module "dynamodb" {
  source = "./modules/dynamodb"
  dynamodb_name = var.dynamodb_name
  billing_mode = var.billing_mode
}

module "sqs" {
  source = "./modules/sqs"
  queue_name = var.queue_name 
  message_retention_seconds = var.message_retention_seconds
  visibility_timeout_seconds = var.visibility_timeout_seconds
  max_message_size = var.max_message_size
  delay_seconds = var.delay_seconds
  lambda_exec_role_name = module.lambda.lambda_exec_role_name
}

module "lambda" {
  source = "./modules/lambda"
  event_source_arn = module.sqs.sqs_event_source_arn
  dynamodb_arn = module.dynamodb.dynamodb_arn
  sqs_queue_arn = module.sqs.sqs_queue_arn
  dynamodb_table_name = module.dynamodb.dynamodb_table_name
  sqs_queue_url = module.sqs.sqs_queue_url
  sns_topic_arn = module.notification.sns_topic_arn
}