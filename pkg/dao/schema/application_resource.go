package schema

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

type ApplicationResource struct {
	ent.Schema
}

func (ApplicationResource) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Time{},
	}
}

func (ApplicationResource) Fields() []ent.Field {
	return []ent.Field{
		oid.Field("instanceID").
			Comment("ID of the application instance to which the resource belongs.").
			NotEmpty().
			Immutable(),
		oid.Field("connectorID").
			Comment("ID of the connector to which the resource deploys.").
			NotEmpty().
			Immutable(),
		oid.Field("compositionID").
			Comment("ID of the application resource to which the resource makes up, " +
				"it presents when mode is discovered.").
			Optional().
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
		field.String("deployerType").
			Comment("Type of deployer.").
			NotEmpty().
			Immutable(),
		field.JSON("status", types.ApplicationResourceStatus{}).
			Comment("Status of the resource.").
			Optional(),
	}
}

func (ApplicationResource) Edges() []ent.Edge {
	return []ent.Edge{
		// Application instance 1-* application resources.
		edge.From("instance", ApplicationInstance.Type).
			Ref("resources").
			Field("instanceID").
			Comment("Application instance to which the resource belongs.").
			Unique().
			Required().
			Immutable(),
		// Connector 1-* application resources.
		edge.From("connector", Connector.Type).
			Ref("resources").
			Field("connectorID").
			Comment("Connector to which the resource deploys.").
			Unique().
			Required().
			Immutable(),
		// Application resource(!discovered) 1-* application resources(discovered).
		edge.To("components", ApplicationResource.Type).
			Comment("Application resources that make up this resource.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}).
			From("composition").
			Field("compositionID").
			Comment("Application resource to which the resource makes up.").
			Unique().
			Immutable(),
	}
}

func (ApplicationResource) Interceptors() []ent.Interceptor {
	type target interface {
		WhereP(...func(*sql.Selector))
	}

	// Filters out not "data" mode and "kubectl_manifest" type resources.
	var filter = ent.TraverseFunc(func(ctx context.Context, query ent.Query) error {
		var t, ok = query.(target)
		if ok {
			t.WhereP(
				sql.FieldNEQ("mode", "data"),
				sql.FieldNEQ("type", "kubectl_manifest"),
			)
		}
		return nil
	})

	return []ent.Interceptor{
		filter,
	}
}
