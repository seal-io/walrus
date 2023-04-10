package property

import (
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

// PropertiesField returns a new ent.Field with type Properties.
func PropertiesField(name string) *otherBuilder {
	return &otherBuilder{
		desc: field.Other(name, Properties{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "json",
				dialect.Postgres: "jsonb",
				dialect.SQLite:   "text",
			}).
			Descriptor(),
	}
}

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

// SchemasField returns a new ent.Field with type Schemas.
func SchemasField(name string) *otherBuilder {
	return &otherBuilder{
		desc: field.Other(name, Schemas{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "json",
				dialect.Postgres: "jsonb",
				dialect.SQLite:   "text",
			}).
			Descriptor(),
	}
}
