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
	"github.com/seal-io/walrus/pkg/dao/types/property"
)

type Resource struct {
	ent.Schema
}

func (Resource) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
		mixin.Status(),
	}
}

func (Resource) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "environment_id", "name").
			Unique(),
	}
}

func (Resource) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("project_id").
			Comment("ID of the project to belong.").
			NotEmpty().
			Immutable(),
		object.IDField("environment_id").
			Comment("ID of the environment to which the resource deploys.").
			NotEmpty().
			Immutable(),
		object.IDField("template_id").
			Comment("ID of the template version to which the resource belong.").
			Optional().
			Nillable(),
		field.String("type").
			Comment("Type of the resource referring to a resource definition type.").
			Immutable().
			Optional().
			Annotations(
				entx.Input(entx.WithCreate(), entx.WithQuery()),
			),
		object.IDField("resource_definition_id").
			Comment("ID of the resource definition to which the resource use.").
			Optional().
			Nillable().
			Annotations(
				entx.SkipIO()),
		object.IDField("resource_definition_matching_rule_id").
			Comment("ID of the resource definition matching rule to which the resource use.").
			Optional().
			Nillable().
			Annotations(
				entx.SkipIO(),
			),
		property.ValuesField("attributes").
			Comment("Attributes to configure the template.").
			Optional(),
		property.ValuesField("computed_attributes").
			Comment("Computed attributes generated from attributes and schemas.").
			Optional(),
		field.JSON("endpoints", types.ResourceEndpoints{}).
			Comment("Endpoints of the resource.").
			Optional().
			StructTag(`json:"endpoints,omitempty,cli-table-column"`),
		field.String("change_comment").
			Comment("Change comment of the resource.").
			Optional().
			Annotations(
				entx.Input(entx.WithCreate(), entx.WithUpdate()),
				entx.SkipOutput()),
	}
}

func (Resource) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* Resources.
		edge.From("project", Project.Type).
			Ref("resources").
			Field("project_id").
			Comment("Project to which the resource belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor)),
		// Environment 1-* Resources.
		edge.From("environment", Environment.Type).
			Ref("resources").
			Field("environment_id").
			Comment("Environment to which the resource belongs.").
			Unique().
			Required().
			Immutable(),
		// TemplateVersion 1-* Resources.
		edge.From("template", TemplateVersion.Type).
			Ref("resources").
			Field("template_id").
			Comment("Template to which the resource belongs.").
			Unique().
			Annotations(
				entx.SkipInput(entx.WithQuery()),
				entx.Input(entx.WithCreate(), entx.WithUpdate())),
		// ResourceDefinition 1-* Resources.
		edge.From("resource_definition", ResourceDefinition.Type).
			Ref("resources").
			Field("resource_definition_id").
			Comment("Definition of the resource.").
			Unique().
			Annotations(
				entx.SkipInput(entx.WithQuery()),
				entx.Input(entx.WithCreate(), entx.WithUpdate()),
				entx.SkipOutput(),
			).
			// Hide the edge from the API, but generate the input for validation and edge resolution.
			// Mapping from type to definition_edge.type is done in the API layer.
			StructTag(`json:"-"`),
		// ResourceDefinitionMatchingRules 1-* Resources.
		edge.From("resource_definition_matching_rule", ResourceDefinitionMatchingRule.Type).
			Ref("resources").
			Field("resource_definition_matching_rule_id").
			Comment("Resource definition matching rule which the resource matches.").
			Unique().
			Annotations(
				entx.SkipInput(entx.WithQuery()),
				entx.Input(entx.WithCreate(), entx.WithUpdate()),
			),
		// Resource 1-* ResourceRevisions.
		edge.To("revisions", ResourceRevision.Type).
			Comment("Revisions that belong to the resource.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Resource 1-* ResourceComponents.
		edge.To("components", ResourceComponent.Type).
			Comment("Components that belong to the resource.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Resource 1-* Resources (dependency).
		edge.To("dependencies", Resource.Type).
			Comment("Dependencies that requires for the resource.").
			Through("resource_relationships", ResourceRelationship.Type).
			Annotations(
				entx.SkipIO()),
	}
}

func (Resource) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}
