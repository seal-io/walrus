package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Token struct {
	schema
}

func (Token) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("casdoorTokenName").
			Unique(),
	}
}

func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.String("casdoorTokenName").
			Comment("The token name of casdoor").
			NotEmpty(),
		field.String("casdoorTokenOwner").
			Comment("The token owner of casdoor").
			NotEmpty(),
		field.String("name").
			Comment("The name of token.").
			NotEmpty(),
		field.Int("expiration").
			Comment("Expiration in seconds.").
			Nillable().
			Optional(),
	}
}
