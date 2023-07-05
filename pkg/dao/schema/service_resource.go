package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/io"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
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
		oid.Field("serviceID").
			Comment("ID of the service to which the resource belongs.").
			NotEmpty().
			Immutable(),
		oid.Field("connectorID").
			Comment("ID of the connector to which the resource deploys.").
			NotEmpty().
			Immutable(),
		oid.Field("compositionID").
			Comment("ID of the parent resource, " +
				"it presents when mode is discovered.").
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
		field.String("deployerType").
			Comment("Type of deployer.").
			NotEmpty().
			Immutable(),
		field.JSON("status", types.ServiceResourceStatus{}).
			Comment("Status of the resource.").
			Optional(),
	}
}

func (ServiceResource) Edges() []ent.Edge {
	return []ent.Edge{
		// Service 1-* service resources.
		edge.From("service", Service.Type).
			Ref("resources").
			Field("serviceID").
			Comment("Service to which the resource belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				io.DisableInput()),
		// Connector 1-* service resources.
		edge.From("connector", Connector.Type).
			Ref("resources").
			Field("connectorID").
			Comment("Connector to which the resource deploys.").
			Unique().
			Required().
			Immutable().
			Annotations(
				io.DisableInput()),
		// Service resource(!discovered) 1-* service resources(discovered).
		edge.To("components", ServiceResource.Type).
			Comment("Sub-resources that make up the resource.").
			From("composition").
			Field("compositionID").
			Comment("Service resource to which the resource makes up.").
			Unique().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade)),
	}
}
