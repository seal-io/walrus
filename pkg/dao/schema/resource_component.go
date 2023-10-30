package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/schema/mixin"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type ResourceComponent struct {
	ent.Schema
}

func (ResourceComponent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ID(),
		mixin.Time(),
	}
}

func (ResourceComponent) Fields() []ent.Field {
	return []ent.Field{
		object.IDField("project_id").
			Comment("ID of the project to belong.").
			NotEmpty().
			Immutable(),
		object.IDField("environment_id").
			Comment("ID of the environment to which the component belongs.").
			NotEmpty().
			Immutable(),
		object.IDField("resource_id").
			Comment("ID of the resource to which the component belongs.").
			NotEmpty().
			Immutable(),
		object.IDField("connector_id").
			Comment("ID of the connector to which the component deploys.").
			NotEmpty().
			Immutable(),
		object.IDField("composition_id").
			Comment("ID of the parent component.").
			Optional().
			Immutable(),
		object.IDField("class_id").
			Comment("ID of the parent class of the component realization.").
			Optional().
			Immutable(),
		field.String("mode").
			Comment("Mode that manages the generated component, " +
				"it is the management way of the deployer to the component, " +
				"which provides by deployer.").
			NotEmpty().
			Immutable(),
		field.String("type").
			Comment("Type of the generated component, " +
				"it is the type of the resource which the deployer observes, " +
				"which provides by deployer.").
			NotEmpty().
			Immutable(),
		field.String("name").
			Comment("Name of the generated component, " +
				"it is the real identifier of the component, " +
				"which provides by deployer.").
			NotEmpty().
			Immutable(),
		field.String("deployer_type").
			Comment("Type of deployer.").
			NotEmpty().
			Immutable(),
		field.String("shape").
			Comment("Shape of the component, it can be class or instance shape.").
			NotEmpty().
			Immutable(),
		field.JSON("status", types.ServiceResourceStatus{}).
			Comment("Status of the component.").
			Optional(),
		field.JSON("keys", &types.ServiceResourceOperationKeys{}).
			Comment("Keys of the component.").
			Optional().
			Annotations(
				entx.SkipInput(),
				entx.SkipStoringField()),
	}
}

func (ResourceComponent) Edges() []ent.Edge {
	return []ent.Edge{
		// Project 1-* ResourceComponents.
		edge.From("project", Project.Type).
			Ref("resource_components").
			Field("project_id").
			Comment("Project to which the component belongs.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.ValidateContext(intercept.WithProjectInterceptor)),
		// Environment 1-* ResourceComponents.
		edge.From("environment", Environment.Type).
			Ref("resource_components").
			Field("environment_id").
			Comment("Environment to which the component deploys.").
			Unique().
			Required().
			Immutable(),
		// Resource 1-* ResourceComponents.
		edge.From("resource", Resource.Type).
			Ref("components").
			Field("resource_id").
			Comment("Resource to which the component belongs.").
			Unique().
			Required().
			Immutable(),
		// Connector 1-* ResourceComponents.
		edge.From("connector", Connector.Type).
			Ref("resource_components").
			Field("connector_id").
			Comment("Connector to which the component deploys.").
			Unique().
			Required().
			Immutable().
			Annotations(
				entx.SkipInput()),
		// ResourceComponent (!discovered) 1-* ResourceComponents (discovered).
		edge.To("components", ResourceComponent.Type).
			Comment("Components that makes up the resource component.").
			From("composition").
			Field("composition_id").
			Unique().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade)),
		// ResourceComponent (class) 1-* ResourceComponents (instance).
		edge.To("instances", ResourceComponent.Type).
			Comment("Instances that realizes the resource component.").
			From("class").
			Field("class_id").
			Unique().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade)),
		// ResourceComponent 1-* ResourceComponents (dependency).
		edge.To("dependencies", ResourceComponent.Type).
			Comment("Dependencies that requires for the resource component.").
			Through("resource_component_relationships", ResourceComponentRelationship.Type),
	}
}

func (ResourceComponent) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.ByProject("project_id"),
	}
}
