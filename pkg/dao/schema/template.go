package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type Template struct {
	ent.Schema
}

func (Template) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time(),
		mixin.LegacyStatus(),
	}
}

func (Template) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("It is also the name of the template.").
			Unique().
			NotEmpty().
			Immutable(),
		field.String("description").
			Comment("Description of the template.").
			Optional(),
		field.String("icon").
			Comment("A URL to an SVG or PNG image to be used as an icon.").
			Optional(),
		field.JSON("labels", map[string]string{}).
			Comment("Labels of the template.").
			Default(map[string]string{}),
		// For terraform deployer, this is a superset of terraform module git source.
		field.String("source").
			Comment("Source of the template.").
			NotEmpty(),
	}
}

func (Template) Edges() []ent.Edge {
	return []ent.Edge{
		// Template 1-* TemplateVersions.
		edge.To("versions", TemplateVersion.Type).
			Comment("Versions of the template.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
	}
}
