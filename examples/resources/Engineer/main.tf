terraform {
  required_providers {
    devops-bootcamp = {
      source = "liatr.io/terraform/devops-bootcamp"
    }
  }
}

provider "devops-bootcamp" {
  endpoint = "http://localhost:8080"
}

# data "devops-bootcamp_engineer" "Ryan" {
#   name = "Ryan"
# }

resource "devops-bootcamp_engineer" "Ryan" {
  name  = "Ryan"
  email = "badEmail@hosetown.com"
}

output "devops_engineers" {
  value = resource.devops-bootcamp_engineer.Ryan
}
