package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

type Variable struct {
	ent.Schema
}

func (Variable) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time(),
		mixin.OwnByProject().Optional(),
	}
}

func (Variable) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("environmentID", "name").
			Unique().
			Annotations(
				entsql.IndexWhere("environment_id IS NOT NULL")),
		index.Fields("projectID", "name").
			Unique().
			Annotations(
				entsql.IndexWhere("project_id IS NOT NULL AND environment_id IS NULL")),
		index.Fields("name").
			Unique().
			Annotations(
				entsql.IndexWhere("project_id IS NULL AND environment_id IS NULL")),
	}
}

func (Variable) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("The name of variable.").
			NotEmpty().
			Immutable(),
		crypto.StringField("value").
			Comment("The value of variable, store in string.").
			NotEmpty(),
		field.Bool("sensitive").
			Comment("The value is sensitive or not.").
			Default(false),
		field.String("description").
			Comment("Description of the variable.").
			Optional(),
		object.IDField("environmentID").
			Comment("ID of the environment to which the variable belongs to.").
			Optional().
			Immutable(),
	}
}

func (Variable) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* variables.
		edge.From("project", Project.Type).
			Ref("variables").
			Field("projectID").
			Comment("Project to which the variable belongs.").
			Unique().
			Immutable(),
		// Environment 1-* variables.
		edge.From("environment", Environment.Type).
			Ref("variables").
			Field("environmentID").
			Comment("Environment to which the variable belongs.").
			Unique().
			Immutable(),
	}
}
