terraform {
  backend "s3" {
    key = "brickrecon/webui/terraform.tfstate"
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
data "aws_caller_identity" "current" {
}

variable "environment" {
  default = "dev"
}

variable "bucket" {
  default = "brickrecon"
}
