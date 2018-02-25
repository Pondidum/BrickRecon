terraform {
  backend "s3" {
    key = "brickrecon/inventory/terraform.tfstate"
    region = "eu-west-1"
  }
}

variable "region" {
  default = "eu-west-1" # Irish region best region!
}

provider "aws" {
  profile = "default"
  region = "${var.region}"
}

// this will fetch our account_id, no need to hard code it
data "aws_caller_identity" "current" {}

variable "brickowl_token" {}
variable "environment" {}
variable "product" {
  default = "brickrecon"
}

locals {
  prefix = "${var.product != "" ? "${var.product}_" : ""}${var.environment}_"
  name = "${local.prefix}inventory"
  sets_table = "${local.name}_sets"
  boids_table = "${local.name}_boids"
  sns_topic = "${local.prefix}kafish_events"
  kafish = "${local.prefix}kafish_writer"
}
