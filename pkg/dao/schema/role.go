package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
)

type Role struct {
	ent.Schema
}

func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("It is also the name of the role.").
			Unique().
			NotEmpty().
			Immutable(),
		field.String("kind").
			Comment("The kind of the role.").
			Default(types.RoleKindSystem).
			Immutable(),
		field.String("description").
			Comment("The detail of the role.").
			Optional(),
		field.JSON("policies", types.RolePolicies{}).
			Comment("The policy list of the role.").
			Default(types.DefaultRolePolicies()),
		field.Bool("session").
			Comment("Indicate whether the role is session level, decide when creating.").
			Default(false).
			Immutable(),
		field.Bool("builtin").
			Comment("Indicate whether the role is builtin, decide when creating.").
			Default(false).
			Immutable(),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		// Subjects *-* roles.
		edge.From("subjects", Subject.Type).
			Ref("roles").
			Comment("Subjects to which the role configures.").
			Through("subjectRoleRelationships", SubjectRoleRelationship.Type),
	}
}
