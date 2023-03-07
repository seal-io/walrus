package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Setting struct {
	schema
}

func (Setting) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}

func (Setting) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("The name of system setting.").
			NotEmpty(),
		field.String("value").
			Comment("The value of system setting, store in string."),
		field.Bool("hidden").
			Comment("Indicate the system setting should be hidden or not, default is visible.").
			Nillable().
			Default(false),
		field.Bool("editable").
			Comment("Indicate the system setting should be edited or not, default is readonly.").
			Nillable().
			Default(false),
		field.Bool("private").
			Comment("Indicate the system setting should be exposed or not, default is exposed.").
			Nillable().
			Default(false),
	}
}
