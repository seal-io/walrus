package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type TemplateVersion struct {
	ent.Schema
}

func (TemplateVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time(),
	}
}

func (TemplateVersion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "version", "template_id").
			Unique(),
	}
}

func (TemplateVersion) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("template_id").
			Comment("ID of the template.").
			NotEmpty().
			Immutable(),
		// Redundant template name reduce the number of queries.
		field.String("name").
			Comment("Name of the template.").
			NotEmpty().
			Immutable(),
		field.String("version").
			Comment("Version of the template.").
			NotEmpty().
			Immutable(),
		// This is the normalized terraform module source that can be directly applied to terraform configuration.
		// For example, when we store multiple versions of a module in a mono repo,
		//   Template.Source = "github.com/foo/bar"
		//   TemplateVersion.Source = "github.com/foo/bar/1.0.0"
		field.String("source").
			Comment("Source of the template.").
			NotEmpty().
			Immutable(),
		field.JSON("schema", &types.TemplateSchema{}).
			Comment("Schema of the template.").
			Default(&types.TemplateSchema{}),
	}
}

func (TemplateVersion) Edges() []ent.Edge {
	return []ent.Edge{
		// Template 1-* TemplateVersions.
		edge.From("template", Template.Type).
			Ref("versions").
			Field("template_id").
			Comment("Template to which the template version belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.SkipInput()),
		// TemplateVersion 1-* Services.
		edge.To("services", Service.Type).
			Comment("Services that belong to the template version.").
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.SkipIO()),
	}
}
