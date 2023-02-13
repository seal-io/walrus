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
			Optional().
			Comment("Status of the resource"),
		field.String("statusMessage").
			Optional().
			Comment("extra message for status, like error details"),
	}
}
