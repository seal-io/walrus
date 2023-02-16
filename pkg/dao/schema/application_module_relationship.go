package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/oid"
	"github.com/seal-io/seal/pkg/dao/schema/mixin"
)

type ApplicationModule struct {
	// ID of module that configure to the application.
	ModuleID string `json:"moduleID"`
	// Name of the module customized to the application.
	Name string `json:"name"`
	// Variables to configure the module.
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type ApplicationModuleRelationship struct {
	relationSchema
}

func (ApplicationModuleRelationship) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (ApplicationModuleRelationship) Annotations() []Annotation {
	return []Annotation{
		field.ID("application_id", "module_id"),
	}
}

func (ApplicationModuleRelationship) Fields() []ent.Field {
	return []ent.Field{
		oid.Field("application_id").
			Comment("ID of the application to which the relationship connects.").
			StructTag(`json:"applicationID"`).
			NotEmpty().
			Immutable(),
		field.String("module_id").
			Comment("ID of the module to which the relationship connects.").
			StructTag(`json:"moduleID"`).
			NotEmpty().
			Immutable(),
		field.String("name").
			Comment("Name of the module customized to the application.").
			NotEmpty().
			Immutable(),
		field.JSON("variables", map[string]interface{}{}).
			Comment("Variables to configure the module.").
			Optional(),
	}
}

func (ApplicationModuleRelationship) Edges() []ent.Edge {
	// NB(thxCode): entc cannot recognize camel case field name on edge with `Through`.
	return []ent.Edge{
		edge.To("application", Application.Type).
			Field("application_id").
			Comment("Applications that connect to the relationship.").
			Unique().
			Required().
			Immutable().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("module", Module.Type).
			Field("module_id").
			Comment("Modules that connect to the relationship.").
			Unique().
			Required().
			Immutable().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}),
	}
}
