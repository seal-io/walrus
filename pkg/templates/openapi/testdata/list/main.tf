variable "list_object_without_default" {
  type = list(object({
    name = string
    age  = number
    email = optional(
      object({
        address = string
        domain  = optional(string, "attr.com")
      }),
      {
        address = "bob_nest",
        domain  = "nest.com"
      })
  }))
}

variable "list_object_with_default_and_nest_object" {
  type = list(object({
    name = string
    age  = number
    email = optional(
      object({
        address = string
        domain  = optional(string, "attr.com")
      }),
      {
        address = "bob_nest",
        domain  = "nest.com"
      })
  }))
  default = [{
    name = "Bob"
    age  = 23
    email = {
      address = "bob"
      domain  = "example.com"
    }
  }]
}