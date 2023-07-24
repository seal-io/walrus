package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/mixin"

	"github.com/seal-io/seal/pkg/dao/schema/io"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

func ID() id {
	return id{}
}

type id struct {
	mixin.Schema
}

func (i id) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("id").
			Immutable().
			Annotations(
				io.DisableInputWhenCreating()),
	}
}

func (id) Hooks() []ent.Hook {
	return []ent.Hook{
		object.IDHook(),
	}
}
