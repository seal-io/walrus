package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/id"
)

type ApplicationResource struct {
	schema
}

func (ApplicationResource) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Status{},
		mixin.Time{},
	}
}

func (ApplicationResource) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("applicationID", "connectorID", "module", "mode", "type", "name").Unique(),
	}
}

func (ApplicationResource) Fields() []ent.Field {
	return []ent.Field{
		id.Field("applicationID").
			Comment("ID of the application to which the resource belongs.").
			NotEmpty().
			Immutable(),
		id.Field("connectorID").
			Comment("ID of the connector to which the resource deploys.").
			NotEmpty().
			Immutable(),
		field.String("module").
			Comment("Name of the module that generates the resource.").
			NotEmpty().
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
	}
}

func (ApplicationResource) Edges() []ent.Edge {
	return []ent.Edge{
		// application 1-* application resources.
		edge.From("application", Application.Type).
			Ref("resources").
			Field("applicationID").
			Comment("Application to which the resource belongs.").
			Unique().
			Required().
			Immutable(),
		// connector 1-* application resources.
		edge.From("connector", Connector.Type).
			Ref("resources").
			Field("connectorID").
			Comment("Connector to which the resource deploys.").
			Unique().
			Required().
			Immutable(),
	}
}
