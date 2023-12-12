package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/compose-spec/compose-go/loader"

	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/interpolation"
	"github.com/seal-io/walrus/utils/strs"
)

const (
	APIGroupResources = "resources"
	APIGroupVariables = "variables"

	ScopeProject     = "project"
	ScopeEnvironment = "environment"
	ScopeGlobal      = "global"
)

var APIGroupCreateSequence = []string{
	APIGroupVariables,
	APIGroupResources,
}

type Object struct {
	APIGroup string              `json:"-"`
	FilePath string              `json:"-"`
	Scope    string              `json:"-"`
	Context  config.ScopeContext `json:"-"`
	Name     string              `json:"name"`
	Value    map[string]any      `json:",inline"`
}

func LoadObjects(c config.ScopeContext, paths []string, recursive bool) (map[string][]Object, error) {
	var allPaths []string

	for _, path := range paths {
		files, err := findYAMLFiles(path, recursive)
		if err != nil {
			return nil, err
		}

		allPaths = append(allPaths, files...)
	}

	manifests := make(map[string]map[string]any)

	for _, path := range allPaths {
		yml, err := LoadManifest(c, path)
		if err != nil {
			return nil, fmt.Errorf("failed to load manifest %s: %w", path, err)
		}

		manifests[path] = yml
	}

	return toObjects(c, manifests)
}

func LoadManifest(c config.ScopeContext, path string) (map[string]any, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	yml, err := loader.ParseYAML(b)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml %s: %w", path, err)
	}

	yml, err = interpolation.Interpolate(c, yml)
	if err != nil {
		return nil, fmt.Errorf("failed to interpolate file %s: %w", path, err)
	}

	return yml, nil
}

func toObject(c config.ScopeContext, path, kind string, obj map[string]any) Object {
	var (
		project, _ = obj["project"].(string)
		env, _     = obj["environment"].(string)
		name, _    = obj["name"].(string)
		scope      string
		objCtx     config.ScopeContext
	)

	var ()

	switch kind {
	case APIGroupVariables:
		// Need to indicate the scope of the variable.
		switch {
		case project != "" && env != "":
			// Indicate the project name and environment name.
			scope = ScopeEnvironment
			objCtx = config.ScopeContext{
				Project:     project,
				Environment: env,
			}
		case project != "":
			// Indicate the project name.
			scope = ScopeProject
			objCtx = config.ScopeContext{
				Project: project,
			}
		default:
			// Without project name.
			scope = ScopeGlobal
		}
	case APIGroupResources:
		scope = ScopeEnvironment

		if project == "" {
			project = c.Project
		}

		if env == "" {
			env = c.Environment
		}

		if project == "" || env == "" {
			panic(
				fmt.Sprintf("need to indicate the project name and environment name for resource %s in file %s",
					name, path))
		}
		objCtx = config.ScopeContext{
			Project:     project,
			Environment: env,
		}
	default:
		panic(fmt.Sprintf("unknown kind %s", kind))
	}

	delete(obj, "project")
	delete(obj, "environment")

	return Object{
		APIGroup: kind,
		FilePath: path,
		Scope:    scope,
		Name:     name,
		Context:  objCtx,
		Value:    obj,
	}
}

func toObjects(c config.ScopeContext, ms map[string]map[string]any) (map[string][]Object, error) {
	objm := make(map[string][]Object)

	for path, m := range ms {
		for _, kind := range APIGroupCreateSequence {
			v, ok := m[kind]
			if !ok {
				continue
			}

			objs, ok := v.([]any)
			if !ok {
				return nil, fmt.Errorf("invalid services: should be array")
			}

			for _, o := range objs {
				obj, ok := o.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("invalid service: should be map")
				}

				objm[kind] = append(objm[kind], toObject(c, path, kind, obj))
			}
		}
	}

	return objm, nil
}

var yamlFileSuffix = []string{".yaml", ".yml"}

func findYAMLFiles(root string, recursive bool) ([]string, error) {
	var yamlFiles []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if !recursive && path != root {
				return filepath.SkipDir
			}
		}

		if strs.HasSuffix(d.Name(), yamlFileSuffix...) {
			abs, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			yamlFiles = append(yamlFiles, abs)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return yamlFiles, nil
}
