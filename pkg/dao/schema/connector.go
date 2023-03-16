package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
)

type Connector struct {
	ent.Schema
}

func (Connector) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Meta{},
		mixin.Time{},
		mixin.State{},
	}
}

func (Connector) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}

func (Connector) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").
			Comment("Type of the connector.").
			NotEmpty().
			Immutable(),
		field.String("configVersion").
			Comment("Connector config version.").
			NotEmpty(),
		field.JSON("configData", map[string]interface{}{}).
			Comment("Connector config data.").
			Default(map[string]interface{}{}).
			Sensitive(),
		field.Bool("enableFinOps").
			Comment("Config whether enable finOps, will install prometheus and opencost while enable."),
		field.JSON("finOpsCustomPricing", types.FinOpsCustomPricing{}).
			Comment("Custom pricing user defined.").
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
		// connector 1-* cluster costs.
		edge.To("clusterCosts", ClusterCost.Type).
			Comment("Cluster costs that linked to the connection").
			Annotations(
				entsql.Annotation{
					OnDelete: entsql.Cascade,
				},
			),
		// connector 1-* allocation costs.
		edge.To("allocationCosts", AllocationCost.Type).
			Comment("Cluster allocation resource costs that linked to the connection.").
			Annotations(
				entsql.Annotation{
					OnDelete: entsql.Cascade,
				},
			),
	}
}
