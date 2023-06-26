package mixin

import (
	stdtime "time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"github.com/seal-io/seal/pkg/dao/schema/io"
	"github.com/seal-io/seal/pkg/dao/types/oid"
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
		oid.Field("id").
			Immutable().
			Annotations(
				io.DisableInputWhenCreating()),
		field.String("name").
			NotEmpty(),
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
					io.Disable()),
		)
	}

	fs = append(fs,
		field.Time("createTime").
			Nillable().
			Default(stdtime.Now).
			Immutable().
			Annotations(
				io.DisableInput()),
	)

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

func (metadata) Hooks() []ent.Hook {
	return []ent.Hook{
		oid.Hook(),
	}
}
