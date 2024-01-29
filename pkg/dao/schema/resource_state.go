package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type ResourceState struct {
	ent.Schema
}

func (ResourceState) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
	}
}

func (ResourceState) Fields() []ent.Field {
	return []ent.Field{
		field.String("data").
			Comment("State data of the resource.").
			Default("").
			Annotations(
				entx.SkipIO(),
			),
		object.IDField("resource_id").
			Comment("ID of the resource to which the state belongs.").
			NotEmpty().
			Immutable(),
	}
}

func (ResourceState) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("resource", Resource.Type).
			Ref("state").
			Field("resource_id").
			Comment("Resource to which the state belongs.").
			Unique().
			Required().
			Immutable(),
	}
}
