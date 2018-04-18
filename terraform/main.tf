# Your AWS credentials file needs to have "default" point to terraform-init account
terraform {
  backend "s3" {
    bucket     = "dr-terraform"
    key        = "remote-state/playground/dr-happy-push/terraform.tfstate"
    region     = "eu-west-1"
    role_arn   = "arn:aws:iam::596354302693:role/Administrator"
    dynamodb_table = "dr-statelock"
  }
  required_version = ">= 0.10.0"
}

provider "aws" {
  region = "${var.aws_region}"
  assume_role {
    role_arn = "${var.assume_role}"
  }
  allowed_account_ids = ["${var.aws_account_id}"]
}