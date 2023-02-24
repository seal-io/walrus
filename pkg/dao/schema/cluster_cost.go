package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/types/id"
)

// ClusterCost holds the schema definition for the cluster hourly cost.
type ClusterCost struct {
	ent.Schema
}

func (ClusterCost) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("startTime", "endTime", "connectorID").
			Unique(),
	}
}

func (ClusterCost) Fields() []ent.Field {
	return []ent.Field{
		field.Time("startTime").
			Comment("Usage start time for current cost").
			Immutable(),
		field.Time("endTime").
			Comment("Usage end time for current cost").
			Immutable(),
		field.Float("minutes").
			Comment("Usage minutes from start time to end time").
			Immutable(),
		id.Field("connectorID").
			Comment("ID of the connector").
			NotEmpty().
			Immutable(),
		field.String("clusterName").
			Comment("Cluster name for current cost").
			NotEmpty().
			Immutable(),
		field.Float("totalCost").
			Comment("Cost number").
			Default(0).
			Min(0),
		field.Int("currency").
			Comment("Cost currency").
			Optional(),
		field.Float("cpuCost").
			Comment("CPU cost for current cost").
			Default(0).
			Min(0),
		field.Float("gpuCost").
			Comment("GPU cost for current cost").
			Default(0).
			Min(0),
		field.Float("ramCost").
			Comment("Ram cost for current cost").
			Default(0).
			Min(0),
		field.Float("storageCost").
			Comment("Storage cost for current cost").
			Default(0).
			Min(0),
		field.Float("allocationCost").
			Comment("Allocation cost for current cost").
			Default(0).
			Min(0),
		field.Float("idleCost").
			Comment("Idle cost for current cost").
			Default(0).
			Min(0),
		field.Float("managementCost").
			Comment("Storage cost for current cost").
			Default(0).
			Min(0),
	}
}

func (ClusterCost) Edges() []ent.Edge {
	return []ent.Edge{
		// connector 1-* cluster costs.
		edge.From("connector", Connector.Type).
			Comment("Connector current cost linked").
			Ref("clusterCosts").
			Field("connectorID").
			Unique().
			Required().
			Immutable(),
	}
}
