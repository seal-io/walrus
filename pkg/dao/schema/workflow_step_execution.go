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

type WorkflowStepExecution struct {
	ent.Schema
}

func (WorkflowStepExecution) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
		mixin.Status(),
	}
}

func (WorkflowStepExecution) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("workflow_step_id").
			Comment("ID of the workflow step that this workflow step execution belongs to.").
			Immutable(),
		object.IDField("workflow_execution_id").
			Comment("ID of the workflow execution that this workflow step execution belongs to.").
			Immutable(),
		object.IDField("workflow_stage_execution_id").
			Comment("ID of the workflow stage execution that this workflow step execution belongs to.").
			Immutable(),
		object.IDField("project_id").
			Comment("ID of the project to belong.").
			NotEmpty().
			Immutable(),
		object.IDField("workflow_id").
			Comment("ID of the workflow that this workflow step execution belongs to.").
			NotEmpty().
			Immutable(),
		field.String("type").
			Comment("Type of the workflow step execution.").
			NotEmpty().
			Immutable(),
		field.JSON("attributes", map[string]any{}).
			Comment("Attributes of the workflow step execution.").
			Optional(),
		field.Int("times").
			Comment("Number of times that this workflow step execution has been executed.").
			NonNegative().
			Default(1).
			Annotations(
				entx.SkipIO()),
		field.Time("execute_time").
			Comment("Time of the step execution started.").
			Optional().
			Annotations(
				entx.SkipInput()),
		field.Int("duration").
			Comment("Execution duration seconds of the workflow step.").
			NonNegative().
			Default(0).
			Annotations(
				entx.SkipInput()),
		field.JSON("retryStrategy", &types.RetryStrategy{}).
			Comment("Retry policy of the workflow step.").
			Optional(),
		field.Int("timeout").
			Comment("Timeout of the workflow step execution.").
			NonNegative().
			Default(0),
		field.Int("order").
			Comment("Order of the workflow step execution.").
			NonNegative().
			Default(0).
			Annotations(
				entx.SkipIO()),
		field.Text("record").
			Comment("Log record of the workflow step execution.").
			Default(""),
	}
}

func (WorkflowStepExecution) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* WorkflowStepExecutions.
		edge.From("project", Project.Type).
			Ref("workflow_step_executions").
			Field("project_id").
			Comment("Project to which the workflow step execution belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor)),
		// WorkflowStageExecution 1-* WorkflowStepExecutions.
		edge.From("stage_execution", WorkflowStageExecution.Type).
			Ref("steps").
			Field("workflow_stage_execution_id").
			Comment("Workflow stage execution that this workflow step execution belongs to.").
			Required().
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.SkipInput(entx.WithCreate(), entx.WithCreate())),
	}
}

func (WorkflowStepExecution) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}
