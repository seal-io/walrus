package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
)

type ResourceDefinitionMatchingRule struct {
	ent.Schema
}

func (ResourceDefinitionMatchingRule) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time().WithoutUpdateTime(),
	}
}

func (ResourceDefinitionMatchingRule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("resource_definition_id", "name").
			Unique(),
		index.Fields("resource_definition_id", "template_id", "name").
			Unique(),
	}
}

func (ResourceDefinitionMatchingRule) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("resource_definition_id").
			Comment("ID of the resource definition to which the relationship connects.").
			NotEmpty().
			Immutable(),
		object.IDField("template_id").
			Comment("ID of the template version to which the relationship connects.").
			NotEmpty(),
		field.String("name").
			Comment("Name of the matching rule.").
			NotEmpty().
			Immutable().
			Annotations(
				entx.Input(entx.WithUpdate())),
		field.JSON("selector", types.Selector{}).
			Comment("Resource selector."),
		property.ValuesField("attributes").
			Comment("Attributes to configure the template.").
			Optional(),
		field.Int("order").
			Comment("Order of the matching rule.").
			Annotations(
				entx.SkipIO()),
	}
}

func (ResourceDefinitionMatchingRule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("resource_definition", ResourceDefinition.Type).
			Field("resource_definition_id").
			Comment("Resource definition that connect to the relationship.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		edge.To("template", TemplateVersion.Type).
			Field("template_id").
			Comment("Template version that connect to the relationship.").
			Unique().
			Required().
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.Input(entx.WithUpdate())),
		edge.To("resources", Resource.Type).
			Comment("Resources that match the rule.").
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.SkipIO()),
	}
}

func (ResourceDefinitionMatchingRule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.SkipClearingOptionalField(),
	}
}
