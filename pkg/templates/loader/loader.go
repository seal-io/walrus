package loader

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/types"
)

// SchemaLoader define the interface for loading schema from template.
type SchemaLoader interface {
	Load(rootDir, templateName, templateVersion string) (*types.Schema, error)
}

// LoadSchema loads schema from template.
func LoadSchema(rootDir, templateName, templateVersion string) (s *types.Schema, err error) {
	// Terraform.
	tf := NewTerraformLoader()

	s, err = tf.Load(rootDir, templateName, templateVersion)
	if err != nil {
		return nil, err
	}

	if s != nil {
		return s, nil
	}

	// Continue with other loaders in the future.
	return nil, fmt.Errorf("no supported schema found for template %s version %s", templateName, templateVersion)
}
