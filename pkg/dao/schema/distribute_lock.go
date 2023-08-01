package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type DistributeLock struct {
	ent.Schema
}

func (DistributeLock) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("ID is the lock key.").
			NotEmpty().
			Immutable(),
		field.Int64("expireAt").
			Comment("Expiration timestamp to prevent the lock be occupied for long time."),
	}
}
