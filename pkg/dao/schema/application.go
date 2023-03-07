package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/id"
)

type Application struct {
	ent.Schema
}

func (Application) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Meta{},
		mixin.Time{},
	}
}

func (Application) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("projectID", "name").
			Unique(),
	}
}

func (Application) Fields() []ent.Field {
	return []ent.Field{
		id.Field("projectID").
			Comment("ID of the project to which the application belongs.").
			NotEmpty().
			Immutable(),
		id.Field("environmentID").
			Comment("ID of the environment to which the application deploys.").
			NotEmpty().
			Immutable(),
	}
}

func (Application) Edges() []ent.Edge {
	return []ent.Edge{
		// project 1-* applications.
		edge.From("project", Project.Type).
			Ref("applications").
			Field("projectID").
			Comment("Project to which this application belongs.").
			Unique().
			Required().
			Immutable(),
		// environment 1-* applications.
		edge.From("environment", Environment.Type).
			Ref("applications").
			Field("environmentID").
			Comment("Environment to which the application belongs.").
			Unique().
			Required().
			Immutable(),
		// application 1-* application resources.
		edge.To("resources", ApplicationResource.Type).
			Comment("Resources that belong to the application.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		// application 1-* application revisions.
		edge.To("revisions", ApplicationRevision.Type).
			Comment("Revisions that belong to this application.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		// applications *-* modules.
		edge.To("modules", Module.Type).
			Comment("Modules that configure to the application.").
			Through("applicationModuleRelationships", ApplicationModuleRelationship.Type),
	}
}
