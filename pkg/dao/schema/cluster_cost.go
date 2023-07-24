package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

// ClusterCost holds the schema definition for the cluster hourly cost.
type ClusterCost struct {
	ent.Schema
}

func (ClusterCost) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("start_time", "end_time", "connector_id").
			Unique(),
	}
}

func (ClusterCost) Fields() []ent.Field {
	return []ent.Field{
		field.Time("start_time").
			Comment("Usage start time for current cost.").
			Immutable(),
		field.Time("end_time").
			Comment("Usage end time for current cost.").
			Immutable(),
		field.Float("minutes").
			Comment("Usage minutes from start time to end time.").
			Immutable(),
		object.IDField("connector_id").
			Comment("ID of the connector.").
			NotEmpty().
			Immutable(),
		field.String("cluster_name").
			Comment("Cluster name for current cost.").
			NotEmpty().
			Immutable(),
		field.Float("total_cost").
			Comment("Cost number.").
			Default(0).
			Min(0),
		field.Int("currency").
			Comment("Cost currency.").
			Optional(),
		field.Float("allocation_cost").
			Comment("Allocation cost for current cost.").
			Default(0).
			Min(0),
		field.Float("idle_cost").
			Comment("Idle cost for current cost.").
			Default(0).
			Min(0),
		field.Float("management_cost").
			Comment("Storage cost for current cost.").
			Default(0).
			Min(0),
	}
}

func (ClusterCost) Edges() []ent.Edge {
	return []ent.Edge{
		// Connector 1-* ClusterCosts.
		edge.From("connector", Connector.Type).
			Comment("Connector current cost linked.").
			Ref("cluster_costs").
			Field("connector_id").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.SkipIO()),
	}
}

func (ClusterCost) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.SkipClearingOptionalField(),
	}
}
