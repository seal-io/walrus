package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Tag struct {
	schema
}

func (Tag) Fields() []ent.Field {
	// keep the json tag in camel case
	return []ent.Field{
		field.Strings("tags").
			Comment("Describe tags.").
			Optional(),
	}
}
