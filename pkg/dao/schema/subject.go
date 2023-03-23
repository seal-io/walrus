package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
)

type Subject struct {
	ent.Schema
}

func (Subject) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Time{},
	}
}

func (Subject) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("kind", "group", "name").
			Unique(),
	}
}

func (Subject) Fields() []ent.Field {
	return []ent.Field{
		field.String("kind").
			Comment("The kind of the subject.").
			Default("user").
			Immutable(),
		field.String("group").
			Comment("The group of the subject.").
			Default("default"),
		field.String("name").
			Comment("The name of the subject.").
			NotEmpty().
			Immutable(),
		field.String("description").
			Comment("The detail of the subject.").
			Optional(),
		field.Bool("mountTo").
			Comment("Indicate whether the user mount to the group.").
			Nillable().
			Default(false),
		field.Bool("loginTo").
			Comment("Indicate whether the user login to the group.").
			Nillable().
			Default(true),
		field.JSON("roles", types.SubjectRoles{}).
			Comment("The role list of the subject.").
			Default(types.DefaultSubjectRoles()),
		field.Strings("paths").
			Comment("The path of the subject from the root group to itself.").
			Default([]string{}),
		field.Bool("builtin").
			Comment("Indicate whether the subject is builtin.").
			Default(false),
	}
}
