package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type Project struct {
	ent.Schema
}

func (Project) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
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
		// Project 1-* Environments.
		edge.To("environments", Environment.Type).
			Comment("Environments that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.SkipIO()),
		// Project 1-* Connectors.
		edge.To("connectors", Connector.Type).
			Comment("Connectors that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		// Project 1-* Roles (Subject).
		edge.To("subject_roles", SubjectRoleRelationship.Type).
			Comment("Roles of a subject that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipOutput()),
		// Project 1-* Services.
		edge.To("services", Service.Type).
			Comment("Services that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.NoAction),
				entx.SkipIO()),
		// Project 1-* ServiceRevisions.
		edge.To("service_revisions", ServiceRevision.Type).
			Comment("ServiceRevisions that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.NoAction),
				entx.SkipIO()),
		// Project 1-* Variables.
		edge.To("variables", Variable.Type).
			Comment("Variables that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
	}
}
