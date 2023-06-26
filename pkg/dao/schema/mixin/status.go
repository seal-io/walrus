package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"github.com/seal-io/seal/pkg/dao/schema/io"
	daostatus "github.com/seal-io/seal/pkg/dao/types/status"
)

func Status() status {
	return status{}
}

type status struct {
	mixin.Schema
}

func (i status) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("status", daostatus.Status{}).
			Optional().
			Annotations(
				io.DisableInput()),
	}
}

func LegacyStatus() legacyStatus {
	return legacyStatus{}
}

type legacyStatus struct {
	mixin.Schema
}

func (i legacyStatus) Fields() []ent.Field {
	return []ent.Field{
		field.String("status").
			Optional().
			Annotations(
				io.DisableInput()),
		field.String("statusMessage").
			Optional().
			Annotations(
				io.DisableInput()),
	}
}
