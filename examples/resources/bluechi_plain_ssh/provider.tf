terraform {
  required_providers {
    bluechi = {
      source  = "bluechi/bluechi"
      version = "1.0.0"
    }
  }
}

provider "bluechi" {
  use_mock = var.use_mock
}
