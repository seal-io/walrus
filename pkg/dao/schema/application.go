package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
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
		oid.Field("projectID").
			Comment("ID of the project to which the application belongs.").
			NotEmpty().
			Immutable(),
		field.JSON("variables", []types.ApplicationVariable{}).
			Comment("Variables definition of the application, " +
				"the variables of instance derived by this definition").
			Optional(),
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
		// application 1-* application instances.
		edge.To("instances", ApplicationInstance.Type).
			Comment("Application instances that belong to this application.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
		// applications *-* modules.
		edge.To("modules", Module.Type).
			Comment("Modules that configure to the application.").
			Through("applicationModuleRelationships", ApplicationModuleRelationship.Type),
	}
}
