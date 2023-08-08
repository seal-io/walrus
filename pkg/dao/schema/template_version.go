package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
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

func (TemplateVersion) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("template_id").
			Comment("ID of the template.").
			NotEmpty().
			Immutable(),
		// Redundant template name reduce the number of queries.
		field.String("template_name").
			Comment("Name of the template.").
			NotEmpty().
			Immutable().Annotations(
			entx.SkipIO(),
		),
		field.String("version").
			Comment("Template version.").
			NotEmpty().
			Immutable(),
		// This is the normalized terraform module source that can be directly applied to terraform configuration.
		// For example, when we store multiple versions of a module in a mono repo,
		//   Template.Source = "github.com/foo/bar"
		//   TemplateVersion.Source = "github.com/foo/bar/1.0.0"
		field.String("source").
			Comment("Template version source.").
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
			Unique().
			Required().
			Immutable(),
	}
}
