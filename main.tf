terraform {
  required_providers {
    dob = {
      source = "liatr.io/terraform/devops-bootcamp"
    }
  }
}

provider "dob" {
  endpoint = "http://localhost:8080"
}

data "dob_engineer" "all" {}

output "engineers" {
  value = [for e in data.dob_engineer.all.engineers : e.id]
}

data "dob_dev" "all" {}

output "devs" {
  value = [for d in data.dob_dev.all.devs : d.id]
}

data "dob_ops" "all" {}

output "ops" {
  value = [for o in data.dob_ops.all.ops : o.id]
}

data "dob_devops" "all" {}

output "devops" {
  value = [for d in data.dob_devops.all.devops : d.id]
}

resource "dob_engineer" "Madi" {
  name = "Madi"
  email = "madi@liatrio.com"
}

resource "dob_engineer" "Colin" {
  name = "Colin"
  email = "colin@liatrio.com"
}

resource "dob_engineer" "Angel" {
  name = "Angel"
  email = "angel@liatrio.com"
}

resource "dob_engineer" "Austin" {
  name = "Austin"
  email = "austin@liatrio.com"
}

resource "dob_engineer" "Jack" {
  name = "Jack"
  email = "jack@liatrio.com"
}

resource "dob_dev" "example" {
  name = "Dev Team #1"
  engineers = [
    dob_engineer.Angel.id, 
    dob_engineer.Colin.id
  ]
}

resource "dob_ops" "example" {
  name = "Ops Team #1"
  engineers = [
    dob_engineer.Madi.id, 
    dob_engineer.Austin.id, 
    dob_engineer.Jack.id
  ]
}
  
resource "dob_devops" "example" {
  devs = [
    dob_dev.example.id
  ]
  ops = [
    dob_ops.example.id
  ]
}
