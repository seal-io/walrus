variable "def_list_string" {
  default = ["a1", "a2"]
}

variable "def_list_tuple" {
  default = ["a1", true]
}

variable "def_list_object" {
  default = [{
    name = "Bob"
    email = {
      domain  = "example.com"
    }
  }]
}

