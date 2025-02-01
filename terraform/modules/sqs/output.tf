output "sqs_queue_url" {
  value = aws_sqs_queue.print_queue.url
}

output "sqs_queue_url_id" {
  value = aws_sqs_queue.print_queue.id 
}

output "sqs_event_source_arn" {
  value = aws_sqs_queue.print_queue.arn 
}

output "sqs_queue_arn" {
  value = aws_sqs_queue.print_queue.arn 
}