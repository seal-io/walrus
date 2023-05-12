package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type Project struct {
	ent.Schema
}

func (Project) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Meta{},
		mixin.Time{},
	}
}

func (Project) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}

func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* applications.
		edge.To("applications", Application.Type).
			Comment("Applications that belong to the project.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
		// Project 1-* secrets.
		edge.To("secrets", Secret.Type).
			Comment("Secrets that belong to the project.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
