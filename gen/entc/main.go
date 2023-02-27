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
