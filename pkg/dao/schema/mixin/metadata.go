package mixin

import (
	stdtime "time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"github.com/seal-io/seal/pkg/dao/entx"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

func Metadata() metadata {
	return metadata{}
}

type metadata struct {
	mixin.Schema

	withoutDescription bool
	withoutLabels      bool
	withoutAnnotations bool
	withoutUpdateTime  bool
}

func (i metadata) WithoutDescription() metadata {
	i.withoutDescription = true
	return i
}

func (i metadata) WithoutLabels() metadata {
	i.withoutLabels = true
	return i
}

func (i metadata) WithoutAnnotations() metadata {
	i.withoutAnnotations = true
	return i
}

func (i metadata) WithoutUpdateTime() metadata {
	i.withoutUpdateTime = true
	return i
}

func (i metadata) Fields() []ent.Field {
	fs := []ent.Field{
		object.IDField("id").
			Immutable(),
		field.String("name").
			NotEmpty().
			Immutable(),
	}

	if !i.withoutDescription {
		fs = append(fs,
			field.String("description").
				Optional(),
		)
	}

	if !i.withoutLabels {
		fs = append(fs,
			field.JSON("labels", map[string]string{}).
				Optional().
				Default(map[string]string{}),
		)
	}

	if !i.withoutAnnotations {
		fs = append(fs,
			field.JSON("annotations", map[string]string{}).
				Optional().
				Default(map[string]string{}).
				Annotations(
					entx.SkipInput(),
					entx.SkipOutput()),
		)
	}

	fs = append(fs,
		field.Time("create_time").
			Nillable().
			Default(stdtime.Now).
			Immutable().
			Annotations(
				entx.SkipInput()),
	)

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

func (metadata) Hooks() []ent.Hook {
	return []ent.Hook{
		object.IDHook(),
	}
}
