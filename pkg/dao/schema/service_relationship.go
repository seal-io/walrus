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

type ServiceRelationship struct {
	ent.Schema
}

func (ServiceRelationship) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time().WithoutUpdateTime(),
	}
}

func (ServiceRelationship) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("service_id", "dependency_id", "path").
			Unique(),
	}
}

func (ServiceRelationship) Fields() []ent.Field {
	return []ent.Field{
		object.Field("service_id").
			Comment("ID of the service that deploys after the dependency finished.").
			StructTag(`json:"serviceID" sql:"serviceID"`).
			NotEmpty().
			Immutable(),
		object.Field("dependency_id").
			Comment("ID of the service that deploys before the service begins.").
			StructTag(`json:"dependencyID" sql:"dependencyID"`).
			NotEmpty().
			Immutable(),
		field.JSON("path", []object.ID{}).
			Comment("ID list of the service includes all dependencies and the service itself.").
			Default([]object.ID{}).
			Immutable().
			Annotations(
				io.DisableInput()),
		field.String("type").
			Comment("Type of the relationship.").
			NotEmpty().
			Immutable().
			Annotations(
				io.DisableInput()),
	}
}

func (ServiceRelationship) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("service", Service.Type).
			Field("service_id").
			Comment("Service to which it currently belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				io.Disable()),
		edge.To("dependency", Service.Type).
			Field("dependency_id").
			Comment("Service to which the dependency belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Restrict)),
	}
}
