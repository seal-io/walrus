package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"golang.org/x/tools/imports"

	"github.com/seal-io/seal/utils/files"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

func init() {
	// goimports prefix
	imports.LocalPrefix = "github.com/seal-io/seal"
}

func main() {
	var err = generate()
	if err != nil {
		log.Fatalf("error generating: %v", err)
	}
}

// generate produces DAO APIs in a safer and cleaner way.
func generate() (err error) {
	// prepare
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %w", err)
	}
	var targetDir = filepath.Join(workingDir, "/pkg/dao/model")
	var schemaDir = filepath.Join(workingDir, "/pkg/dao/schema")
	var templateDir = filepath.Join(workingDir, "/pkg/dao/template")
	generatingDir := files.TempDir("seal-dao-generated-*")
	defer func() {
		_ = os.RemoveAll(generatingDir)
	}()
	var newGeneratedDir = filepath.Join(generatingDir, "/new")
	var oldGeneratedDir = filepath.Join(generatingDir, "/old")

	header, err := os.ReadFile(filepath.Join(workingDir, "/hack/boilerplate/go.txt"))
	if err != nil {
		return err
	}

	// generate
	var cfg = gen.Config{
		Features: gen.AllFeatures,
		Header:   string(header),
		Hooks:    []gen.Hook{tagFields("json")},
		Target:   newGeneratedDir,
		Schema:   "github.com/seal-io/seal/pkg/dao/schema",
		Package:  "github.com/seal-io/seal/pkg/dao/model",
	}
	err = entc.Generate(schemaDir, &cfg, entc.TemplateDir(templateDir))
	if err != nil {
		return err
	}

	// save new generated
	err = os.Rename(targetDir, oldGeneratedDir)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			return fmt.Errorf("error cleaning stale generated files: %w", err)
		}
	}
	defer func() {
		if err != nil {
			_ = os.Rename(oldGeneratedDir, targetDir)
		}
	}()
	err = os.Rename(newGeneratedDir, targetDir)
	if err != nil {
		return fmt.Errorf("error move new generated files to %s: %w", targetDir, err)
	}

	return
}

// tagFields tags all fields defined in the schema with the given struct-tag if the field has not been tagged.
func tagFields(def string) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, n := range g.Nodes {
				for _, f := range n.Fields {
					var name = f.Name
					var defaultTag = fmt.Sprintf(`%s:"%s,omitempty"`, def, name)
					// respect the customized struct tag.
					if f.StructTag != defaultTag {
						continue
					}
					// otherwise configure omitempty tag if match the following conditions.
					name = strs.CamelizeDownFirst(name)
					var tag = fmt.Sprintf(`%s:"%s"`, def, name)
					if f.Optional || f.Default || // creation
						f.UpdateDefault || // modification
						f.Nillable { // storing
						tag = fmt.Sprintf(`%s:"%s,omitempty"`, def, name)
					}
					f.StructTag = tag
				}
			}
			return next.Generate(g)
		})
	}
}
