package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type Workflow struct {
	ent.Schema
}

func (Workflow) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
	}
}

func (Workflow) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "name").
			Unique(),
	}
}

func (Workflow) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("project_id").
			Comment("ID of the project that this workflow belongs to.").
			NotEmpty().
			Immutable(),
		object.IDField("environment_id").
			Comment("ID of the environment that this workflow belongs to.").
			Optional().
			Immutable(),
		field.String("type").
			Comment("Type of the workflow.").
			NotEmpty().
			Immutable(),
		field.Int("parallelism").
			Comment("Number of task pods that can be executed in parallel of workflow.").
			Positive().
			Default(10),
		field.Int("timeout").
			Comment("Timeout seconds of the workflow.").
			NonNegative().
			Default(0),
		field.Int("version").
			Comment("Execution version of the workflow.").
			NonNegative().
			Default(0).
			Annotations(
				entx.SkipInput()),
		field.JSON("variables", types.WorkflowVariables{}).
			Comment("Configs of workflow variables.").
			Optional(),
	}
}

func (Workflow) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* Environments.
		edge.From("project", Project.Type).
			Ref("workflows").
			Field("project_id").
			Comment("Project to which the workflow belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor)),
		// Workflow 1-* WorkflowStages.
		edge.To("stages", WorkflowStage.Type).
			Comment("Stages that belong to this workflow.").
			Annotations(
				entsql.OnDelete(entsql.Cascade)),
		// Workflow 1-* WorkflowExecutions.
		edge.To("executions", WorkflowExecution.Type).
			Comment("Workflow executions that belong to this workflow.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipInput()),
	}
}

func (Workflow) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}
