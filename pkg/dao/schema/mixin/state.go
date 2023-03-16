package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

type State struct {
	schema
}

func (State) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("status", status.Status{}).
			Comment("Status of the object.").
			Optional(),
	}
}
