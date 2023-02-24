package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Status struct {
	schema
}

func (Status) Fields() []ent.Field {
	return []ent.Field{
		field.String("status").
			Comment("Status of the resource.").
			Optional(),
		field.String("statusMessage").
			Comment("Extra message for status, like error details.").
			Optional(),
	}
}
