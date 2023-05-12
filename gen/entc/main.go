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
	configStyle()
	configTemplateFuncs()
	configTemplate()
}

func main() {
	var err = generate()
	if err != nil {
		log.Fatalf("error generating: %v", err)
	}
}

// generate produces DAO APIs in a safer and cleaner way.
func generate() (err error) {
	// Prepare.
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

	// Generate.
	var feats = []gen.Feature{
		gen.FeatureSnapshot,
		gen.FeatureSchemaConfig,
		gen.FeatureLock,
		gen.FeatureModifier,
		gen.FeatureExecQuery,
		gen.FeatureUpsert,
		gen.FeatureVersionedMigration,
	}
	var cfg = gen.Config{
		Features: feats,
		Header:   string(header),
		Target:   newGeneratedDir,
		Schema:   "github.com/seal-io/seal/pkg/dao/schema",
		Package:  "github.com/seal-io/seal/pkg/dao/model",
	}
	err = entc.Generate(schemaDir, &cfg, entc.TemplateDir(templateDir))
	if err != nil {
		return err
	}

	// Save new generated.
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

// configStyle configures the style of generation.
func configStyle() {
	// Goimports prefix.
	imports.LocalPrefix = "github.com/seal-io/seal"
}

// configTemplateFuncs configures the functions of template generation.
func configTemplateFuncs() {
	// Override.
	gen.Funcs["camel"] = strs.CamelizeDownFirst
	gen.Funcs["snake"] = strs.Underscore
	gen.Funcs["pascal"] = strs.Camelize
	// Extend.
	gen.Funcs["getInputFields"] = getInputFields
	gen.Funcs["getInputEdges"] = getInputEdges
	gen.Funcs["getOutputFields"] = getOutputFields
	gen.Funcs["getOutputEdges"] = getOutputEdges
}

// configTemplate configures the template of generation.
func configTemplate() {
	var pkgf = func(s string) func(t *gen.Type) string {
		return func(t *gen.Type) string { return fmt.Sprintf(s, t.PackageDir()) }
	}
	// Generate io file for per model.
	gen.Templates = append(gen.Templates, gen.TypeTemplate{
		Name:   "io",
		Format: pkgf("%s_io.go"),
		ExtendPatterns: []string{
			// Combine the go templates that matches the following patterns together,
			// render and output to the file path formatted by `pkgf`.
			"io",
			"io/additional",
			"io/additional/*",
		},
	})
}
