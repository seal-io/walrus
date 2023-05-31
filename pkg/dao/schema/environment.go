package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type Environment struct {
	ent.Schema
}

func (Environment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.OwnByProject{},
		mixin.Meta{},
		mixin.Time{},
	}
}

func (Environment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("projectID", "name").
			Unique(),
	}
}

func (Environment) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* environments.
		edge.From("project", Project.Type).
			Ref("environments").
			Field("projectID").
			Comment("Project to which the environment belongs.").
			Unique().
			Required().
			Immutable(),
		// Environments *-* connectors.
		edge.To("connectors", Connector.Type).
			Comment("Connectors that configure to the environment.").
			Through("environmentConnectorRelationships", EnvironmentConnectorRelationship.Type),
		// Environment 1-* services.
		edge.To("services", Service.Type).
			Comment("Services that belong to the environment.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
		// Environment 1-* service revisions.
		edge.To("serviceRevisions", ServiceRevision.Type).
			Comment("Services revisions that belong to the environment.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
	}
}
