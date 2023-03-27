package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Perspective holds the schema definition for cost perspectives.
type Perspective struct {
	ent.Schema
}

func (Perspective) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Time{},
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
		field.String("name").
			Comment("Name for current perspective.").
			Unique().
			NotEmpty(),
		field.String("startTime").
			Comment("Start time for current perspective.").
			NotEmpty(),
		field.String("endTime").
			Comment("End time for current perspective.").
			NotEmpty(),
		field.Bool("builtin").
			Comment("Is builtin Perspective.").
			Default(false),
		field.JSON("allocationQueries", []types.QueryCondition{}).
			Comment("Indicated the perspective included allocation queries, record the used query condition.").
			Default([]types.QueryCondition{}),
	}
}
