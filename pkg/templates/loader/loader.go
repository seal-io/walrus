package loader

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/types"
)

// SchemaLoader define the interface for loading schema from template.
type SchemaLoader interface {
	Load(rootDir, templateName string, mode Mode) (*types.TemplateVersionSchema, error)
}

// LoadSchemaPreferFile loads schema from template, prefer load from schema.yaml file.
func LoadSchemaPreferFile(rootDir, templateName string) (s *types.TemplateVersionSchema, err error) {
	return LoadSchema(rootDir, templateName, ModeSchemaFile)
}

// LoadOriginalSchema loads schema generate from original template.
func LoadOriginalSchema(rootDir, templateName string) (s *types.TemplateVersionSchema, err error) {
	return LoadSchema(rootDir, templateName, ModeOriginal)
}

// LoadAndMergeSchema loads the schema from the template,
// and merges it with the schemas from schema.yaml and the generated schema.
func LoadAndMergeSchema(rootDir, templateName string) (s *types.TemplateVersionSchema, err error) {
	return LoadSchema(rootDir, templateName, ModeMerge)
}

// LoadSchema loads schema from template.
func LoadSchema(rootDir, templateName string, mode Mode) (s *types.TemplateVersionSchema, err error) {
	// Terraform.
	tf := NewTerraformLoader()

	s, err = tf.Load(rootDir, templateName, mode)
	if err != nil {
		return nil, err
	}

	if s != nil {
		return s, nil
	}

	// Continue with other loaders in the future.
	return nil, fmt.Errorf("no supported schema found for template %s", templateName)
}

type Mode uint8

const (
	ModeSchemaFile Mode = iota
	ModeOriginal
	ModeMerge
)
