# Config definition for S3 bucket
module "s3" {
  source = "./modules/s3"
  bucket_name = var.bucket_name
  environment = var.environment
  force_destroy = var.force_destroy
  s3_ksm_description = var.s3_ksm_description
  s3_bucket_iam_role = var.s3_bucket_iam_role
  enable_key_rotation = var.enable_key_rotation
  deletion_window_in_days = var.deletion_window_in_days
}

# Config definition for DynamoDB
module "dynamodb" {
  source = "./modules/dynamodb"
  dynamodb_name = var.dynamodb_name
  billing_mode = var.billing_mode
  dynamodb_iam_role_name = var.dynamodb_iam_role_name
  dynamodb_iam_policy_name = var.dynamodb_iam_policy_name
  dynamodb_iam_policy_description = var.dynamodb_iam_policy_description
}

# Config definition fot SQS 
module "sqs" {
  source = "./modules/sqs"
  queue_name = var.queue_name 
  delay_seconds = var.delay_seconds
  max_message_size = var.max_message_size
  sqs_lambda_policy_name = var.sqs_lambda_policy_name
  sqs_policy_description = var.sqs_policy_description
  message_retention_seconds = var.message_retention_seconds
  # visibility_timeout_seconds = var.visibility_timeout_seconds 
  lambda_exec_role_name = module.lambda.lambda_exec_role_name
}

# Config definition for SNS
module "sns" {
  source = "./modules/sns"
  sns_topic_name = var.sns_topic_name
  sns_email_protocol = var.sns_email_protocol
  sns_policy_description = var.sns_policy_description
  sns_lambda_policy_name = var.sns_lambda_policy_name
  ses_email_identity = module.ses.ses_email_identity
  lambda_exec_role_name = module.lambda.lambda_exec_role_name
}

# Config defintion for SES
module "ses" {
  source = "./modules/ses"
  region = var.region
  ses_email = var.ses_email
  ses_policy_name = var.ses_policy_name
  ses_policy_description = var.ses_policy_description
  lambda_exec_role_name = module.lambda.lambda_exec_role_name
}

# Config defintion for Lambda
module "lambda" {
  source = "./modules/lambda"
  region = var.region
  ses_email = var.ses_email
  lambda_polic_name = var.lambda_polic_name
  lambda_policy_description = var.lambda_policy_description
  lambda_exec_role = var.lambda_exec_role
  sqs_queue_arn = module.sqs.sqs_queue_arn
  sqs_queue_url = module.sqs.sqs_queue_url
  sns_topic_arn = module.sns.sns_topic_arn
  dynamodb_arn = module.dynamodb.dynamodb_arn
  event_source_arn = module.sqs.sqs_event_source_arn
  dynamodb_table_name = module.dynamodb.dynamodb_table_name
  s3_processed_image_bucket_id = module.s3.processed_image_bucket_id
}