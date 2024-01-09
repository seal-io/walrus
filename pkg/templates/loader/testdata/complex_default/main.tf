variable "object_with_default" {
  type = object({
    name = string
    age  = number
  })
  default = {
    name = "Bob"
    age  = 23
  }
}

variable "object_with_attr_default" {
  type = object({
    name = optional(string, "Bob")
    age  = optional(number, 23)
  })
}

variable "object_with_nest_object" {
  type = object({
    name = string
    age  = number
    email = optional(object({
      address = string
      domain  = optional(string, "attr.com")
    }))
  })
}

variable "object_with_default_and_nest_object" {
  type = object({
    name = string
    age  = number
    email = optional(object({
      address = string
      domain  = optional(string, "attr.com")
    }))
  })
  default = {
    name = "Bob"
    age  = 23
    email = {
      address = "bob"
      domain  = "example.com"
    }
  }
}

variable "object_with_default_and_nest_object2" {
  type = object({
    name = string
    age  = number
    email = optional(object({
      address = string
      domain  = optional(string, "attr.com")
    }))
  })
  default = {
    name = "Bob"
    age  = 23
    email = {
      address = "bob"
    }
  }
}

variable "object_with_default_and_nest_object3" {
  type = object({
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
  })
  default = {
    name = "Bob"
    age  = 23
    email = {
      address = "bob"
      domain  = "example.com"
    }
  }
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
