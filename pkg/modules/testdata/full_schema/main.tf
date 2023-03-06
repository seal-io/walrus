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

# foo is a test variable
#
#@label "Foo Label"
# @options ["F1","F2","F3"]
#   @group        "Test Group"
variable "foo" {
  type = string
  default = "foo"
}

// bar is another test variable using a different comment style
//@label "Bar Label"
// some comments
// @options ["B1","B2","B3"]
//   @group        "Test Group"
// Using tab characters for indents as the following also works
// 	@show_if	"foo=F1"
// some other comments
variable "bar" {
  type = string
  default = "bar"
}

variable "thee" {
  type = string
  default = "thee"
}

output "first" {
  value = null_resource.test.id
  description = "The first output."
}

output "second" {
  value = "some value"
  description = "The second output."
  sensitive = true
}
