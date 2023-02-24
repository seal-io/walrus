package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/id"
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
		id.Field("applicationID").
			Comment("ID of the application to which the revision belongs.").
			NotEmpty().
			Immutable(),
		id.Field("environmentID").
			Comment("ID of the environment to which the application deploys.").
			NotEmpty().
			Immutable(),
		field.JSON("modules", []types.ApplicationModule{}).
			Comment("Application modules.").
			Default([]types.ApplicationModule{}),
		field.JSON("inputVariables", map[string]interface{}{}).
			Comment("Input variables of the revision.").
			Default(map[string]interface{}{}),
		field.String("inputPlan").
			Comment("Input plan of the revision."),
		field.String("output").
			Comment("Output of the revision."),
		field.String("deployerType").
			Comment("type of deployer").
			Default(types.AppRevisionDeployerTypeTF),
		field.Int("duration").
			Comment("deployment duration(seconds) of the of application revision").
			Default(0),
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
		// environment 1-* application revisions.
		edge.From("environment", Environment.Type).
			Ref("revisions").
			Field("environmentID").
			Comment("Environment to which the revision deploys.").
			Unique().
			Required().
			Immutable(),
	}
}
