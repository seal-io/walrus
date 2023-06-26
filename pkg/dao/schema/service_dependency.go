package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/io"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

type ServiceDependency struct {
	ent.Schema
}

func (ServiceDependency) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time().WithoutUpdateTime(),
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
			Comment("ID of the service that dependent by the service specified by serviceID.").
			NotEmpty(),
		field.JSON("path", []oid.ID{}).
			Comment("ID list (from root to leaf) of the service that " +
				"dependent by the service specified by serviceID."),
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
			Comment("Service to which the dependency belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				io.DisableInput()),
	}
}
