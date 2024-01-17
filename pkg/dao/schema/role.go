package schema

import (
	"context"
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
)

type Role struct {
	ent.Schema
}

func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time(),
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
			Immutable().
			StructTag(`json:"kind,omitempty,cli-table-column"`),
		field.String("description").
			Comment("The detail of the role.").
			Optional(),
		field.JSON("policies", types.RolePolicies{}).
			Comment("The policy list of the role.").
			Default(types.DefaultRolePolicies()),
		field.Strings("applicable_environment_types").
			Comment("The environment type list of the role to apply, only for system kind role.").
			Default([]string{}).
			Optional(),
		field.Bool("session").
			Comment("Indicate whether the role is session level, decide when creating.").
			Default(false).
			Immutable().
			Annotations(
				entx.SkipIO()),
		field.Bool("builtin").
			Comment("Indicate whether the role is builtin, decide when creating.").
			Default(false).
			Immutable().
			StructTag(`json:"builtin,omitempty,cli-table-column"`).
			Annotations(
				entx.SkipInput()),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		// Subjects *-* Roles.
		edge.From("subjects", Subject.Type).
			Ref("roles").
			Comment("Subjects to which the role configures.").
			Through("subject_role_relationships", SubjectRoleRelationship.Type).
			Annotations(
				entx.SkipIO()),
	}
}

func (Role) Hooks() []ent.Hook {
	// Normalize policies.
	normalizePolicies := func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !m.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne) {
				return n.Mutate(ctx, m)
			}

			if v, ok := m.Field("policies"); ok && len(v.(types.RolePolicies)) != 0 {
				policies := v.(types.RolePolicies).Normalize().Deduplicate().Sort()

				err := m.SetField("policies", policies)
				if err != nil {
					return nil, fmt.Errorf("error normalizing policies: %w", err)
				}
			}

			return n.Mutate(ctx, m)
		})
	}

	return []ent.Hook{
		normalizePolicies,
	}
}
