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
		mixin.Meta{},
		mixin.Time{},
	}
}

func (Environment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}

func (Environment) Edges() []ent.Edge {
	return []ent.Edge{
		// Environments *-* connectors.
		edge.To("connectors", Connector.Type).
			Comment("Connectors that configure to the environment.").
			Through("environmentConnectorRelationships", EnvironmentConnectorRelationship.Type),
		// Environment 1-* application instances.
		edge.To("instances", ApplicationInstance.Type).
			Comment("Application instances that belong to the environment.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
		// Environment 1-* application revisions.
		edge.To("revisions", ApplicationRevision.Type).
			Comment("Application revisions that belong to the environment.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
	}
}
