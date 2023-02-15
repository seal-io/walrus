package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type Module struct {
	schema
}

func (Module) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Meta{},
		mixin.Status{},
		mixin.Time{},
	}
}

func (Module) Fields() []ent.Field {
	return []ent.Field{
		field.String("source").
			Comment("Source of the module."),
		field.String("version").
			Comment("Version of the module."),
		field.JSON("inputSchema", map[string]interface{}{}).
			Comment("Input schema of the module.").
			Optional(),
		field.JSON("outputSchema", map[string]interface{}{}).
			Comment("Output schema of the module.").
			Optional(),
	}
}
