package mixin

import (
	"context"
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"github.com/seal-io/walrus/pkg/dao/entx"
	daostatus "github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/strs"
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
				entx.SkipInput()),
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
				entx.SkipInput()),
		field.String("status_message").
			Optional().
			Annotations(
				entx.SkipInput()),
	}
}

func (legacyStatus) Hooks() []ent.Hook {
	// Normalize special chars in status message.
	normalizeStatusMessage := func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !m.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne) {
				return n.Mutate(ctx, m)
			}

			if v, ok := m.Field("status_message"); ok && v.(string) != "" {
				err := m.SetField("status_message", strs.NormalizeSpecialChars(v.(string)))
				if err != nil {
					return nil, fmt.Errorf("error normalizing status message: %w", err)
				}
			}

			return n.Mutate(ctx, m)
		})
	}

	return []ent.Hook{
		normalizeStatusMessage,
	}
}
