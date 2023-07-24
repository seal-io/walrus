package view

import (
	"bytes"
	"fmt"
	"path/filepath"

	"entgo.io/ent/entc/gen"

	"github.com/seal-io/seal/pkg/dao/entx/extension"
)

func Generate(g *gen.Graph) (err error) {
	var generated extension.Assets

	templates = loadTemplate()

	// Generate.
	for _, n := range g.Nodes {
		for _, tmpl := range Templates {
			// Generate.
			buf := bytes.NewBuffer(nil)
			if err = templates.ExecuteTemplate(buf, tmpl.Name, n); err != nil {
				return fmt.Errorf("genearte %q: %w", tmpl.Name, err)
			}

			generated.Add(
				filepath.Join(g.Config.Target, tmpl.Format(n)),
				buf.Bytes())
		}
	}

	for _, tmpl := range GraphTemplates {
		if tmpl.Skip != nil && tmpl.Skip(g) {
			continue
		}

		buf := bytes.NewBuffer(nil)
		if err = templates.ExecuteTemplate(buf, tmpl.Name, g); err != nil {
			return fmt.Errorf("genearte %q: %w", tmpl.Name, err)
		}

		generated.Add(
			filepath.Join(g.Config.Target, tmpl.Format),
			buf.Bytes())
	}

	// Write.
	if err = generated.Write(); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	// Format.
	if err = generated.Format(); err != nil {
		return fmt.Errorf("format: %w", err)
	}

	return
}
