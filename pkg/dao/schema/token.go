package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type Token struct {
	ent.Schema
}

func (Token) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time().WithoutUpdateTime(),
	}
}

func (Token) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("subject_id").
			Comment("ID of the subject to belong.").
			NotEmpty().
			Immutable(),
		field.String("kind").
			Comment("The kind of token.").
			Default(types.TokenKindAPI).
			Immutable(),
		field.String("name").
			Comment("The name of token.").
			Unique().
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
		field.String("access_token").
			Comment("AccessToken is the token used for authentication.").
			Optional().
			Annotations(
				entx.SkipInput(),
				entx.SkipStoringField()),
	}
}

func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		// Subject 1-* Tokens.
		edge.From("subject", Subject.Type).
			Ref("tokens").
			Field("subject_id").
			Comment("Subject to which the token belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.SkipInput(),
				entx.ValidateContext(intercept.WithSubjectInterceptor)),
	}
}

func (Token) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.BySubject("subject_id"),
	}
}

func (Token) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.SkipClearingOptionalField(),
	}
}
