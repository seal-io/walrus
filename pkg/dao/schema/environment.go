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
			Immutable().
			StructTag(`json:"type,omitempty,cli-table-column"`),
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
		// Environment 1-* Resources.
		edge.To("resources", Resource.Type).
			Comment("Resources that belong to the environment.").
			StructTag(`json:"resources,omitempty,cli-ignore"`).
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.SkipInput(entx.WithUpdate(), entx.WithQuery()),
				entx.SkipOutput()),
		// Environment 1-* ResourceRevisions.
		edge.To("resource_revisions", ResourceRevision.Type).
			Comment("ResourceRevisions that belong to the environment.").
			Annotations(
				entsql.OnDelete(entsql.NoAction),
				entx.SkipIO()),
		// Environment 1-* ResourceComponents.
		edge.To("resource_components", ResourceComponent.Type).
			Comment("ResourceComponents that belong to the environment.").
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
