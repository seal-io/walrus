package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/schema/mixin"
	"github.com/seal-io/seal/pkg/dao/types"
)

type ModuleVersion struct {
	ent.Schema
}

func (ModuleVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID{},
		mixin.Time{},
	}
}

func (ModuleVersion) Fields() []ent.Field {
	return []ent.Field{
		field.String("moduleID").
			Comment("ID of the module.").
			NotEmpty().
			Immutable(),
		field.String("version").
			Comment("Module version.").
			NotEmpty().
			Immutable(),
		// This is the normalized terraform module source that can be directly applied to terraform configuration.
		// For example, when we store multiple versions of a module in a mono repo,
		//   Module.Source = "github.com/foo/bar"
		//   ModuleVersion.Source = "github.com/foo/bar/1.0.0"
		field.String("source").
			Comment("Module version source.").
			NotEmpty().
			Immutable(),
		field.JSON("schema", &types.ModuleSchema{}).
			Comment("Schema of the module.").
			Default(&types.ModuleSchema{}),
	}
}

func (ModuleVersion) Edges() []ent.Edge {
	return []ent.Edge{
		// Module 1-* module versions.
		edge.From("module", Module.Type).
			Ref("versions").
			Field("moduleID").
			Unique().
			Required().
			Immutable(),
	}
}
