package loader

import (
	"github.com/seal-io/walrus/pkg/dao/types"
)

// SchemaLoader define the interface for loading schema from template.
type SchemaLoader interface {
	Load(rootDir, templateName string, mode Mode) (*types.TemplateVersionSchema, error)
}

// LoadFileSchema loads schema from schema.yaml file.
func LoadFileSchema(rootDir, templateName string) (s *types.TemplateVersionSchema, err error) {
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

	return tf.Load(rootDir, templateName, mode)
}

type Mode uint8

const (
	ModeSchemaFile Mode = iota
	ModeOriginal
	ModeMerge
)
