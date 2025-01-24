terraform {
  required_providers {
    aws = {
      version = "~>5.0.0"
      source  = "hashicorp/aws"
    }
  }
}

provider "aws" {
  profile = "tf-user"
  region  = var.region
}