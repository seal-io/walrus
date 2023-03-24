package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// AllocationCost holds the schema definition for the cluster allocated resource hourly cost.
type AllocationCost struct {
	ent.Schema
}

func (AllocationCost) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("startTime", "endTime", "connectorID"),
		index.Fields("startTime", "endTime", "connectorID", "fingerprint").
			Unique(),
	}
}

func (AllocationCost) Fields() []ent.Field {
	return []ent.Field{
		field.Time("startTime").
			Comment("Usage start time for current cost.").
			Immutable(),
		field.Time("endTime").
			Comment("Usage end time for current cost.").
			Immutable(),
		field.Float("minutes").
			Comment("Usage minutes from start time to end time.").
			Immutable(),
		oid.Field("connectorID").
			Comment("ID of the connector.").
			NotEmpty().
			Immutable(),
		field.String("name").
			Comment("Resource name for current cost, could be __unmounted__.").
			Immutable(),
		field.String("fingerprint").
			Comment("String generated from resource properties, used to identify this cost.").
			Immutable(),
		// for k8s
		field.String("clusterName").
			Comment("Cluster name for current cost.").
			Optional().
			Immutable(),
		field.String("namespace").
			Comment("Namespace for current cost.").
			Optional().
			Immutable(),
		field.String("node").
			Comment("Node for current cost.").
			Optional().
			Immutable(),
		field.String("controller").
			Comment("Controller name for the cost linked resource.").
			Optional().
			Immutable(),
		field.String("controllerKind").
			Comment("Controller kind for the cost linked resource, deployment, statefulSet etc.").
			Optional().
			Immutable(),
		field.String("pod").
			Comment("Pod name for current cost.").
			Optional().
			Immutable(),
		field.String("container").
			Comment("Container name for current cost.").
			Optional().
			Immutable(),
		field.JSON("pvs", map[string]types.PVCost{}).
			Comment("PV list for current cost linked.").
			Default(map[string]types.PVCost{}).
			Immutable(),
		field.JSON("labels", map[string]string{}).
			Comment("Labels for the cost linked resource.").
			Default(map[string]string{}).
			Immutable(),
		// cost
		field.Float("totalCost").
			Comment("Cost number.").
			Default(0).
			Min(0),
		field.Int("currency").
			Comment("Cost currency.").
			Optional(),
		field.Float("cpuCost").
			Comment("Cpu cost for current cost.").
			Default(0).
			Min(0),
		field.Float("cpuCoreRequest").
			Comment("Cpu core requested.").
			Default(0).
			Min(0).
			Immutable(),
		field.Float("gpuCost").
			Comment("GPU cost for current cost.").
			Default(0).
			Min(0),
		field.Float("gpuCount").
			Comment("GPU core count.").
			Default(0).
			Min(0).
			Immutable(),
		field.Float("ramCost").
			Comment("Ram cost for current cost.").
			Default(0).
			Min(0),
		field.Float("ramByteRequest").
			Comment("Ram requested in byte.").
			Default(0).
			Min(0).
			Immutable(),
		field.Float("pvCost").
			Comment("PV cost for current cost linked.").
			Default(0).
			Min(0),
		field.Float("pvBytes").
			Comment("PV bytes for current cost linked.").
			Default(0).
			Min(0),
		field.Float("loadBalancerCost").
			Comment("LoadBalancer cost for current cost linked.").
			Default(0).
			Min(0),
		// usage
		field.Float("cpuCoreUsageAverage").
			Comment("CPU core average usage.").
			Default(0).
			Min(0),
		field.Float("cpuCoreUsageMax").
			Comment("CPU core max usage.").
			Default(0).
			Min(0),
		field.Float("ramByteUsageAverage").
			Comment("Ram average usage in byte.").
			Default(0).
			Min(0),
		field.Float("ramByteUsageMax").
			Comment("Ram max usage in byte.").
			Default(0).
			Min(0),
	}
}

func (AllocationCost) Edges() []ent.Edge {
	return []ent.Edge{
		// connector 1-* allocation cost.
		edge.From("connector", Connector.Type).
			Comment("Connector current cost linked.").
			Ref("allocationCosts").
			Field("connectorID").
			Unique().
			Required().
			Immutable(),
	}
}
