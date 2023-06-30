package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/io"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
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
		oid.Field("serviceID").
			Comment("ID of the service to which the revision belongs.").
			NotEmpty().
			Immutable(),
		oid.Field("environmentID").
			Comment("ID of the environment to which the service deploys.").
			NotEmpty().
			Immutable(),
		field.String("templateID").
			Comment("ID of the template.").
			NotEmpty().
			Immutable(),
		field.String("templateVersion").
			Comment("Version of the template.").
			NotEmpty(),
		property.ValuesField("attributes").
			Comment("Attributes to configure the template.").
			Optional(),
		crypto.MapField[string, string]("secrets").
			Comment("Secrets of the revision.").
			Default(crypto.Map[string, string]{}),
		crypto.MapField[string, string]("variables").
			Comment("Variables of the revision.").
			Default(crypto.Map[string, string]{}),
		field.String("inputPlan").
			Comment("Input plan of the revision.").
			Sensitive(),
		field.String("output").
			Comment("Output of the revision.").
			Sensitive(),
		field.String("deployerType").
			Comment("Type of deployer.").
			Default(types.DeployerTypeTF),
		field.Int("duration").
			Comment("Duration in seconds of the revision deploying.").
			Default(0),
		field.JSON("previousRequiredProviders", []types.ProviderRequirement{}).
			Comment("Previous provider requirement of the revision.").
			Default([]types.ProviderRequirement{}),
		field.Strings("tags").
			Comment("Tags of the revision.").
			Default([]string{}),
	}
}

func (ServiceRevision) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* service revisions.
		edge.From("project", Project.Type).
			Ref("serviceRevisions").
			Field("projectID").
			Comment("Project to which the revision belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				io.DisableInput()),
		// Environment 1-* service revisions.
		edge.From("environment", Environment.Type).
			Ref("serviceRevisions").
			Field("environmentID").
			Comment("Environment to which the revision deploys.").
			Unique().
			Required().
			Immutable().
			Annotations(
				io.DisableInput()),
		// Service 1-* service revisions.
		edge.From("service", Service.Type).
			Ref("revisions").
			Field("serviceID").
			Comment("Service to which the revision belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				io.DisableInput()),
	}
}
