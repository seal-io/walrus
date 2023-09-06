package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"github.com/seal-io/walrus/pkg/dao/entx"
	daostatus "github.com/seal-io/walrus/pkg/dao/types/status"
)

func Status() status {
	return status{}
}

type status struct {
	mixin.Schema
}

func (i status) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("status", daostatus.Status{}).
			Optional().
			Annotations(
				entx.SkipInput()),
	}
}
