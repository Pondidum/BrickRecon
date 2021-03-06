resource "aws_dynamodb_table" "inventory_sets" {
  name = "${local.sets_table}"

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

resource "aws_dynamodb_table" "boid_cache" {
  name = "${local.boids_table}"

  write_capacity = 5 # guesses!
  read_capacity = 20 # also!

  hash_key = "boid"

  attribute {
    name = "boid"
    type = "S"
  }

  tags {
    env = "${var.environment}"
    service = "${var.product}"
  }
}
