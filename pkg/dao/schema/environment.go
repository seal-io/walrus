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
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type Environment struct {
	ent.Schema
}

func (Environment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
	}
}

func (Environment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "name").
			Unique(),
	}
}

func (Environment) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("project_id").
			Comment("ID of the project to belong.").
			NotEmpty().
			Immutable(),
		field.String("type").
			Comment("Type of the environment.").
			NotEmpty().
			Immutable(),
	}
}

func (Environment) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* Environments.
		edge.From("project", Project.Type).
			Ref("environments").
			Field("project_id").
			Comment("Project to which the environment belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor)),
		// Environments *-* Connectors.
		edge.To("connectors", Connector.Type).
			Comment("Connectors that configure to the environment.").
			Through("environment_connector_relationships", EnvironmentConnectorRelationship.Type),
		// Environment 1-* Services.
		edge.To("services", Service.Type).
			Comment("Services that belong to the environment.").
			StructTag(`json:"services,omitempty,cli-ignore"`).
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.SkipInput(entx.WithUpdate(), entx.WithQuery()),
				entx.SkipOutput()),
		// Environment 1-* ServiceRevisions.
		edge.To("service_revisions", ServiceRevision.Type).
			Comment("ServicesRevisions that belong to the environment.").
			Annotations(
				entsql.OnDelete(entsql.NoAction),
				entx.SkipIO()),
		// Environment 1-* ServiceResources.
		edge.To("service_resources", ServiceResource.Type).
			Comment("ServiceResources that belong to the environment.").
			Annotations(
				entsql.OnDelete(entsql.NoAction),
				entx.SkipIO()),
		// Environment 1-* Variables.
		edge.To("variables", Variable.Type).
			Comment("Variables that belong to the environment.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipInput(entx.WithUpdate(), entx.WithQuery()),
				entx.SkipOutput()),
	}
}

func (Environment) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}
