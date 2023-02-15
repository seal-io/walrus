package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/oid"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type Environment struct {
	schema
}

func (Environment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Meta{},
		mixin.Time{},
	}
}

func (Environment) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("connectorIDs", []oid.ID{}).
			Comment("ID of connectors of the environment.").
			Optional(),
		field.JSON("variables", map[string]interface{}{}).
			Comment("Variables of the environment.").
			Optional(),
	}
}
