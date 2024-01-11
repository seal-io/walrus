package types

import (
	"context"
	"errors"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"

	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
)

const (
	VariableSchemaKey = "variables"
	OutputSchemaKey   = "outputs"
)

// Schema specifies the openAPI schema with variables and outputs.
type Schema struct {
	OpenAPISchema *openapi3.T `json:"openAPISchema"`
}

// Validate reports if the schema is valid.
func (s *Schema) Validate() error {
	if s.OpenAPISchema == nil {
		return nil
	}

	// workaround: inject paths and version since kin-openapi/openapi3 need it.
	s.OpenAPISchema.Paths = openapi3.Paths{}
	if s.OpenAPISchema.Info != nil && s.OpenAPISchema.Info.Version == "" {
		s.OpenAPISchema.Info.Version = "v0.0.0"
	}

	if err := s.OpenAPISchema.Validate(
		context.Background(),
		openapi3.DisableSchemaDefaultsValidation(),
	); err != nil {
		return err
	}

	return nil
}

func (s *Schema) IsEmpty() bool {
	return s.OpenAPISchema == nil || s.OpenAPISchema.Components == nil || len(s.OpenAPISchema.Components.Schemas) == 0
}

// Expose returns the UI schema of the schema.
func (s *Schema) Expose(skipProps ...string) UISchema {
	vs := s.VariableSchema()
	if vs == nil {
		return UISchema{}
	}

	for _, v := range skipProps {
		delete(vs.Properties, v)
	}

	// In order to prevent the remove ext affect the original schema, serialize and deserialize to copy the schema.
	b, err := json.Marshal(vs)
	if err != nil {
		log.Warnf("error marshal variable schema while expost: %v", err)
		return UISchema{}
	}

	var cps openapi3.Schema

	err = json.Unmarshal(b, &cps)
	if err != nil {
		log.Warnf("error unmarshal variable schema while expost: %v", err)
		return UISchema{}
	}

	return UISchema{
		OpenAPISchema: &openapi3.T{
			OpenAPI: s.OpenAPISchema.OpenAPI,
			Info:    s.OpenAPISchema.Info,
			Components: &openapi3.Components{
				Schemas: map[string]*openapi3.SchemaRef{
					VariableSchemaKey: {
						Value: openapi.RemoveExtOriginal(&cps),
					},
				},
			},
		},
	}
}

// VariableSchema returns the variables' schema.
func (s *Schema) VariableSchema() *openapi3.Schema {
	if s.OpenAPISchema == nil ||
		s.OpenAPISchema.Components == nil ||
		s.OpenAPISchema.Components.Schemas == nil ||
		s.OpenAPISchema.Components.Schemas[VariableSchemaKey] == nil ||
		s.OpenAPISchema.Components.Schemas[VariableSchemaKey].Value == nil {
		return nil
	}

	return s.OpenAPISchema.Components.Schemas[VariableSchemaKey].Value
}

func (s *Schema) SetVariableSchema(v *openapi3.Schema) {
	s.ensureInit()
	s.OpenAPISchema.Components.Schemas[VariableSchemaKey].Value = v
}

func (s *Schema) RemoveVariableContext() {
	if s.IsEmpty() {
		return
	}

	variableSchema := openapi.RemoveVariableContext(s.VariableSchema())
	s.SetVariableSchema(variableSchema)
}

func (s *Schema) SetOutputSchema(v *openapi3.Schema) {
	s.ensureInit()
	s.OpenAPISchema.Components.Schemas[OutputSchemaKey].Value = v
}

func (s *Schema) ensureInit() {
	if s.OpenAPISchema == nil {
		s.OpenAPISchema = &openapi3.T{}
	}

	if s.OpenAPISchema.Components == nil {
		s.OpenAPISchema.Components = &openapi3.Components{}
	}

	if s.OpenAPISchema.Components.Schemas == nil {
		s.OpenAPISchema.Components.Schemas = openapi3.Schemas{}
	}
}

// OutputSchema returns the outputs' schema.
func (s *Schema) OutputSchema() *openapi3.Schema {
	if s.OpenAPISchema == nil ||
		s.OpenAPISchema.Components == nil ||
		s.OpenAPISchema.Components.Schemas == nil ||
		s.OpenAPISchema.Components.Schemas[OutputSchemaKey] == nil ||
		s.OpenAPISchema.Components.Schemas[OutputSchemaKey].Value == nil {
		return nil
	}

	return s.OpenAPISchema.Components.Schemas[OutputSchemaKey].Value
}

// Intersect sets variables & outputs schema of s to intersection of s and s2.
func (s *Schema) Intersect(s2 *Schema) {
	if s2.OpenAPISchema == nil {
		return
	}

	variableSchema := openapi.IntersectSchema(s.VariableSchema(), s2.VariableSchema())
	s.SetVariableSchema(variableSchema)
	outputSchema := openapi.IntersectSchema(s.OutputSchema(), s2.OutputSchema())
	s.SetOutputSchema(outputSchema)
}

// UISchema include the UI schema that users can customize.
type UISchema Schema

// IsEmpty reports if the schema is empty.
func (s UISchema) IsEmpty() bool {
	return s.OpenAPISchema == nil ||
		s.OpenAPISchema.Components == nil ||
		len(s.OpenAPISchema.Components.Schemas) == 0 ||
		s.OpenAPISchema.Components.Schemas[VariableSchemaKey] == nil ||
		s.OpenAPISchema.Components.Schemas[VariableSchemaKey].Value == nil
}

// VariableSchema returns the variables' schema.
func (s *UISchema) VariableSchema() *openapi3.Schema {
	if s.IsEmpty() {
		return nil
	}

	return s.OpenAPISchema.Components.Schemas[VariableSchemaKey].Value
}

// SetVariableSchema sets the variables' schema.
func (s *UISchema) SetVariableSchema(v *openapi3.Schema) {
	if s.OpenAPISchema == nil {
		return
	}

	s.OpenAPISchema.Components.Schemas[VariableSchemaKey].Value = v
}

// Validate reports if the ui schema is valid.
func (s *UISchema) Validate() error {
	w := Schema{
		OpenAPISchema: s.OpenAPISchema,
	}

	return w.Validate()
}

// TemplateVersionSchema include the internal template variables schema and template data.
type TemplateVersionSchema struct {
	Schema `json:",inline"`

	// TemplateVersionSchemaData include the data of this template version.
	TemplateVersionSchemaData `json:",inline"`
}

// Validate reports if the schema is valid.
func (s *TemplateVersionSchema) Validate() error {
	if err := s.Schema.Validate(); err != nil {
		return err
	}

	providerExist := len(s.RequiredProviders) != 0

	var (
		variablesExist bool
		outputsExist   bool
	)

	switch {
	case s.OpenAPISchema == nil || s.OpenAPISchema.Components == nil:
		variablesExist = false
		outputsExist = false
	default:
		variablesExist = s.OpenAPISchema.Components.Schemas["variables"] != nil
		outputsExist = s.OpenAPISchema.Components.Schemas["outputs"] != nil
	}

	if !providerExist && !variablesExist && !outputsExist {
		return errors.New("invalid schema: at least one of requiredProviders, variables, outputs must be specified")
	}

	return nil
}

// TemplateVersionSchemaData include the data of this template version.
type TemplateVersionSchemaData struct {
	// Readme specifies the readme of this template.
	Readme string `json:"readme"`
	// RequiredProviders specifies the required providers of this template.
	RequiredProviders []ProviderRequirement `json:"requiredProviders"`
}

// ProviderRequirement include the required provider.
type ProviderRequirement struct {
	*tfconfig.ProviderRequirement

	Name string `json:"name"`
}
