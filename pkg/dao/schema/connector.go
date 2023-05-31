package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
)

type Connector struct {
	ent.Schema
}

func (Connector) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Meta{},
		mixin.OwnByProject{}.AsOptional(),
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
		crypto.PropertiesField("configData").
			Comment("Connector config data.").
			Default(crypto.Properties{}),
		field.Bool("enableFinOps").
			Comment("Config whether enable finOps, will install prometheus and opencost while enable."),
		field.JSON("finOpsCustomPricing", &types.FinOpsCustomPricing{}).
			Comment("Custom pricing user defined.").
			Optional(),
		field.String("category").
			Comment("Category of the connector.").
			NotEmpty(),
	}
}

func (Connector) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* connectors.
		edge.From("project", Project.Type).
			Ref("connectors").
			Field("projectID").
			Comment("Project to which the connector belongs.").
			Unique().
			Immutable(),
		// Environments *-* connectors.
		edge.From("environments", Environment.Type).
			Ref("connectors").
			Comment("Environments to which the connector configures.").
			Through("environmentConnectorRelationships", EnvironmentConnectorRelationship.Type),
		// Connector 1-* service resources.
		edge.To("resources", ServiceResource.Type).
			Comment("Service resources that use the connector.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
		// Connector 1-* cluster costs.
		edge.To("clusterCosts", ClusterCost.Type).
			Comment("Cluster costs that linked to the connection").
			Annotations(
				entsql.Annotation{
					OnDelete: entsql.Cascade,
				},
			),
		// Connector 1-* allocation costs.
		edge.To("allocationCosts", AllocationCost.Type).
			Comment("Cluster allocation resource costs that linked to the connection.").
			Annotations(
				entsql.Annotation{
					OnDelete: entsql.Cascade,
				},
			),
	}
}
