# example terraform configuration

variable "list_any_with_default" {
  type        = list(any)
  default = [
    {
      name = "default-name"
    }
  ]
}

variable "map_any_with_default" {
  type        = map(any)

  default = {
    name = "default-name"
  }
}

variable "list_map_any_with_default" {
  type        = list(map(any))

  default = [
    {
      name = "default-name"
    }
  ]
}

variable "object_with_any_default" {
  type = object({
    any_data = optional(any, {
      port    = 80
      headers = {
        "X-Forwarded-Proto" = "https"
      }
    })
  })
}

variable "list_object_with_any_default" {
  type = list(object({
    any_data = optional(any, {
      port    = 80
      headers = {
        "X-Forwarded-Proto" = "https"
      }
    })
  }))
}

variable "map_object_with_any_default" {
  type = map(object({
    any_data = optional(any, {
      port    = 80
      headers = {
        "X-Forwarded-Proto" = "https"
      }
    })
  }))
}
