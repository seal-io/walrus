package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type WorkflowStep struct {
	ent.Schema
}

func (WorkflowStep) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
	}
}

func (WorkflowStep) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").
			Comment("Type of the workflow step.").
			NotEmpty().
			Immutable(),
		object.IDField("project_id").
			Comment("ID of the project to belong.").
			NotEmpty().
			Immutable(),
		object.IDField("workflow_id").
			Comment("ID of the workflow that this workflow step belongs to.").
			NotEmpty().
			Immutable(),
		object.IDField("workflow_stage_id").
			Comment("ID of the stage that this workflow step belongs to.").
			Immutable(),
		field.JSON("attributes", map[string]any{}).
			Comment("Attributes of the workflow step.").
			Optional(),
		field.JSON("inputs", map[string]any{}).
			Comment("Inputs of the workflow step.").
			Optional(),
		field.JSON("outputs", map[string]any{}).
			Comment("Outputs of the workflow step.").
			Optional(),
		field.Int("order").
			Comment("Order of the workflow step.").
			NonNegative().
			Default(0).
			Annotations(
				entx.SkipIO()),
		field.JSON("dependencies", []object.ID{}).
			Comment("ID list of the workflow steps that this workflow step depends on.").
			Default([]object.ID{}),
		field.JSON("retryStrategy", &types.RetryStrategy{}).
			Comment("Retry policy of the workflow step.").
			Optional(),
		field.Int("timeout").
			Comment("Timeout seconds of the workflow step, 0 means no timeout.").
			NonNegative().
			Default(0),
	}
}

func (WorkflowStep) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* WorkflowSteps.
		edge.From("project", Project.Type).
			Ref("workflow_steps").
			Field("project_id").
			Comment("Project to which the step belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor)),
		// WorkflowStage 1-* WorkflowSteps.
		edge.From("stage", WorkflowStage.Type).
			Ref("steps").
			Field("workflow_stage_id").
			Comment("Workflow stage that this workflow step belongs to.").
			Required().
			Unique().
			Immutable(),
	}
}

func (WorkflowStep) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}
