package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

type Template struct {
	ent.Schema
}

func (Template) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata().
			WithoutAnnotations(),
		mixin.Status(),
	}
}

func (Template) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}

func (Template) Fields() []ent.Field {
	return []ent.Field{
		field.String("icon").
			Comment("A URL to an SVG or PNG image to be used as an icon.").
			Annotations(entx.SkipInput()).
			Optional(),
		// For terraform deployer, this is a superset of terraform module git source.
		field.String("source").
			Comment("Source of the template.").
			NotEmpty(),
		object.IDField("catalog_id").
			Comment("ID of the template catalog.").
			Optional().
			Immutable(),
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
		// Catalog 1-* Templates.
		edge.From("catalog", Catalog.Type).
			Ref("templates").
			Field("catalog_id").
			Comment("Catalog to which the template belongs.").
			Unique().
			Immutable().
			Annotations(entx.SkipInput()),
	}
}
