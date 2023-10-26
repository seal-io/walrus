package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
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
		index.Fields("project_id", "name").
			Unique().
			Annotations(
				entsql.IndexWhere("project_id IS NOT NULL")),
		index.Fields("name").
			Unique().
			Annotations(
				entsql.IndexWhere("project_id IS NULL")),
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
		object.IDField("project_id").
			Comment("ID of the project to belong, empty means for all projects.").
			Immutable().
			Optional(),
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
		// Project 1-* Catalogs.
		edge.From("project", Project.Type).
			Ref("catalogs").
			Field("project_id").
			Comment("Project to which the catalog belongs.").
			Unique().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor)),
	}
}

func (Catalog) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProjectOptional("project_id"),
	}
}
