terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }

    bluechi = {
      source  = "bluechi/bluechi"
      version = "1.0.0"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = "us-east-2"
}

provider "bluechi" {
  use_mock = var.use_mock
}
