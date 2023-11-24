# example terraform configuration

variable "any" {
}

variable "any_map" {
  type    = map(any)
  default = null
}

variable "string_map" {
  type    = map(string)
  default = {
    a = "a"
    b = "1"
    c = "true"
  }
}

variable "string_slice" {
  type    = list(string)
  default = [
    "x", "y", "z"
  ]
}

variable "object" {
  type = object({
    a = string
    b = number
    c = bool
  })
  default = {
    a = "a"
    b = 1
    c = true
  }
}

variable "object_nested" {
  type = object({
    a = string
    b = list(object({
      c = bool
    }))
  })
  default = {
    a = "a"
    b = [
      {
        c = true
      }
    ]
  }
}

variable "list_object" {
  type = list(object({
    a = string
    b = number
    c = bool
  }))
}

variable "tuple" {
  type = tuple([string, bool, number])
}

variable "object_tuple" {
  type = object({
    data = optional(tuple([string, bool]), ["foo", true])
  })
}
