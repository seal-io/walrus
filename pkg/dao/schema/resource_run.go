package schema

import (
	"context"
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/utils/strs"
)

type ResourceRun struct {
	ent.Schema
}

func (ResourceRun) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time().WithoutUpdateTime(),
		mixin.Status(),
	}
}

func (ResourceRun) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("project_id").
			Comment("ID of the project to belong.").
			NotEmpty().
			Immutable(),
		object.IDField("environment_id").
			Comment("ID of the environment to which the run belongs.").
			NotEmpty().
			Immutable(),
		object.IDField("resource_id").
			Comment("ID of the resource to which the run belongs.").
			NotEmpty().
			Immutable(),
		field.String("template_name").
			Comment("Name of the template.").
			NotEmpty().
			Immutable(),
		field.String("template_version").
			Comment("Version of the template.").
			NotEmpty(),
		object.IDField("template_id").
			Comment("ID of the template.").
			NotEmpty().
			Immutable(),
		property.ValuesField("attributes").
			Comment("Attributes to configure the template.").
			Optional(),
		property.ValuesField("computed_attributes").
			Comment("Computed attributes generated from attributes and schemas.").
			Optional(),
		crypto.MapField[string, string]("variables").
			Comment("Variables of the run.").
			Default(crypto.Map[string, string]{}),
		field.JSON("input_configs", map[string]types.ResourceRunConfigData{}).
			Comment("Input configs of the run.").
			Sensitive(),
		field.String("deployer_type").
			Comment("Type of deployer.").
			Default(types.DeployerTypeTF),
		field.Int("duration").
			Comment("Duration in seconds of the run deploying.").
			Default(0),
		field.JSON("previous_required_providers", []types.ProviderRequirement{}).
			Comment("Previous provider requirement of the run.").
			Default([]types.ProviderRequirement{}),
		field.Text("plan_record").
			Comment("Record of the run plan.").
			Optional(),
		field.Text("record").
			Comment("Record of the run.").
			Optional(),
		field.String("change_comment").
			Comment("Change comment of the run.").
			Optional(),
		field.String("created_by").
			Comment("User who created the run.").
			Annotations(entx.SkipInput()),
		field.String("type").
			Comment("Type of the run.").
			Annotations(entx.SkipInput()),
		field.Bool("preview").
			Comment("If the run is preview.").
			Default(false),
		field.JSON("annotations", map[string]string{}).
			Optional().
			Default(map[string]string{}).
			Annotations(
				entx.SkipInput(),
				entx.SkipOutput()),
		field.JSON("component_changes", []*types.ResourceComponentChange{}).
			Comment("Changes of the resource components.").
			Optional(),
		field.JSON("component_change_summary", types.ResourceComponentChangeSummary{}).
			Comment("Change summary of the resource.").
			Optional(),
	}
}

func (ResourceRun) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* ResourceRuns.
		edge.From("project", Project.Type).
			Ref("resource_runs").
			Field("project_id").
			Comment("Project to which the run belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor)),
		// Environment 1-* ResourceRuns.
		edge.From("environment", Environment.Type).
			Ref("resource_runs").
			Field("environment_id").
			Comment("Environment to which the run deploys.").
			Unique().
			Required().
			Immutable(),
		// Resource 1-* ResourceRuns.
		edge.From("resource", Resource.Type).
			Ref("runs").
			Field("resource_id").
			Comment("Resource to which the run belongs.").
			Unique().
			Required().
			Immutable(),
	}
}

func (ResourceRun) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}

func (ResourceRun) Hooks() []ent.Hook {
	// Normalize special chars in status message.
	normalizeStatusMessage := func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !m.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne) {
				return n.Mutate(ctx, m)
			}

			if v, ok := m.Field("record"); ok && v.(string) != "" {
				err := m.SetField("record", strs.NormalizeSpecialChars(v.(string)))
				if err != nil {
					return nil, fmt.Errorf("error normalizing record: %w", err)
				}
			}

			return n.Mutate(ctx, m)
		})
	}

	return []ent.Hook{
		normalizeStatusMessage,
	}
}
