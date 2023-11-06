package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
)

type Project struct {
	ent.Schema
}

func (Project) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
	}
}

func (Project) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}

func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* Environments.
		edge.To("environments", Environment.Type).
			Comment("Environments that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.SkipIO()),
		// Project 1-* Connectors.
		edge.To("connectors", Connector.Type).
			Comment("Connectors that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Project 1-* Roles (Subject).
		edge.To("subject_roles", SubjectRoleRelationship.Type).
			Comment("Roles of a subject that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipOutput()),
		// Project 1-* Services.
		edge.To("services", Service.Type).
			Comment("Services that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.NoAction),
				entx.SkipIO()),
		// Project 1-* ServiceResources.
		edge.To("service_resources", ServiceResource.Type).
			Comment("ServiceResources that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.NoAction),
				entx.SkipIO()),
		// Project 1-* ServiceRevisions.
		edge.To("service_revisions", ServiceRevision.Type).
			Comment("ServiceRevisions that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.NoAction),
				entx.SkipIO()),
		// Project 1-* Variables.
		edge.To("variables", Variable.Type).
			Comment("Variables that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Project 1-* Templates.
		edge.To("templates", Template.Type).
			Comment("Templates that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Project 1-* TemplateVersions.
		edge.To("template_versions", TemplateVersion.Type).
			Comment("TemplateVersions that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Project 1-* Catalogs.
		edge.To("catalogs", Catalog.Type).
			Comment("Catalogs that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Project 1-* Workflows.
		edge.To("workflows", Workflow.Type).
			Comment("Workflows that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Project 1-* WorkflowStages.
		edge.To("workflow_stages", WorkflowStage.Type).
			Comment("WorkflowStages that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Project 1-* WorkflowSteps.
		edge.To("workflow_steps", WorkflowStep.Type).
			Comment("WorkflowSteps that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Project 1-* WorkflowExecutions.
		edge.To("workflow_executions", WorkflowExecution.Type).
			Comment("WorkflowExecutions that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Project 1-* WorkflowStageExecutions.
		edge.To("workflow_stage_executions", WorkflowStageExecution.Type).
			Comment("WorkflowStageExecutions that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Project 1-* WorkflowStepExecutions.
		edge.To("workflow_step_executions", WorkflowStepExecution.Type).
			Comment("WorkflowStepExecutions that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
	}
}

func (Project) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.ValidateContext(intercept.WithProjectInterceptor),
	}
}

func (Project) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("id"),
	}
}
