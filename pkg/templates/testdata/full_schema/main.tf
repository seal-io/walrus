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
  type    = string
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
  type    = string
  default = "bar"
  description = "description of bar."
}

variable "thee" {
  type    = string
  default = "thee"
}

// number_options_var is a test variable with number options
// @options [1, 2, 3]
variable "number_options_var" {
  type    = number
  default = 1
}

// subgroup1_1 is another test variable using sub group
// @label "Subgroup1_1 Label"
// some comments
// @group "Test Subgroup/Subgroup 1"
// some other comments
variable "subgroup1_1" {
  type    = string
  default = "subgroup1_1"
}

// subgroup1_2 is another test variable using sub group
// @label "Subgroup1_2 Label"
// some comments
// @group "Test Subgroup/Subgroup 1"
// some other comments
variable "subgroup1_2" {
  type    = string
  default = "subgroup1_2"
}

// subgroup2_1 is another test variable using sub group
// @label "Subgroup2_1 Label"
// some comments
// @group "Test Subgroup/Subgroup 2"
// some other comments
variable "subgroup2_1" {
  type    = string
  default = "subgroup2_1"
}

// subgroup2_1_hidden is another test variable using sub group and should be hidden
// @hidden
// some comments
// @group "Test Subgroup/Subgroup 2"
// some other comments
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
