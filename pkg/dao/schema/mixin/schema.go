package mixin

import (
	entschema "entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/mixin"
)

type schema struct {
	mixin.Schema
}

func (schema) Annotations() []entschema.Annotation {
	// Tag edges field with omitempty.
	return []entschema.Annotation{
		edge.Annotation{
			StructTag: `json:"edges,omitempty"`,
		},
	}
}
