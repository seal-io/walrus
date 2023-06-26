package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/io"
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
		// Project 1-* environments.
		edge.To("environments", Environment.Type).
			Comment("Environments that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				io.Disable()),
		// Project 1-* connectors.
		edge.To("connectors", Connector.Type).
			Comment("Connectors that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				io.Disable()),
		// Project 1-* secrets.
		edge.To("secrets", Secret.Type).
			Comment("Secrets that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				io.Disable()),
		// Project 1-* subject roles.
		edge.To("subjectRoles", SubjectRoleRelationship.Type).
			Comment("Subject roles that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				io.Disable()),
		// Project 1-* services.
		edge.To("services", Service.Type).
			Comment("Services that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.NoAction),
				io.Disable()),
		// Project 1-* service revisions.
		edge.To("serviceRevisions", ServiceRevision.Type).
			Comment("Service revisions that belong to the project.").
			Annotations(
				entsql.OnDelete(entsql.NoAction),
				io.Disable()),
	}
}
