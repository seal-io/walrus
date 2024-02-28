variable "list_string" {
  type = list(string)
  default = [
    "Bob",
    "Mia"
  ]
}

variable "list_number" {
  type = list(number)
  default = [
    1,
    2
  ]
}

variable "list_bool" {
  type = list(bool)
  default = [
    true,
    false
  ]
}


variable "list_map" {
  type = list(map(string))
  default = [
    {
      name = "Bob"
    },
    {
      name = "Mia"
    }
  ]
}

variable "list_list_string" {
  type = list(list(string))
  default = [
   [
     "Bob",
     "Mia"
   ]
  ]
}

variable "list_object_with_default" {
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

variable "list_object_with_default2" {
  type = list(object({
    name = string
    age  = optional(number, 20)
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
  }]
}

variable "list_map_object" {
  type = list(map(object({
        name = string
        age  = optional(number, 20)
        email = optional(
          object({
            address = string
            domain  = optional(string, "attr.com")
          }),
          {
            address = "bob_nest",
            domain  = "nest.com"
          })
  })))
  default = [
    {
      a = {
        name = "Bob"
        email = {
          domain = "example.com"
        }
      }
    },
    {
      a = {
        name = "Mia"
      }
    }
  ]
}

variable "list_map_object_with_map" {
  type = list(map(object({
        name = string
        age  = optional(number, 20)
        email = optional(
          object({
            address = string
            domain  = optional(string, "attr.com")
          }),
          {
            address = "bob_nest",
            domain  = "nest.com"
          })
        labels = optional(
          map(string),
          {
            "job" = "teacher"
          })
  })))
  default = [
    {
      a = {
        name = "Bob"
        labels = {
          "a" = "a"
        }
      }
    },
    {
      a = {
        name = "Mia"
      }
    }
  ]
}


variable "list_map_object_with_list" {
  type = list(map(object({
        name = string
        age  = optional(number, 20)
        email = optional(
          object({
            address = string
            domain  = optional(string, "attr.com")
          }),
          {
            address = "bob_nest",
            domain  = "nest.com"
          })
        labels = optional(
          list(string),
          [
            "label1", "label1"
          ])
  })))
  default = [
    {
      a = {
        name = "Bob"
        labels = [
          "label3"
        ]
      }
    },
    {
      a = {
        name = "Mia"
      }
    }
  ]
}

variable "list_list_object" {
  type = list(list(object({
    name = string
    age  = optional(number, 20)
    email = optional(
      object({
        address = string
        domain  = optional(string, "attr.com")
      }),
      {
        address = "bob_nest",
        domain  = "nest.com"
      })
  })))
  default = [
    [
      {
        name = "Bob"
      },
      {
        name = "Mia"
      }],
    [
      {
        name = "Bob2"
      },
      {
        name = "Mia2"
      }
    ]
  ]
}

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


