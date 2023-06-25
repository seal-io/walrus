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
			Comment("Name of the resource.").
			NotEmpty(),
		field.String("description").
			Comment("Description of the resource.").
			Optional(),
		field.JSON("labels", map[string]string{}).
			Comment("Labels of the resource.").
			Default(map[string]string{}),
		field.JSON("annotations", map[string]string{}).
			Comment("Annotation of the resource.").
			Sensitive().
			Default(map[string]string{}),
	}
}
