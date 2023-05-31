package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

type Service struct {
	ent.Schema
}

func (Service) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.OwnByProject{},
		mixin.Time{},
	}
}

func (Service) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("environmentID", "name").
			Unique(),
	}
}

func (Service) Fields() []ent.Field {
	return []ent.Field{
		oid.Field("environmentID").
			Comment("ID of the environment to which the service deploys.").
			NotEmpty().
			Immutable(),
		field.JSON("template", types.TemplateVersionRef{}).
			Comment("Template ID and version."),
		property.ValuesField("attributes").
			Comment("Attributes to configure the template.").
			Optional(),
		field.String("name").
			Comment("Name of the service.").
			NotEmpty().
			Immutable(),
		field.JSON("status", status.Status{}).
			Comment("Status of the service.").
			Optional(),
	}
}

func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		// Environment 1-* services.
		edge.From("environment", Environment.Type).
			Ref("services").
			Field("environmentID").
			Comment("Environment to which the service belongs.").
			Unique().
			Required().
			Immutable(),
		// Project 1-* services.
		edge.From("project", Project.Type).
			Ref("services").
			Field("projectID").
			Comment("Project to which the service belongs.").
			Unique().
			Required().
			Immutable(),
		// Service 1-* service revisions.
		edge.To("revisions", ServiceRevision.Type).
			Comment("Revisions that belong to the service.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		// Service 1-* service resources.
		edge.To("resources", ServiceResource.Type).
			Comment("Resources that belong to the service.").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
