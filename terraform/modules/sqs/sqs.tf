resource "aws_sqs_queue" "print_queue" {
  name                      = "photo-print-queue"
  visibility_timeout_seconds = 60
  message_retention_seconds = 86400
  max_message_size          = 262144
  delay_seconds             = 0
}
