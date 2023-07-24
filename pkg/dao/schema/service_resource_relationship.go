package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/io"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

type ServiceResourceRelationship struct {
	ent.Schema
}

func (ServiceResourceRelationship) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time().WithoutUpdateTime(),
	}
}

func (ServiceResourceRelationship) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("service_resource_id", "dependency_id", "type").
			Unique(),
	}
}

func (ServiceResourceRelationship) Fields() []ent.Field {
	return []ent.Field{
		object.Field("service_resource_id").
			Comment("ID of the service resource.").
			StructTag(`json:"serviceResourceID" sql:"serviceResourceID"`).
			NotEmpty().
			Immutable(),
		object.Field("dependency_id").
			Comment("ID of the resource that resource depends on.").
			StructTag(`json:"dependencyID" sql:"dependencyID"`).
			NotEmpty().
			Immutable(),
		field.String("type").
			Comment("Type of the relationship.").
			NotEmpty().
			Immutable(),
	}
}

func (ServiceResourceRelationship) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("serviceResource", ServiceResource.Type).
			Field("service_resource_id").
			Comment("Service Resource to which it currently belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				io.Disable()),
		edge.To("dependency", ServiceResource.Type).
			Field("dependency_id").
			Comment("Service resource to which the dependency belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade)),
	}
}
