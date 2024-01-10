variable "map_object_with_default_and_nest_object" {
  type = map(object({
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
  default = {
    "ab" : {
      name = "Bob"
      age  = 23
      email = {
        address = "bob"
        domain  = "example.com"
      }
    }
  }
}