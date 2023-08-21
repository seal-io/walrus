package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
)

type Catalog struct {
	ent.Schema
}

func (Catalog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Metadata(),
		mixin.Status(),
	}
}

func (Catalog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}

func (Catalog) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").
			Comment("Type of the catalog.").
			NotEmpty().
			Immutable(),
		field.String("source").
			Comment("Source of the catalog.").
			NotEmpty().
			Immutable(),
		field.JSON("sync", &types.CatalogSync{}).
			Comment("Sync information of the catalog.").
			Optional().
			Annotations(
				entx.SkipInput()),
	}
}

func (Catalog) Edges() []ent.Edge {
	return []ent.Edge{
		// Catalog 1-* Templates.
		edge.To("templates", Template.Type).
			Comment("Templates that belong to this catalog.").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entx.SkipIO()),
	}
}
