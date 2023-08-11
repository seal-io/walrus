package schema

import (
	"context"
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/utils/strs"
)

// CostReport holds the schema definition for the cluster resource item hourly cost.
type CostReport struct {
	ent.Schema
}

func (CostReport) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("start_time", "end_time", "connector_id", "fingerprint").
			Unique(),
	}
}

func (CostReport) Fields() []ent.Field {
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
		field.String("name").
			Comment("Resource name for current cost, could be __unmounted__.").
			Immutable(),
		field.String("fingerprint").
			Comment("String generated from resource properties, used to identify this cost.").
			Immutable(),
		// For k8s.
		field.String("cluster_name").
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
		field.String("controller_kind").
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
		// Cost.
		field.Float("totalCost").
			Comment("Cost number.").
			Default(0),
		field.Int("currency").
			Comment("Cost currency.").
			Optional(),
		field.Float("cpu_cost").
			Comment("Cpu cost for current cost.").
			Default(0),
		field.Float("cpu_core_request").
			Comment("Cpu core requested.").
			Default(0).
			Immutable(),
		field.Float("gpu_cost").
			Comment("GPU cost for current cost.").
			Default(0),
		field.Float("gpu_count").
			Comment("GPU core count.").
			Default(0).
			Immutable(),
		field.Float("ram_cost").
			Comment("Ram cost for current cost.").
			Default(0),
		field.Float("ram_byte_request").
			Comment("Ram requested in byte.").
			Default(0).
			Immutable(),
		field.Float("pv_cost").
			Comment("PV cost for current cost linked.").
			Default(0),
		field.Float("pv_bytes").
			Comment("PV bytes for current cost linked.").
			Default(0),
		field.Float("load_balancer_cost").
			Comment("LoadBalancer cost for current cost linked.").
			Default(0),
		// Usage.
		field.Float("cpu_core_usage_average").
			Comment("CPU core average usage.").
			Default(0),
		field.Float("cpu_core_usage_max").
			Comment("CPU core max usage.").
			Default(0),
		field.Float("ram_byte_usage_average").
			Comment("Ram average usage in byte.").
			Default(0),
		field.Float("ram_byte_usage_max").
			Comment("Ram max usage in byte.").
			Default(0),
	}
}

func (CostReport) Edges() []ent.Edge {
	return []ent.Edge{
		// Connector 1-* CostReports.
		edge.From("connector", Connector.Type).
			Comment("Connector current cost item linked.").
			Ref("cost_reports").
			Field("connector_id").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.SkipIO()),
	}
}

func (CostReport) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.SkipClearingOptionalField(),
	}
}

func (CostReport) Hooks() []ent.Hook {
	// Generate fingerprint by conditions.
	generateFingerprint := func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !m.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne) {
				return n.Mutate(ctx, m)
			}

			var cn, nd, ns, pd, cd, name string

			if v, ok := m.Field("cluster_name"); ok {
				cn = v.(string)
			}

			if v, ok := m.Field("node"); ok {
				nd = v.(string)
			}

			if v, ok := m.Field("namespace"); ok {
				ns = v.(string)
			}

			if v, ok := m.Field("pod"); ok {
				pd = v.(string)
			}

			if v, ok := m.Field("container"); ok {
				cd = v.(string)
			}

			if v, ok := m.Field("name"); ok {
				name = v.(string)
			}

			var err error
			if types.IsIdleOrManagementCost(name) {
				err = m.SetField("fingerprint", strs.Join("/", cn, name))
			} else {
				err = m.SetField("fingerprint", strs.Join("/", cn, nd, ns, pd, cd))
			}

			if err != nil {
				return nil, fmt.Errorf("error generating fingerprint: %w", err)
			}

			return n.Mutate(ctx, m)
		})
	}

	return []ent.Hook{
		generateFingerprint,
	}
}
