package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/types/id"
)

type EnvironmentConnectorRelationship struct {
	relationSchema
}

func (EnvironmentConnectorRelationship) Indexes() []ent.Index {
	// NB(thxCode): entc cannot allow more than two fields composite as primary key through `field.ID`,
	// so we keep the default increment primary key generated from entc,
	// and use another unique key composited fields as the real primary key.
	return []ent.Index{
		// one environment can include one connector per time.
		index.Fields("environment_id", "connector_id").
			Unique(),
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
