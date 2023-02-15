package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type Connector struct {
	schema
}

func (Connector) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Meta{},
		mixin.Status{},
		mixin.Time{},
	}
}

func (Connector) Fields() []ent.Field {
	return []ent.Field{
		field.String("driver").
			Comment("Driver type of the connector."),
		field.String("configVersion").
			Comment("Connector config version."),
		field.JSON("configData", map[string]interface{}{}).
			Comment("Connector config data.").
			Optional(),
	}
}
