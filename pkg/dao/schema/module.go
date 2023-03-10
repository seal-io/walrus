package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
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
		field.String("source").
			Comment("Source of the module.").
			NotEmpty(),
		field.String("version").
			Comment("Version of the module.").
			Optional(),
		field.JSON("schema", &types.ModuleSchema{}).
			Comment("Schema of the module.").
			Default(&types.ModuleSchema{}),
	}
}

func (Module) Edges() []ent.Edge {
	return []ent.Edge{
		// applications *-* modules.
		edge.From("applications", Application.Type).
			Ref("modules").
			Comment("Applications to which the module configures.").
			Through("applicationModuleRelationships", ApplicationModuleRelationship.Type),
	}
}
