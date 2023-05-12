package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type Module struct {
	ent.Schema
}

func (Module) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Status{},
		mixin.Time{},
	}
}

func (Module) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("It is also the name of the module.").
			Unique().
			NotEmpty().
			Immutable(),
		field.String("description").
			Comment("Description of the module.").
			Optional(),
		field.String("icon").
			Comment("A URL to an SVG or PNG image to be used as an icon.").
			Optional(),
		field.JSON("labels", map[string]string{}).
			Comment("Labels of the module.").
			Default(map[string]string{}),
		// For terraform deployer, this is a superset of terraform module git source.
		field.String("source").
			Comment("Source of the module.").
			NotEmpty(),
	}
}

func (Module) Edges() []ent.Edge {
	return []ent.Edge{
		// Applications *-* modules.
		edge.From("applications", Application.Type).
			Ref("modules").
			Comment("Applications to which the module configures.").
			Through("applicationModuleRelationships", ApplicationModuleRelationship.Type),
		// Module 1-* module versions.
		edge.To("versions", ModuleVersion.Type).
			Comment("versions of the module.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
