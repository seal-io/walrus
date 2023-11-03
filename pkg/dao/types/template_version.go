package types

import (
	"context"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"

	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/utils/errorx"
)

// UISchema include the ui schema of template.
type UISchema struct {
	// T is the openapi schema.
	openapi3.T `json:",inline"`
}

// IsEmpty reports if the schema is empty.
func (s UISchema) IsEmpty() bool {
	return s.Components == nil ||
		len(s.Components.Schemas) == 0 ||
		s.Components.Schemas["variables"] == nil ||
		s.Components.Schemas["variables"].Value == nil
}

// Validate reports if the schema is valid.
func (s UISchema) Validate() error {
	if s.IsEmpty() {
		return nil
	}
	// workaround: inject paths and version since kin-openapi/openapi3 need it.
	s.Paths = openapi3.Paths{}
	if s.Info != nil && s.Info.Version == "" {
		s.Info.Version = "mock"
	}

	err := s.T.Validate(context.Background())
	if err != nil {
		return errorx.NewHttpError(http.StatusBadRequest, err.Error())
	}

	return nil
}

// VariableSchemas returns the variable schemas.
func (s UISchema) VariableSchemas() *openapi3.Schema {
	if s.IsEmpty() {
		return nil
	}

	return s.Components.Schemas["variables"].Value
}

// Schema include the internal template variables schema and template data.
type Schema struct {
	// TemplateVersionSchemaData include the data of this template version.
	TemplateVersionSchemaData `json:",inline"`

	// OpenAPISchema specifies the openapi schema of input variables.
	OpenAPISchema *openapi3.T `json:"openAPISchema"`
}

// IsEmpty reports if the schema is empty.
func (s Schema) IsEmpty() bool {
	return s.OpenAPISchema == nil || s.OpenAPISchema.Components == nil || len(s.OpenAPISchema.Components.Schemas) == 0
}

// Validate reports if the schema is valid.
func (s Schema) Validate() error {
	if s.OpenAPISchema == nil {
		return nil
	}
	// workaround: inject paths since kin-openapi/openapi3 need it.
	s.OpenAPISchema.Paths = openapi3.Paths{}

	err := s.OpenAPISchema.Validate(context.Background())
	if err != nil {
		return errorx.NewHttpError(http.StatusBadRequest, err.Error())
	}

	providerExist := len(s.RequiredProviders) != 0

	variablesExist := s.OpenAPISchema.Components != nil && s.OpenAPISchema.Components.Schemas["variables"] != nil

	outputsExist := s.OpenAPISchema.Components != nil && s.OpenAPISchema.Components.Schemas["outputs"] != nil

	if !providerExist && !variablesExist && !outputsExist {
		return errorx.NewHttpError(http.StatusBadRequest,
			"invalid schema: at least one of requiredProviders, variables, outputs must be specified")
	}

	return nil
}

// VariableSchemas returns the variable schemas.
func (s *Schema) VariableSchemas() *openapi3.Schema {
	if s.OpenAPISchema == nil ||
		s.OpenAPISchema.Components == nil ||
		s.OpenAPISchema.Components.Schemas == nil ||
		s.OpenAPISchema.Components.Schemas["variables"] == nil ||
		s.OpenAPISchema.Components.Schemas["variables"].Value == nil {
		return nil
	}

	return s.OpenAPISchema.Components.Schemas["variables"].Value
}

// OutputSchemas returns the output schemas.
func (s *Schema) OutputSchemas() *openapi3.Schema {
	if s.OpenAPISchema == nil ||
		s.OpenAPISchema.Components == nil ||
		s.OpenAPISchema.Components.Schemas == nil ||
		s.OpenAPISchema.Components.Schemas["outputs"] == nil ||
		s.OpenAPISchema.Components.Schemas["outputs"].Value == nil {
		return nil
	}

	return s.OpenAPISchema.Components.Schemas["outputs"].Value
}

// Expose returns the exposed schema of this internal schema.
func (s *Schema) Expose() UISchema {
	vs := s.VariableSchemas()
	if vs == nil {
		return UISchema{}
	}

	return UISchema{
		T: openapi3.T{
			OpenAPI: s.OpenAPISchema.OpenAPI,
			Info:    s.OpenAPISchema.Info,
			Components: &openapi3.Components{
				Schemas: map[string]*openapi3.SchemaRef{
					"variables": {
						Value: openapi.RemoveExtOriginal(vs),
					},
				},
			},
		},
	}
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
