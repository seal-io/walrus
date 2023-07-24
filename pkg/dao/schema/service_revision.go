package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/pkg/dao/types/property"
)

type ServiceRevision struct {
	ent.Schema
}

func (ServiceRevision) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time().WithoutUpdateTime(),
		mixin.OwnByProject(),
		mixin.LegacyStatus(),
	}
}

func (ServiceRevision) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("environment_id").
			Comment("ID of the environment to which the service deploys.").
			NotEmpty().
			Immutable(),
		object.IDField("service_id").
			Comment("ID of the service to which the revision belongs.").
			NotEmpty().
			Immutable(),
		field.String("template_id").
			Comment("ID of the template.").
			NotEmpty().
			Immutable(),
		field.String("template_version").
			Comment("Version of the template.").
			NotEmpty(),
		property.ValuesField("attributes").
			Comment("Attributes to configure the template.").
			Optional(),
		crypto.MapField[string, string]("variables").
			Comment("Variables of the revision.").
			Default(crypto.Map[string, string]{}),
		field.String("input_plan").
			Comment("Input plan of the revision.").
			Sensitive(),
		field.String("output").
			Comment("Output of the revision.").
			Sensitive(),
		field.String("deployer_type").
			Comment("Type of deployer.").
			Default(types.DeployerTypeTF),
		field.Int("duration").
			Comment("Duration in seconds of the revision deploying.").
			Default(0),
		field.JSON("previous_required_providers", []types.ProviderRequirement{}).
			Comment("Previous provider requirement of the revision.").
			Default([]types.ProviderRequirement{}),
		field.Strings("tags").
			Comment("Tags of the revision.").
			Default([]string{}),
	}
}

func (ServiceRevision) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* ServiceRevisions.
		edge.From("project", Project.Type).
			Ref("service_revisions").
			Field("project_id").
			Comment("Project to which the revision belongs.").
			Unique().
			Required().
			Immutable(),
		// Environment 1-* ServiceRevisions.
		edge.From("environment", Environment.Type).
			Ref("service_revisions").
			Field("environment_id").
			Comment("Environment to which the revision deploys.").
			Unique().
			Required().
			Immutable(),
		// Service 1-* ServiceRevisions.
		edge.From("service", Service.Type).
			Ref("revisions").
			Field("service_id").
			Comment("Service to which the revision belongs.").
			Unique().
			Required().
			Immutable(),
	}
}
