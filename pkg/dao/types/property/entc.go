package property

import (
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

// ValuesField returns a new ent.Field with type Values.
func ValuesField(name string) *otherBuilder {
	return &otherBuilder{
		desc: field.Other(name, Values{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "json",
				dialect.Postgres: "jsonb",
				dialect.SQLite:   "text",
			}).
			Descriptor(),
	}
}
