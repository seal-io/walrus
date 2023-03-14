package mixin

import (
	"entgo.io/ent"

	"github.com/seal-io/seal/pkg/dao/types/oid"
)

type ID struct {
	schema
}

func (ID) Fields() []ent.Field {
	// keep the json tag in camel case
	return []ent.Field{
		oid.Field("id").
			Immutable(),
	}
}

func (ID) Hooks() []ent.Hook {
	return []ent.Hook{
		oid.Hook(),
	}
}
