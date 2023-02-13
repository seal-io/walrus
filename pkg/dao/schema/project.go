package schema

import (
	"entgo.io/ent"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type Project struct {
	schema
}

func (Project) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Meta{},
		mixin.Time{},
	}
}
