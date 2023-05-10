package oid

import (
	"context"
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"github.com/sony/sonyflake"

	"github.com/seal-io/seal/utils/vars"
)

// Config holds the config of the types.ID generation.
var Config = vars.SetOnce[sonyflake.Settings]{}

// Field returns a new ent.Field with type types.ID.
func Field(name string) *fieldBuilder {
	return &fieldBuilder{
		desc: field.String(name).
			GoType(ID("")).
			SchemaType(map[string]string{
				dialect.MySQL:    "bigint",
				dialect.Postgres: "bigint",
				dialect.SQLite:   "integer",
			}).
			Descriptor(),
	}
}

// Hook returns a new ent.Hook for generating the types.ID.
func Hook() ent.Hook {
	type setter interface {
		SetID(ID)
	}
	g := sonyflake.NewSonyflake(Config.Get())
	return func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !m.Op().Is(ent.OpCreate) {
				return n.Mutate(ctx, m)
			}

			is, ok := m.(setter)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation %T", m)
			}
			id, err := g.NextID()
			if err != nil {
				return "", fmt.Errorf("error generating id: %w", err)
			}
			is.SetID(New(id))
			return n.Mutate(ctx, m)
		})
	}
}
