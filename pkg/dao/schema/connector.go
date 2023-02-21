package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
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
		field.String("type").
			Comment("Type of the connector."),
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
		edge.From("environments", Environment.Type).
			Ref("connectors").
			Comment("Environments to which the connector configures.").
			Through("environmentConnectorRelationships", EnvironmentConnectorRelationship.Type),
		// connector 1-* application resources.
		edge.To("resources", ApplicationResource.Type).
			Comment("Resources that belong to the application.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
	}
}
