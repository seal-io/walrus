package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type WorkflowStageExecution struct {
	ent.Schema
}

func (WorkflowStageExecution) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
		mixin.Status(),
	}
}

func (WorkflowStageExecution) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("project_id").
			Comment("ID of the project to belong.").
			NotEmpty().
			Immutable(),
		object.IDField("workflow_id").
			Comment("ID of the workflow that this workflow execution belongs to.").
			NotEmpty().
			Immutable(),
		object.IDField("workflow_stage_id").
			Comment("ID of the workflow stage that this workflow stage execution belongs to.").
			Immutable(),
		object.IDField("workflow_execution_id").
			Comment("ID of the workflow execution that this workflow stage execution belongs to.").
			Immutable(),
		field.Time("execute_time").
			Comment("Time of the stage execution started.").
			Optional().
			Annotations(
				entx.SkipInput()),
		field.Int("duration").
			Comment("Duration of the workflow stage execution.").
			NonNegative().
			Default(0).
			Annotations(
				entx.SkipIO()),
		field.Int("order").
			Comment("Order of the workflow stage execution.").
			NonNegative().
			Default(0).
			Annotations(
				entx.SkipIO()),
	}
}

func (WorkflowStageExecution) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* WorkflowStageExecution.
		edge.From("project", Project.Type).
			Ref("workflow_stage_executions").
			Field("project_id").
			Comment("Project to which the workflow stage executions belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor),
				entx.SkipInput()),
		// WorkflowStageExecution 1-* WorkflowStepExecutions.
		edge.To("steps", WorkflowStepExecution.Type).
			Comment("Workflow step executions that belong to this workflow stage execution.").
			Annotations(
				entsql.OnDelete(entsql.Cascade)),
		// WorkflowExecution 1-* WorkflowStageExecutions.
		edge.From("workflow_execution", WorkflowExecution.Type).
			Ref("stages").
			Field("workflow_execution_id").
			Comment("Workflow execution that this workflow stage execution belongs to.").
			Required().
			Unique().
			Immutable(),
	}
}

func (WorkflowStageExecution) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}
