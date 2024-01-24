package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
)

type Setting struct {
	ent.Schema
}

func (Setting) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time(),
	}
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
			NotEmpty().
			Immutable().
			Annotations(
				entx.SkipInput(
					entx.WithCreate(),
					entx.WithUpdate())),
		crypto.StringField("value").
			Comment("The value of system setting, store in string.").
			StructTag(`json:"value,omitempty,cli-table-column"`),
		field.Bool("hidden").
			Comment("Indicate the system setting should be hidden or not, default is visible.").
			Nillable().
			Optional().
			Default(false).
			Annotations(
				entx.SkipInput()),
		field.Bool("editable").
			Comment("Indicate the system setting should be edited or not, default is readonly.").
			Nillable().
			Optional().
			Default(false).
			Annotations(
				entx.SkipInput()),
		field.Bool("sensitive").
			Comment("Indicate the system setting should be sanitized or not before exposing, default is not.").
			Nillable().
			Optional().
			Default(false).
			Annotations(
				entx.SkipInput()),
		field.Bool("private").
			Comment("Indicate the system setting should be exposed or not, default is exposed.").
			Nillable().
			Optional().
			Default(false).
			Annotations(
				entx.SkipIO()),
		field.Bool("configured").
			Comment("Configured indicates the setting whether to be configured.").
			Optional().
			Annotations(
				entx.SkipInput(),
				entx.SkipStoringField()),
	}
}

func (Setting) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.SkipClearingOptionalField(),
	}
}
