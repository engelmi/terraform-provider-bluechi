terraform {
  required_providers {
    bluechi = {
      source = "bluechi/bluechi"
    }
  }
}

provider "bluechi" {
  use_mock = var.use_mock
}
