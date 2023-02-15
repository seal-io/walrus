package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/oid"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type ApplicationResource struct {
	schema
}

func (ApplicationResource) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Meta{},
		mixin.Status{},
		mixin.Time{},
	}
}

func (ApplicationResource) Fields() []ent.Field {
	return []ent.Field{
		oid.Field("applicationID").
			Comment("ID of the application to which the revision belongs.").
			Immutable(),
		field.String("module").
			Comment("Module that generates the resource."),
		field.String("type").
			Comment("Resource type."),
	}
}
