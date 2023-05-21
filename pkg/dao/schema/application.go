package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/property"
)

type Application struct {
	ent.Schema
}

func (Application) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.OwnByProject{},
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
		property.SchemasField("variables").
			Comment("Variables definition of the application, " +
				"the variables of instance derived by this definition").
			Optional(),
	}
}

func (Application) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* applications.
		edge.From("project", Project.Type).
			Ref("applications").
			Field("projectID").
			Comment("Project to which the application belongs.").
			Unique().
			Required().
			Immutable(),
		// Application 1-* application instances.
		edge.To("instances", ApplicationInstance.Type).
			Comment("Application instances that belong to the application.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
		// Applications *-* modules.
		edge.To("modules", Module.Type).
			Comment("Modules that configure to the application.").
			Through("applicationModuleRelationships", ApplicationModuleRelationship.Type),
	}
}
