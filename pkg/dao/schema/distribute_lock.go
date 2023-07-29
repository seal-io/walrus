package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type DistributeLock struct {
	ent.Schema
}

func (DistributeLock) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

func (DistributeLock) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("ID is the lock key.").
			Unique().
			NotEmpty().
			Immutable(),
		field.Int64("expireAt").
			Comment("Expiration timestamp to prevent the lock be occupied for long time."),
	}
}
