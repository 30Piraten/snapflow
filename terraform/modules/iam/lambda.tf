locals {
  filepath = "${path.root}/frontend/config/build/dummyPrinter.zip"
}

//
resource "null_resource" "build_lambda" {
  provisioner "local-exec" {
    command = "${path.root}/frontend/config/script/script.sh"
  }
}


resource "aws_lambda_function" "dummy_print_service" {
  function_name = "dummyPrinter"
  filename = local.filepath
  handler = "main"
  runtime = "provided.al2"
  source_code_hash = filebase64sha256(local.filepath)
  role = aws_iam_role.lambda_role.arn 
}

resource "aws_iam_role" "lambda_exec_role" {
  name = "lambda-exec-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
        {
            Action = "sts:AssumeRole"
            Effect = "Allow"
            Principal = {
                Service = "lambda.amazonaws.com"
            }
        }
    ]
  })
}

resource "aws_lambda_event_source_mapping" "sqs_to_lambda" {
  batch_size = 10 
  event_source_arn = var.event_source_arn #aws_sqs_queue.print_job.arn # var.event_source_arn
  function_name = aws_lambda_function.dummy_print_service.arn 
}