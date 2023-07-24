package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

type EnvironmentConnectorRelationship struct {
	ent.Schema
}

func (EnvironmentConnectorRelationship) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time().WithoutUpdateTime(),
	}
}

func (EnvironmentConnectorRelationship) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("environment_id", "connector_id").
			Unique(),
	}
}

func (EnvironmentConnectorRelationship) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("environment_id").
			Comment("ID of the environment to which the relationship connects.").
			NotEmpty().
			Immutable(),
		object.IDField("connector_id").
			Comment("ID of the connector to which the relationship connects.").
			NotEmpty().
			Immutable(),
	}
}

func (EnvironmentConnectorRelationship) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("environment", Environment.Type).
			Field("environment_id").
			Comment("Environment that connect to the relationship.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
		edge.To("connector", Connector.Type).
			Field("connector_id").
			Comment("Connector that connect to the relationship.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entx.Input(entx.WithUpdate())),
	}
}

func (EnvironmentConnectorRelationship) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.SkipClearingOptionalField(),
	}
}
