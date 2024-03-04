package loader

import (
	"github.com/seal-io/walrus/pkg/dao/types"
)

type SchemaGroup struct {
	Schema   *types.TemplateVersionSchema `json:"schema"`
	UISchema *types.UISchema              `json:"uiSchema"`
}

// SchemaLoader define the interface for loading schema from template.
type SchemaLoader interface {
	Load(rootDir, templateName string) (*SchemaGroup, error)
}

// LoadSchema loads schema from template.
func LoadSchema(rootDir, templateName string) (s *SchemaGroup, err error) {
	// Terraform.
	tf := NewTerraformLoader()

	return tf.Load(rootDir, templateName)
}
