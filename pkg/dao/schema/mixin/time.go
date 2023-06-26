package mixin

import (
	stdtime "time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/seal-io/seal/pkg/dao/schema/io"
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
		index.Fields("createTime"),
	}
}

func (i time) Fields() []ent.Field {
	fs := []ent.Field{
		field.Time("createTime").
			Nillable().
			Default(stdtime.Now).
			Immutable().
			Annotations(
				io.DisableInput()),
	}

	if !i.withoutUpdateTime {
		fs = append(fs,
			field.Time("updateTime").
				Nillable().
				Default(stdtime.Now).
				UpdateDefault(stdtime.Now).
				Annotations(
					io.DisableInput()),
		)
	}

	return fs
}
