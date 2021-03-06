resource "aws_dynamodb_table" "event_store" {
  name = "${local.table_name}"

  write_capacity = 5 # guesses!
  read_capacity = 20 # also!

  hash_key = "eventId"
  range_key = "timestamp"

  global_secondary_index = {
    name = "eventType"
    hash_key = "eventType"
    range_key = "timestamp"
    projection_type = "ALL"
    write_capacity = 5 # guesses!
    read_capacity = 10 # also!
  }

  attribute {
    name = "eventId"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "N"
  }

  attribute {
    name = "eventType"
    type = "S"
  }

  tags {
    env = "${var.environment}"
    service = "${var.product}"
  }
}
