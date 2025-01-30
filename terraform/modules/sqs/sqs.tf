resource "aws_sqs_queue" "print_queue" {
  name                      = var.queue_name
  visibility_timeout_seconds = var.visibility_timeout_seconds
  message_retention_seconds = var.message_retention_seconds
  max_message_size          = var.max_message_size
  delay_seconds             = var.delay_seconds
}
