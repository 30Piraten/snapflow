# Terraform provider configuration
terraform {
  required_providers {
    aws = {
      version = "~>5.84.0"
      source  = "hashicorp/aws"
    }
  }
}

# Using my AWS configured profile
provider "aws" {
  profile = "tf-user"
  region  = var.region
}