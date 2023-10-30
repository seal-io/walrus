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

type ResourceRelationship struct {
	ent.Schema
}

func (ResourceRelationship) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time().WithoutUpdateTime(),
	}
}

func (ResourceRelationship) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("resource_id", "dependency_id", "path").
			Unique(),
	}
}

func (ResourceRelationship) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("resource_id").
			Comment("ID of the resource that deploys after the dependency finished.").
			NotEmpty().
			Immutable(),
		object.IDField("dependency_id").
			Comment("ID of the resource that deploys before the resource begins.").
			NotEmpty().
			Immutable(),
		field.JSON("path", []object.ID{}).
			Comment("ID list of the resource includes all dependencies and the resource itself.").
			Default([]object.ID{}).
			Immutable().
			Annotations(
				entx.SkipInput()),
		field.String("type").
			Comment("Type of the relationship.").
			NotEmpty().
			Immutable().
			Annotations(
				entx.SkipInput()),
	}
}

func (ResourceRelationship) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("resource", Resource.Type).
			Field("resource_id").
			Comment("Resource to which it currently belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		edge.To("dependency", Resource.Type).
			Field("dependency_id").
			Comment("Resource to which the dependency belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.Input(entx.WithUpdate())),
	}
}

func (ResourceRelationship) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.SkipClearingOptionalField(),
	}
}
