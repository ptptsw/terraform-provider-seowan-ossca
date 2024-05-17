terraform {
  required_providers {
    seowan = {
      source = "registry.terraform.io/study/seowan-ossca"
    }
  }
}


provider "seowan" {
  username = "testuser3"
  password = "test123"
  host     = "http://localhost:19090"
}

# resource "seowan_food" "apple" {
#   items = [{
#       name = "apple",
#       price = 10000,
#   }
#   ]
# }

resource "seowan_food" "dessert" {
  items = [{
      name = "cake",
      price = 15000,
  },{
    name = "cookie",
    price = 2000,
  }
  ]
}

resource "seowan_food" "drinks" {
  items = [{
      name = "water1",
      price = 1500,
  },{
    name = "coffee2",
    price = 1000,
  }
  ]
}

# output "apple_food" {
#   value = seowan_food.apple
# }

# output "dessert_food" {
#   value = seowan_food.dessert
# }