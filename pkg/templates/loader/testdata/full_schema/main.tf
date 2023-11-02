# example terraform configuration

terraform {
  required_providers {
    null = {
      source = "hashicorp/null"
    }

    mycloud = {
      source  = "mycorp/mycloud"
      version = "~> 1.0"
    }
  }
}

resource "null_resource" "test" {}

variable "foo" {
  type    = string
  default = "foo"
}

variable "bar" {
  type    = string
  default = "bar"
  description = "description of bar."
}

variable "thee" {
  type    = string
  default = "thee"
}

variable "number_options_var" {
  type    = number
  default = 1
}

variable "subgroup1_1" {
  type    = string
  default = "subgroup1_1"
}

variable "subgroup1_2" {
  type    = string
  default = "subgroup1_2"
}

variable "subgroup2_1" {
  type    = string
  default = "subgroup2_1"
}

variable "subgroup2_1_hidden" {
  type    = string
  default = ""
}

output "first" {
  value       = null_resource.test.id
  description = "The first output."
}

output "second" {
  value       = "some value"
  description = "The second output."
  sensitive   = true
}
