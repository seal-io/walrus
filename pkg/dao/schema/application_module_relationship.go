package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	ents "entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
)

type ApplicationModuleRelationship struct {
	ent.Schema
}

func (ApplicationModuleRelationship) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (ApplicationModuleRelationship) Annotations() []ents.Annotation {
	return []ents.Annotation{
		field.ID("application_id", "module_id", "name"),
	}
}

func (ApplicationModuleRelationship) Fields() []ent.Field {
	return []ent.Field{
		oid.Field("application_id").
			Comment("ID of the application to which the relationship connects.").
			StructTag(`json:"applicationID" sql:"applicationID"`).
			NotEmpty().
			Immutable(),
		field.String("module_id").
			Comment("ID of the module to which the relationship connects.").
			StructTag(`json:"moduleID" sql:"moduleID"`).
			NotEmpty().
			Immutable(),
		field.String("version").
			Comment("Version of the module to which the relationship connects.").
			NotEmpty(),
		field.String("name").
			Comment("Name of the module customized to the application.").
			NotEmpty().
			Immutable(),
		property.ValuesField("attributes").
			Comment("Attributes to configure the module.").
			Optional(),
	}
}

func (ApplicationModuleRelationship) Edges() []ent.Edge {
	// NB(thxCode): entc cannot recognize camel case field name on edge with `Through`.
	return []ent.Edge{
		edge.To("application", Application.Type).
			Field("application_id").
			Comment("Application that connect to the relationship.").
			Unique().
			Required().
			Immutable().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("module", Module.Type).
			Field("module_id").
			Comment("Module that connect to the relationship.").
			Unique().
			Required().
			Immutable().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
	}
}
