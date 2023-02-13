package schema

import (
	"entgo.io/ent"
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
		mixin.Meta{},
		mixin.Status{},
		mixin.Time{},
	}
}

func (ApplicationRevision) Fields() []ent.Field {
	return []ent.Field{
		oid.Field("applicationID", "ID of the application to which the revision belongs"),
		oid.Field("environmentID", "ID of the environment to which the application deploys"),
		field.JSON("modules", []ApplicationModule{}).
			Comment("Application modules"),
		field.JSON("inputVariables", map[string]interface{}{}).
			Comment("Input variables of the revision"),
		field.String("inputPlan").
			Comment("Input plan of the revision"),
		field.String("output").
			Comment("Output of the revision"),
	}
}
