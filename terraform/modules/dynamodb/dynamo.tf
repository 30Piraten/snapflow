# DynamoDB config

resource "random_id" "id" {
    byte_length = 8
}

resource "aws_dynamodb_table" "customer_data_table" {
  name         = "${var.dynamodb_name}-${random_id.id.hex}"
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

  global_secondary_index {
    name = "PhotoStatusIndex"
    hash_key = "photo_status"
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