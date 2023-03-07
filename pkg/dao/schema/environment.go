package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type Environment struct {
	schema
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

func (Environment) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("variables", map[string]interface{}{}).
			Comment("Variables of the environment.").
			Optional(),
	}
}

func (Environment) Edges() []ent.Edge {
	return []ent.Edge{
		// environments *-* connectors.
		edge.To("connectors", Connector.Type).
			Comment("Connectors that configure to the environment.").
			Through("environmentConnectorRelationships", EnvironmentConnectorRelationship.Type),
		// environment 1-* applications.
		edge.To("applications", Application.Type).
			Comment("Applications that belong to the environment.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
		// environment 1-* application revisions.
		edge.To("revisions", ApplicationRevision.Type).
			Comment("Revisions that belong to the environment.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
	}
}
