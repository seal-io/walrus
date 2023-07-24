package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

type ServiceResource struct {
	ent.Schema
}

func (ServiceResource) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time(),
		mixin.OwnByProject(),
	}
}

func (ServiceResource) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("service_id").
			Comment("ID of the service to which the resource belongs.").
			NotEmpty().
			Immutable(),
		object.IDField("connector_id").
			Comment("ID of the connector to which the resource deploys.").
			NotEmpty().
			Immutable(),
		object.IDField("composition_id").
			Comment("ID of the parent resource.").
			Optional().
			Immutable(),
		object.IDField("class_id").
			Comment("ID of the parent class of the resource realization.").
			Optional().
			Immutable(),
		field.String("mode").
			Comment("Mode that manages the generated resource, " +
				"it is the management way of the deployer to the resource, " +
				"which provides by deployer.").
			NotEmpty().
			Immutable(),
		field.String("type").
			Comment("Type of the generated resource, " +
				"it is the type of the resource which the deployer observes, " +
				"which provides by deployer.").
			NotEmpty().
			Immutable(),
		field.String("name").
			Comment("Name of the generated resource, " +
				"it is the real identifier of the resource, " +
				"which provides by deployer.").
			NotEmpty().
			Immutable(),
		field.String("deployer_type").
			Comment("Type of deployer.").
			NotEmpty().
			Immutable(),
		field.String("shape").
			Comment("Shape of the resource, it can be class or instance shape.").
			NotEmpty().
			Immutable(),
		field.JSON("status", types.ServiceResourceStatus{}).
			Comment("Status of the resource.").
			Optional(),
		field.JSON("keys", &types.ServiceResourceOperationKeys{}).
			Comment("Keys of the resource.").
			Optional().
			Annotations(
				entx.SkipInput(),
				entx.SkipStoringField()),
	}
}

func (ServiceResource) Edges() []ent.Edge {
	return []ent.Edge{
		// Service 1-* ServiceResources.
		edge.From("service", Service.Type).
			Ref("resources").
			Field("service_id").
			Comment("Service to which the resource belongs.").
			Unique().
			Required().
			Immutable(),
		// Connector 1-* ServiceResources.
		edge.From("connector", Connector.Type).
			Ref("resources").
			Field("connector_id").
			Comment("Connector to which the resource deploys.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.SkipInput()),
		// ServiceResource (!discovered) 1-* ServiceResources (discovered).
		edge.To("components", ServiceResource.Type).
			Comment("Components that makes up the service resource.").
			From("composition").
			Field("composition_id").
			Unique().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade)),
		// ServiceResource (class) 1-* ServiceResources (instance).
		edge.To("instances", ServiceResource.Type).
			Comment("Instances that realizes the service resource.").
			From("class").
			Field("class_id").
			Unique().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade)),
		// ServiceResource 1-* ServiceResource (dependency).
		edge.To("dependencies", ServiceResource.Type).
			Comment("Dependencies that requires for the service resource.").
			Through("service_resource_relationships", ServiceResourceRelationship.Type),
	}
}
