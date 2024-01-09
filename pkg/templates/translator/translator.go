package translator

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
)

// Translator translates between original template language and go types with openapi schema.
type Translator interface {
	// SchemaOfOriginalType generates openAPI schema from original type.
	SchemaOfOriginalType(typ any, opts Options) *openapi3.Schema
	// ToGoTypeValues converts values to go types.
	ToGoTypeValues(values map[string]json.RawMessage, schema openapi3.Schema) (map[string]any, error)
}

type Options struct {
	Name          string
	DefaultValue  any
	DefaultObject any
	Description   string
	Sensitive     bool
	Order         int
	TypeExpress   any
	// IgnoreWidget used for skipping add widget for any type in sub schema.
	IgnoreWidget bool
}

// SchemaOfType generates openAPI schema from original type.
func SchemaOfType(typ any, opts Options) (schema openapi3.Schema) {
	var s *openapi3.Schema

	// Terraform.
	tf := NewTerraformTranslator()

	s = tf.SchemaOfOriginalType(typ, opts)
	if s != nil {
		return *s
	}

	// Continue with other translator in the future.

	// No translator found.
	log.Warnf("no supported translator found for type %v, name %s", typ, opts.Name)

	// Default unknown type.
	s = openapi3.NewSchema().
		WithDefault(opts.DefaultValue)
	s.Title = opts.Name
	s.Description = opts.Description
	s.WriteOnly = opts.Sensitive

	return *s
}

// ToGoTypeValues converts values to go types.
func ToGoTypeValues(values map[string]json.RawMessage, schema openapi3.Schema) (r map[string]any, err error) {
	// Terraform.
	tf := NewTerraformTranslator()

	r, err = tf.ToGoTypeValues(values, schema)
	if err != nil {
		return nil, err
	}

	if r != nil {
		return r, nil
	}

	// Continue with other translator in the future.

	// No translator found.
	return nil, fmt.Errorf("no supported translator found for convert %v to go type", values)
}
