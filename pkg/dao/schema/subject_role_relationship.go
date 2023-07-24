package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

type SubjectRoleRelationship struct {
	ent.Schema
}

func (SubjectRoleRelationship) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time().WithoutUpdateTime(),
		mixin.OwnByProject().Optional(),
	}
}

func (SubjectRoleRelationship) Indexes() []ent.Index {
	return []ent.Index{
		// NB(thxCode): since null project subject roles belongs to the organization(beyond any project),
		// single unique constraint index cannot cover null column value,
		// so we leverage conditional indexes to run this case.
		index.Fields("project_id", "subject_id", "role_id").
			Unique().
			Annotations(
				entsql.IndexWhere("project_id IS NOT NULL")),
		index.Fields("subject_id", "role_id").
			Unique().
			Annotations(
				entsql.IndexWhere("project_id IS NULL")),
	}
}

func (SubjectRoleRelationship) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("subject_id").
			Comment("ID of the subject to which the relationship connects.").
			NotEmpty().
			Immutable(),
		field.String("role_id").
			Comment("ID of the role to which the relationship connects.").
			NotEmpty().
			Immutable(),
	}
}

func (SubjectRoleRelationship) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* Roles (Subject).
		edge.From("project", Project.Type).
			Ref("subject_roles").
			Field("project_id").
			Comment("Project to which the subject role belongs.").
			Unique().
			Immutable().
			Annotations(
				entx.SkipInput()),
		// Subject 1-1 Role.
		edge.To("subject", Subject.Type).
			Field("subject_id").
			Comment("Subject that connect to the relationship.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.Input(entx.WithUpdate())),
		edge.To("role", Role.Type).
			Field("role_id").
			Comment("Role that connect to the relationship.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.Input(entx.WithUpdate())),
	}
}

func (SubjectRoleRelationship) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.SkipClearingOptionalField(),
	}
}
