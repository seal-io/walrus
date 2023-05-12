package mixin

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Time struct {
	schema
}

func (Time) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("updateTime"),
	}
}

func (Time) Fields() []ent.Field {
	// Keep the json tag in camel case.
	return []ent.Field{
		field.Time("createTime").
			Comment("Describe creation time.").
			Nillable().
			Default(time.Now).
			Immutable(),
		field.Time("updateTime").
			Comment("Describe modification time.").
			Nillable().
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

type CreateTime struct {
	schema
}

func (CreateTime) Fields() []ent.Field {
	// Keep the json tag in camel case.
	return []ent.Field{
		field.Time("createTime").
			Comment("Describe creation time.").
			Nillable().
			Default(time.Now).
			Immutable(),
	}
}

type UpdateTime struct {
	schema
}

func (UpdateTime) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("updateTime"),
	}
}

func (UpdateTime) Fields() []ent.Field {
	// Keep the json tag in camel case.
	return []ent.Field{
		field.Time("updateTime").
			Comment("Describe modification time.").
			Nillable().
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
