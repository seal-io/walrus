package mixin

import (
	stdtime "time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/seal-io/walrus/pkg/dao/entx"
)

func Time() time {
	return time{}
}

type time struct {
	mixin.Schema

	withoutUpdateTime bool
}

func (i time) WithoutUpdateTime() time {
	i.withoutUpdateTime = true
	return i
}

func (time) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("create_time"),
	}
}

func (i time) Fields() []ent.Field {
	fs := []ent.Field{
		field.Time("create_time").
			Nillable().
			Default(stdtime.Now).
			Immutable().
			Annotations(
				entx.SkipInput()),
	}

	if !i.withoutUpdateTime {
		fs = append(fs,
			field.Time("update_time").
				Nillable().
				Default(stdtime.Now).
				UpdateDefault(stdtime.Now).
				Annotations(
					entx.SkipInput()),
		)
	}

	return fs
}
