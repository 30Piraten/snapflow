output "dynamodb_table_name" {
  value = aws_dynamodb_table.customer_data_table.name
}

output "dynamodb_arn" {
  value = aws_dynamodb_table.customer_data_table.arn 
}