# DynamoDB config

resource "aws_dynamodb_table" "processed_data_table" {
  name         = var.dynamodb_name
  billing_mode = var.billing_mode
  hash_key     = "UserId"
  range_key    = "SnapFlow"

  point_in_time_recovery {
    enabled = true
  }

  attribute {
    name = "UserId"
    type = "S"
  }

  attribute {
    name = "processed_s3_location"
    type = "S"
  }

  attribute {
    name = "photo_status"
    type = "S"
  }

  attribute {
    name = "upload_timestamp"
    type = "N"
  }

  tags = {
    Name        = "SnapFlowTable"
    Environment = "Production"
  }
}