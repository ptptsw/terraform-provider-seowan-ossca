terraform {
  required_providers {
    seowan = {
      source = "registry.terraform.io/study/seowan-ossca"
    }
  }
}

provider "seowan" {
  username = "testuser"
  password = "test123"
  host     = "http://localhost:19090"
}

resource "seowan_food" "apple" {
  items = [{
      name = "apple",
      price = 10000,
  }
  ]
}

output "apple_food" {
  value = seowan_food.apple
}
