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

type WorkflowStage struct {
	ent.Schema
}

func (WorkflowStage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
	}
}

func (WorkflowStage) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("project_id").
			Comment("ID of the project to belong.").
			NotEmpty().
			Immutable(),
		object.IDField("workflow_id").
			Comment("ID of the workflow that this workflow stage belongs to.").
			NotEmpty().
			Immutable(),
		field.JSON("dependencies", []object.ID{}).
			Comment("ID list of the workflow stages that this workflow stage depends on.").
			Default([]object.ID{}),
		field.Int("order").
			Comment("Order of the workflow stage.").
			NonNegative().
			Default(0).
			Annotations(
				entx.SkipIO()),
	}
}

func (WorkflowStage) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* WorkflowStages.
		edge.From("project", Project.Type).
			Ref("workflow_stages").
			Field("project_id").
			Comment("Project to which the workflow stage belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor)),
		// WorkflowStage 1-* WorkflowSteps.
		edge.To("steps", WorkflowStep.Type).
			Comment("Workflow steps that belong to this workflow stage.").
			Annotations(
				entsql.OnDelete(entsql.Cascade)),
		// Workflow 1-* WorkflowStages.
		edge.From("workflow", Workflow.Type).
			Ref("stages").
			Field("workflow_id").
			Comment("Workflow that this workflow stage belongs to.").
			Required().
			Unique().
			Immutable(),
	}
}

func (WorkflowStage) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}
