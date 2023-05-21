package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

type ApplicationInstance struct {
	ent.Schema
}

func (ApplicationInstance) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.OwnByProject{},
		mixin.Time{},
	}
}

func (ApplicationInstance) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("applicationID", "environmentID", "name").
			Unique(),
	}
}

func (ApplicationInstance) Fields() []ent.Field {
	return []ent.Field{
		oid.Field("applicationID").
			Comment("ID of the application to which the instance belongs.").
			NotEmpty().
			Immutable(),
		oid.Field("environmentID").
			Comment("ID of the environment to which the instance deploys.").
			NotEmpty().
			Immutable(),
		field.String("name").
			Comment("Name of the instance.").
			NotEmpty().
			Immutable(),
		property.ValuesField("variables").
			Comment("Variables of the instance.").
			Optional(),
		field.JSON("status", status.Status{}).
			Comment("Status of the instance.").
			Optional(),
	}
}

func (ApplicationInstance) Edges() []ent.Edge {
	return []ent.Edge{
		// Application 1-* application instances.
		edge.From("application", Application.Type).
			Ref("instances").
			Field("applicationID").
			Comment("Application to which the instance belongs.").
			Unique().
			Required().
			Immutable(),
		// Environment 1-* application instances.
		edge.From("environment", Environment.Type).
			Ref("instances").
			Field("environmentID").
			Comment("Environment to which the instance belongs.").
			Unique().
			Required().
			Immutable(),
		// Application instance 1-* application revisions.
		edge.To("revisions", ApplicationRevision.Type).
			Comment("Application revisions that belong to the instance.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		// Application instance 1-* application resources.
		edge.To("resources", ApplicationResource.Type).
			Comment("Application resources that belong to the instance.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
