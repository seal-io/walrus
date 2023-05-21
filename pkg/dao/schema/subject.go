package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
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
		index.Fields("kind", "domain", "name").
			Unique(),
	}
}

func (Subject) Fields() []ent.Field {
	return []ent.Field{
		field.String("kind").
			Comment("The kind of the subject.").
			Default(types.SubjectKindUser).
			Immutable(),
		field.String("domain").
			Comment("The domain of the subject.").
			Default(types.SubjectDomainBuiltin),
		field.String("name").
			Comment("The name of the subject.").
			NotEmpty().
			Immutable(),
		field.String("description").
			Comment("The detail of the subject.").
			Optional(),
		field.Bool("builtin").
			Comment("Indicate whether the subject is builtin, decide when creating.").
			Default(false).
			Immutable(),
	}
}

func (Subject) Edges() []ent.Edge {
	return []ent.Edge{
		// Subject 1-* tokens.
		edge.To("tokens", Token.Type).
			Comment("Tokens that belong to the subject.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		// Subjects *-* roles.
		edge.To("roles", Role.Type).
			Comment("Roles that configure to the subject.").
			Through("subjectRoleRelationships", SubjectRoleRelationship.Type),
	}
}
