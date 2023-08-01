package schema

import (
	"context"
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/schema/intercept"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

type Connector struct {
	ent.Schema
}

func (Connector) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
		mixin.Status(),
	}
}

func (Connector) Indexes() []ent.Index {
	return []ent.Index{
		// NB(thxCode): since null project connector belongs to the organization(beyond any project),
		// single unique constraint index cannot cover null column value,
		// so we leverage conditional indexes to run this case.
		index.Fields("project_id", "name").
			Unique().
			Annotations(
				entsql.IndexWhere("project_id IS NOT NULL")),
		index.Fields("name").
			Unique().
			Annotations(
				entsql.IndexWhere("project_id IS NULL")),
	}
}

func (Connector) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("project_id").
			Comment("ID of the project to belong, empty means for all projects.").
			Immutable().
			Optional(),
		field.String("category").
			Comment("Category of the connector.").
			Immutable().
			NotEmpty().
			Annotations(
				entx.Input()),
		field.String("type").
			Comment("Type of the connector.").
			NotEmpty().
			Immutable().
			Annotations(
				entx.Input()),
		field.String("config_version").
			Comment("Connector config version.").
			NotEmpty(),
		crypto.PropertiesField("config_data").
			Comment("Connector config data.").
			Default(crypto.Properties{}).
			Optional().
			Annotations(
				entx.SkipClearingOptionalField()),
		field.Bool("enable_fin_ops").
			Comment("Config whether enable finOps, will install prometheus and opencost while enable."),
		field.JSON("fin_ops_custom_pricing", &types.FinOpsCustomPricing{}).
			Comment("Custom pricing user defined.").
			Optional().
			Annotations(
				entx.SkipClearingOptionalField()),
	}
}

func (Connector) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* Connectors.
		edge.From("project", Project.Type).
			Ref("connectors").
			Field("project_id").
			Comment("Project to which the connector belongs.").
			Unique().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor)),
		// Environments *-* Connectors.
		edge.From("environments", Environment.Type).
			Ref("connectors").
			Comment("Environments to which the connector configures.").
			Through("environment_connector_relationships", EnvironmentConnectorRelationship.Type).
			Annotations(
				entx.SkipIO()),
		// Connector 1-* ServiceResources.
		edge.To("resources", ServiceResource.Type).
			Comment("ServiceResources that use the connector.").
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.SkipIO()),
		edge.To("cost_reports", CostReport.Type).
			Comment("CostReports that linked to the connection.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
	}
}

func (Connector) Hooks() []ent.Hook {
	// Set default pricing for Kubernetes connector.
	defaultK8sFinOps := func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !m.Op().Is(ent.OpCreate) {
				return n.Mutate(ctx, m)
			}

			if v, ok := m.Field("type"); !ok || v.(string) != types.ConnectorTypeK8s {
				return n.Mutate(ctx, m)
			}

			if v, ok := m.Field("fin_ops_custom_pricing"); !ok || v.(*types.FinOpsCustomPricing).IsZero() {
				err := m.SetField("fin_ops_custom_pricing", types.DefaultFinOpsCustomPricing())
				if err != nil {
					return nil, fmt.Errorf("error setting default finOps custom pricing: %w", err)
				}
			}

			return n.Mutate(ctx, m)
		})
	}

	return []ent.Hook{
		defaultK8sFinOps,
	}
}

func (Connector) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}
