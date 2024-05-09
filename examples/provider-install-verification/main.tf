terraform {
  required_providers {
    seowan = {
      source = "registry.terraform.io/study/seowan-ossca"
    }
  }
}

provider "seowan" {}

# data "hashicups_coffees" "example" {}
