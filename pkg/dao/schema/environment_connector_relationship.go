package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	ents "entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/id"
)

type EnvironmentConnectorRelationship struct {
	ent.Schema
}

func (EnvironmentConnectorRelationship) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.CreateTime{},
	}
}

func (EnvironmentConnectorRelationship) Annotations() []ents.Annotation {
	return []ents.Annotation{
		field.ID("environment_id", "connector_id"),
	}
}

func (EnvironmentConnectorRelationship) Fields() []ent.Field {
	return []ent.Field{
		id.Field("environment_id").
			Comment("ID of the environment to which the relationship connects.").
			StructTag(`json:"environmentID"`).
			NotEmpty().
			Immutable(),
		id.Field("connector_id").
			Comment("ID of the connector to which the relationship connects.").
			StructTag(`json:"connectorID"`).
			NotEmpty().
			Immutable(),
	}
}

func (EnvironmentConnectorRelationship) Edges() []ent.Edge {
	// NB(thxCode): entc cannot recognize camel case field name on edge with `Through`.
	return []ent.Edge{
		edge.To("environment", Environment.Type).
			Field("environment_id").
			Comment("Environments that connect to the relationship.").
			Unique().
			Required().
			Immutable().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("connector", Connector.Type).
			Field("connector_id").
			Comment("Connectors that connect to the relationship.").
			Unique().
			Required().
			Immutable().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
	}
}
