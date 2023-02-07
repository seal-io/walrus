package oid

import (
	"context"
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/entc/integration/ent/hook"
	"entgo.io/ent/schema/field"
	"github.com/sony/sonyflake"
)

func Field() ent.Field {
	return field.Other("id", ID("")).
		Comment("The ID of resource.").
		SchemaType(map[string]string{
			dialect.MySQL:    "bigint",
			dialect.Postgres: "bigint",
			dialect.SQLite:   "integer",
		}).
		Immutable()
}

func Hook() ent.Hook {
	type setter interface {
		SetID(ID)
	}
	var g = sonyflake.NewSonyflake(sonyflake.Settings{})
	var h = func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			is, ok := m.(setter)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation %T", m)
			}
			var id, err = g.NextID()
			if err != nil {
				return "", fmt.Errorf("error generating id: %w", err)
			}
			is.SetID(New(id))
			return n.Mutate(ctx, m)
		})
	}
	return hook.On(h, ent.OpCreate)
}
