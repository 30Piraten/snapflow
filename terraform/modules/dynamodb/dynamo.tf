# DynamoDB config
resource "aws_dynamodb_table" "customer_data_table" {
  name         = var.dynamodb_name
  billing_mode = var.billing_mode

  # Primary Key (partition key and sort key)
  hash_key     = "customer_email"
  range_key    = "photo_id"

  # Enable point-in-time recovery
  point_in_time_recovery {
    enabled = true
  }

  # Define attributes
  attribute {
    name = "customer_name"
    type = "S"
  }

  attribute {
    name = "customer_email"
    type = "S"
  }

  attribute {
    name = "photo_id"
    type = "S"
  }

  attribute {
    name = "photo_size"
    type = "S"
  }

  attribute {
    name = "paper_type"
    type = "S"
  }

  attribute {
    name = "processed_location"
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

  global_secondary_index {
    name = "PhotoStatusIndex"
    hash_key = "photo_status"
    projection_type = "ALL"
  }

  global_secondary_index {
    name = "CustomerEmailIndex"
    hash_key = "customer_email"
    projection_type = "ALL"
  }

  global_secondary_index {
    name = "CustomerNameIndex"
    hash_key = "customer_name"
    projection_type = "ALL"
  }

  global_secondary_index {
    name = "PhotoSizeIndex"
    hash_key = "photo_size"
    projection_type = "ALL"
  }

  global_secondary_index {
    name = "PapeType"
    hash_key = "paper_type"
    projection_type = "ALL"
  }

  global_secondary_index {
    name = "ProcessedLocationIndex"
    hash_key = "processed_location"
    projection_type = "ALL"
  }

  local_secondary_index {
    name = "UploadTimestampIndex"
    range_key = "upload_timestamp"
    projection_type = "ALL"
  }

  tags = {
    Name        = "CustomerPhotos"
    Environment = "Production"
  }
}