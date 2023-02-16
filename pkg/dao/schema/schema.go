package schema

import (
	"entgo.io/ent"
	ents "entgo.io/ent/schema"

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

type relationSchema struct {
	ent.Schema
}

func (relationSchema) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.CreateTime{},
	}
}

type Annotation = ents.Annotation
