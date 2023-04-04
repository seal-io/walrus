package schema

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
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
		// application instance 1-* application resources.
		edge.From("instance", ApplicationInstance.Type).
			Ref("resources").
			Field("instanceID").
			Comment("Application instance to which the resource belongs.").
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

func (ApplicationResource) Interceptors() []ent.Interceptor {
	// filters out none "managed" mode resources.
	var filter = func(n model.Querier) model.Querier {
		return model.QuerierFunc(func(ctx context.Context, query model.Query) (model.Value, error) {
			var t, ok = query.(*model.ApplicationResourceQuery)
			if ok {
				// TODO: temporary store these resource but hidden while show the resource
				t.Where(applicationresource.And(
					applicationresource.Mode(types.ApplicationResourceModeManaged),
					applicationresource.TypeNEQ("kubectl_manifest"),
				))
			}
			return n.Query(ctx, query)
		})
	}

	return []ent.Interceptor{
		model.InterceptFunc(filter),
	}
}
