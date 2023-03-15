package crypto

import (
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/utils/cryptox"
	"github.com/seal-io/seal/utils/vars"
)

// EncryptorConfig holds the config of the String encryption.
var EncryptorConfig = vars.NewSetOnce[cryptox.Encryptor](cryptox.Null())

// StringField returns a new ent.Field with type String.
func StringField(name string) *stringBuilder {
	return &stringBuilder{
		desc: field.String(name).
			GoType(String("")).
			SchemaType(map[string]string{
				dialect.MySQL:    "blob",
				dialect.Postgres: "bytea",
				dialect.SQLite:   "blob",
			}).
			Descriptor(),
	}
}

// BytesField returns a new ent.Field with type Bytes.
func BytesField(name string) *bytesBuilder {
	return &bytesBuilder{
		desc: field.Bytes(name).
			GoType(Bytes(nil)).
			SchemaType(map[string]string{
				dialect.MySQL:    "blob",
				dialect.Postgres: "bytea",
				dialect.SQLite:   "blob",
			}).
			Descriptor(),
	}
}

// MapField returns a new ent.Field with type Map.
func MapField[K comparable, V any](name string) *otherBuilder {
	return &otherBuilder{
		desc: field.Other(name, Map[K, V]{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "blob",
				dialect.Postgres: "bytea",
				dialect.SQLite:   "blob",
			}).
			Descriptor(),
	}
}

// SliceField returns a new ent.Field with type Slice.
func SliceField[T any](name string) *otherBuilder {
	return &otherBuilder{
		desc: field.Other(name, Slice[T]{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "blob",
				dialect.Postgres: "bytea",
				dialect.SQLite:   "blob",
			}).
			Descriptor(),
	}
}
