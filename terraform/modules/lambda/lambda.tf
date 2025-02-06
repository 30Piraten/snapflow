# Lambda function definition
resource "aws_lambda_function" "dummy_print_service" {
  function_name = "dummyPrinter"
  filename = "../src/lambda/dummyprinter.zip"
  handler = "bootstrap"
  runtime = "provided.al2"
  role = aws_iam_role.lambda_exec_role.arn 
  timeout = 60
  memory_size = 128

  environment {
    variables = {
      SQS_QUEUE_URL = var.sqs_queue_url
      SNS_TOPIC_ARN = var.sns_topic_arn
      DYNAMODB_TABLE_NAME = var.dynamodb_table_name 
      AWS_LAMBDA_EXEC_WRAPPER = "/opt/bootstrap"
    }
  }
}

# AWS caller identity to retrive account ID
data "aws_caller_identity" "account" {}

# Lambda event source mapping for SQS-Lambda
resource "aws_lambda_event_source_mapping" "sqs_to_lambda" {
  batch_size = 10 
  event_source_arn = var.event_source_arn 
  function_name = aws_lambda_function.dummy_print_service.arn 
  enabled = true 
}