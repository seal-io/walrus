package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
)

type Role struct {
	ent.Schema
}

func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Time{},
	}
}

func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("domain", "name").
			Unique(),
	}
}

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("domain").
			Comment("The domain of the role.").
			Default("system").
			Immutable(),
		field.String("name").
			Comment("The name of the role.").
			NotEmpty().
			Immutable(),
		field.String("description").
			Comment("The detail of the role.").
			Optional(),
		field.JSON("policies", types.RolePolicies{}).
			Comment("The policy list of the role.").
			Default(types.DefaultRolePolicies()),
		field.Bool("builtin").
			Comment("Indicate whether the subject is builtin, decide when creating.").
			Default(false).
			Immutable(),
		field.Bool("session").
			Comment("Indicate whether the subject is session level, decide when creating.").
			Default(false).
			Immutable(),
	}
}
