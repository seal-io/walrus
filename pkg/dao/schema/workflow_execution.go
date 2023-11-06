package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type WorkflowExecution struct {
	ent.Schema
}

func (WorkflowExecution) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
		mixin.Status(),
	}
}

func (WorkflowExecution) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("project_id").
			Comment("ID of the project to belong.").
			NotEmpty().
			Immutable(),
		field.Int("version").
			Comment("Version of the workflow execution.").
			NonNegative().
			Annotations(
				entx.SkipInput()),
		field.String("type").
			Comment("Type of the workflow execution.").
			NotEmpty().
			Immutable().
			Annotations(
				entx.SkipInput()),
		object.IDField("workflow_id").
			Comment("ID of the workflow that this workflow execution belongs to.").
			NotEmpty().
			Immutable(),
		object.IDField("subject_id").
			Comment("ID of the subject that create workflow execution.").
			Immutable().
			Annotations(
				entx.SkipInput()),
		field.Time("execute_time").
			Comment("Time of the workflow execution started.").
			Optional().
			Annotations(
				entx.SkipInput()),
		field.Int("times").
			Comment("Number of times that this workflow execution has been executed.").
			NonNegative().
			Default(1),
		field.Int("duration").
			Comment("Duration seconds of the workflow execution.").
			NonNegative().
			Default(0),
		field.Int("parallelism").
			Comment("Number of task pods that can be executed in parallel of workflow.").
			Positive().
			Default(10),
		field.Int("timeout").
			Comment("Timeout of the workflow execution.").
			NonNegative().
			Default(0),
		field.JSON("trigger", types.WorkflowExecutionTrigger{}).
			Comment("Trigger of the workflow execution.").
			Default(types.WorkflowExecutionTrigger{}).
			Annotations(
				entx.SkipInput()),
	}
}

func (WorkflowExecution) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* WorkflowExecutions.
		edge.From("project", Project.Type).
			Ref("workflow_executions").
			Field("project_id").
			Comment("Project to which the workflow execution belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor),
				entx.SkipInput(entx.WithCreate(), entx.WithUpdate())),
		// WorkflowExecution 1-* WorkflowStageExecutions.
		edge.To("stages", WorkflowStageExecution.Type).
			Comment("Workflow stage executions that belong to this workflow execution.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipInput(entx.WithCreate(), entx.WithUpdate())),

		// Workflow 1-* WorkflowExecutions.
		edge.From("workflow", Workflow.Type).
			Ref("executions").
			Field("workflow_id").
			Comment("Workflow that this workflow execution belongs to.").
			Required().
			Unique().
			Immutable().
			Annotations(
				entx.SkipInput(entx.WithCreate(), entx.WithUpdate())),
	}
}

func (WorkflowExecution) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}
