resource "aws_dynamodb_table" "inventory_sets" {
  name = "${local.table_name}"

  write_capacity = 5 # guesses!
  read_capacity = 20 # also!

  hash_key = "setNumber"

  attribute {
    name = "setNumber"
    type = "S"
  }

  tags {
    env = "${var.environment}"
    service = "${var.product}"
  }
}
