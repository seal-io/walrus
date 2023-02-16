package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/oid"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type ApplicationRevision struct {
	schema
}

func (ApplicationRevision) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Status{},
		mixin.Time{},
	}
}

func (ApplicationRevision) Fields() []ent.Field {
	return []ent.Field{
		oid.Field("applicationID").
			Comment("ID of the application to which the revision belongs.").
			Immutable(),
		oid.Field("environmentID").
			Comment("ID of the environment to which the application deploys, " +
				"uses for redundancy but not correlation constraint.").
			NotEmpty().
			Immutable(),
		field.JSON("modules", []ApplicationModule{}).
			Comment("Application modules."),
		field.JSON("inputVariables", map[string]interface{}{}).
			Comment("Input variables of the revision."),
		field.String("inputPlan").
			Comment("Input plan of the revision."),
		field.String("output").
			Comment("Output of the revision."),
	}
}

func (ApplicationRevision) Edges() []ent.Edge {
	return []ent.Edge{
		// application 1-* application revisions.
		edge.From("application", Application.Type).
			Ref("revisions").
			Field("applicationID").
			Comment("Application to which the revision belongs.").
			Unique().
			Required().
			Immutable(),
	}
}
