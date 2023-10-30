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
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type ResourceComponentRelationship struct {
	ent.Schema
}

func (ResourceComponentRelationship) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time().WithoutUpdateTime(),
	}
}

func (ResourceComponentRelationship) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("resource_component_id", "dependency_id", "type").
			Unique(),
	}
}

func (ResourceComponentRelationship) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("resource_component_id").
			Comment("ID of the resource component.").
			NotEmpty().
			Immutable(),
		object.IDField("dependency_id").
			Comment("ID of the resource that resource depends on.").
			NotEmpty().
			Immutable(),
		field.String("type").
			Comment("Type of the relationship.").
			NotEmpty().
			Immutable(),
	}
}

func (ResourceComponentRelationship) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("resource_component", ResourceComponent.Type).
			Field("resource_component_id").
			Comment("ResourceComponent to which it currently belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		edge.To("dependency", ResourceComponent.Type).
			Field("dependency_id").
			Comment("ResourceComponent to which the dependency belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.Input(entx.WithUpdate())),
	}
}

func (ResourceComponentRelationship) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.SkipClearingOptionalField(),
	}
}
