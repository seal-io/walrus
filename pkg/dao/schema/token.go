package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
)

type Token struct {
	ent.Schema
}

func (Token) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.OwnBySubject{},
		mixin.CreateTime{},
	}
}

func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.String("kind").
			Comment("The kind of token.").
			Default(types.TokenKindAPI).
			Immutable(),
		field.String("name").
			Comment("The name of token.").
			NotEmpty().
			Immutable(),
		field.Time("expiration").
			Comment("The time of expiration, empty means forever.").
			Nillable().
			Optional().
			Immutable(),
		crypto.StringField("value").
			Comment("The value of token, store in string.").
			NotEmpty().
			Immutable().
			Sensitive(),
	}
}

func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		// Subject 1-* tokens.
		edge.From("subject", Subject.Type).
			Ref("tokens").
			Field("subjectID").
			Comment("Subject to which the token belongs.").
			Unique().
			Required().
			Immutable(),
	}
}
