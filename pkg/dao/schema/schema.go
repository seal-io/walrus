package schema

import (
	"entgo.io/ent"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type schema struct {
	ent.Schema
}

func (schema) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Time{},
	}
}
