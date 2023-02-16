package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

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

func (Connector) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
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

func (Connector) Edges() []ent.Edge {
	return []ent.Edge{
		// environments *-* connectors.
		edge.From("environment", Environment.Type).
			Ref("connectors").
			Comment("Environments to which the connector configures.").
			Through("environmentConnectorRelationships", EnvironmentConnectorRelationship.Type),
	}
}
