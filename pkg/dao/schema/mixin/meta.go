package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Meta struct {
	schema
}

func (Meta) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("Name of the resource."),
		field.String("description").
			Comment("Description of the resource.").
			Optional(),
		field.JSON("labels", map[string]string{}).
			Comment("Labels of the resource.").
			Optional(),
	}
}
