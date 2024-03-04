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
)

type ResourceDefinition struct {
	ent.Schema
}

func (ResourceDefinition) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
	}
}

func (ResourceDefinition) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("name").
			Unique(),
	}
}

func (ResourceDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").
			Comment("Type of the resources generated from the resource definition.").
			Immutable().
			StructTag(`json:"type,omitempty,cli-table-column"`),
		field.JSON("schema", types.Schema{}).
			Comment("Generated schema of the resource definition.").
			Default(types.Schema{}).
			Annotations(
				entx.SkipInput(),
			),
		field.JSON("ui_schema", &types.UISchema{}).
			Comment("UI schema of the resource definition.").
			Default(&types.UISchema{}).
			Optional().
			Annotations(
				entx.SkipClearingOptionalField(),
				entx.Input(entx.WithUpdate())),
		field.Bool("builtin").
			Comment("Indicate whether the resource definition is builtin, decided when creating.").
			Default(false).
			Immutable().
			Annotations(entx.SkipInput()).
			StructTag(`json:"builtin,omitempty,cli-table-column"`),
	}
}

func (ResourceDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		// ResourceDefinition *-* TemplateVersions.
		edge.To("matching_rules", TemplateVersion.Type).
			Comment("Template versions that configure to the resource definition.").
			Through("resource_definition_matching_rules", ResourceDefinitionMatchingRule.Type),
		// ResourceDefinition 1-* Resources.
		edge.To("resources", Resource.Type).
			Comment("Resources that use the definition.").
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.SkipIO()),
	}
}
