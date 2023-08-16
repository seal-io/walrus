package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
)

// Perspective holds the schema definition for cost perspectives.
type Perspective struct {
	ent.Schema
}

func (Perspective) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
	}
}

func (Perspective) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}

func (Perspective) Fields() []ent.Field {
	return []ent.Field{
		field.String("start_time").
			Comment("Start time for the perspective.").
			NotEmpty(),
		field.String("end_time").
			Comment("End time for the perspective.").
			NotEmpty(),
		field.Bool("builtin").
			Comment("Is builtin perspective.").
			Default(false),
		field.JSON("cost_queries", []types.QueryCondition{}).
			Comment("Indicated the perspective included cost queries, record the used query condition.").
			Default([]types.QueryCondition{}),
	}
}
