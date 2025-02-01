# resource "null_resource" "build_lambda" {
#   provisioner "local-exec" {
#     command = "${path.module}/script/script.sh"
#   }
# }

resource "aws_lambda_function" "dummy_print_service" {
  function_name = "dummyPrinter"
  filename = "../frontend/lambda/dummyprinter.zip"
  # filename = "../frontend/test/dummyprinter.zip"
  handler = "bootstrap"
  runtime = "provided.al2"
  role = aws_iam_role.lambda_exec_role.arn 
  # timeout = 60
  memory_size = 128

  # depends_on = [ null_resource.build_lambda ]

  environment {
    variables = {
      SQS_QUEUE_URL = var.sqs_queue_url
      DYNAMODB_TABLE = var.dynamodb_table_name 
    }
  }
}

resource "aws_lambda_event_source_mapping" "sqs_to_lambda" {
  batch_size = 10 
  event_source_arn = var.event_source_arn #aws_sqs_queue.print_job.arn # var.event_source_arn
  function_name = aws_lambda_function.dummy_print_service.arn 
  enabled = true 
}