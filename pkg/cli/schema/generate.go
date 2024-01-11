package schema

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/seal-io/walrus/pkg/templates/loader"
)

const (
	schemaFileName = "schema.yaml"
)

func Generate(opts GenerateOption) error {
	if opts.Dir == "" {
		return fmt.Errorf("dir is empty")
	}

	_, err := os.Stat(opts.Dir)
	if err != nil {
		return err
	}

	tmplName := filepath.Base(opts.Dir)

	s, err := loader.LoadOriginalSchema(opts.Dir, tmplName)
	if err != nil {
		return err
	}

	if s == nil || s.IsEmpty() {
		return fmt.Errorf("no supported schema found for template %s", tmplName)
	}

	us, err := loader.LoadFileSchema(opts.Dir, tmplName)
	if err != nil {
		return err
	}

	b, err := FormattedOpenAPI(s, us)
	if err != nil {
		return err
	}

	// Write to file.
	f, err := os.Create(filepath.Join(opts.Dir, schemaFileName))
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", f.Name(), err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", f.Name(), err)
	}

	return nil
}
