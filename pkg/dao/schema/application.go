package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/oid"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type Application struct {
	schema
}

func (Application) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Meta{},
		mixin.Time{},
	}
}

func (Application) Fields() []ent.Field {
	return []ent.Field{
		oid.Field("projectID").
			Comment("ID of the project to which the application belongs.").
			Immutable(),
		oid.Field("environmentID").
			Comment("ID of the environment to which the application deploys.").
			Immutable(),
		field.JSON("modules", []ApplicationModule{}).
			Comment("Application modules."),
	}
}

type ApplicationModule struct {
	Module    oid.ID
	Name      string
	Variables map[string]interface{}
}
