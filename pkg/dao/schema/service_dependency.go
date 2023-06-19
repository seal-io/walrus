package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

type ServiceDependency struct {
	ent.Schema
}

func (ServiceDependency) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.CreateTime{},
	}
}

func (ServiceDependency) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("serviceID", "dependentID", "path").
			Unique(),
	}
}

func (ServiceDependency) Fields() []ent.Field {
	return []ent.Field{
		oid.Field("serviceID").
			Comment("ID of the service.").
			NotEmpty().
			Immutable(),
		oid.Field("dependentID").
			Comment("service ID is dependent by the service.").
			NotEmpty(),
		field.JSON("path", []oid.ID{}).
			Comment("dependency path of service."),
		field.String("type").
			Comment("Type of the service dependency.").
			NotEmpty().
			Immutable(),
	}
}

func (ServiceDependency) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("service", Service.Type).
			Ref("dependencies").
			Field("serviceID").
			Comment("Services of the dependency.").
			Unique().
			Required().
			Immutable(),
	}
}
