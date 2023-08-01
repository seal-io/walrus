package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/schema/intercept"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/pkg/dao/types/property"
)

type Service struct {
	ent.Schema
}

func (Service) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
		mixin.Status(),
	}
}

func (Service) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "environment_id", "name").
			Unique(),
	}
}

func (Service) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("project_id").
			Comment("ID of the project to belong.").
			NotEmpty().
			Immutable(),
		object.IDField("environment_id").
			Comment("ID of the environment to which the service deploys.").
			NotEmpty().
			Immutable(),
		object.IDField("template_id").
			Comment("ID of the template version to which the service belong.").
			NotEmpty(),
		property.ValuesField("attributes").
			Comment("Attributes to configure the template.").
			Optional(),
	}
}

func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* Services.
		edge.From("project", Project.Type).
			Ref("services").
			Field("project_id").
			Comment("Project to which the service belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor)),
		// Environment 1-* Services.
		edge.From("environment", Environment.Type).
			Ref("services").
			Field("environment_id").
			Comment("Environment to which the service belongs.").
			Unique().
			Required().
			Immutable(),
		// TemplateVersion 1-* Services.
		edge.From("template", TemplateVersion.Type).
			Ref("services").
			Field("template_id").
			Comment("Template to which the service belongs.").
			Unique().
			Required().
			Annotations(
				entx.SkipInput(entx.WithQuery()),
				entx.Input(entx.WithCreate(), entx.WithUpdate())),
		// Service 1-* ServiceRevisions.
		edge.To("revisions", ServiceRevision.Type).
			Comment("Revisions that belong to the service.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Service 1-* ServiceResources.
		edge.To("resources", ServiceResource.Type).
			Comment("Resources that belong to the service.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Service 1-* Services (dependency).
		edge.To("dependencies", Service.Type).
			Comment("Dependencies that requires for the service.").
			Through("service_relationships", ServiceRelationship.Type).
			Annotations(
				entx.SkipIO()),
	}
}

func (Service) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}
