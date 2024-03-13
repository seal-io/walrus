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

// Loader is an interface that defines methods for loading objects from various sources.
// Objects are loaded either from a list of file paths or from a byte slice.
type Loader interface {
	LoadFromFiles(filePaths []string, recursive bool) (ObjectSet, error)
	LoadFromByte(b []byte) (ObjectSet, error)
}

// DefaultLoader creates a Loader with default configurations.
func DefaultLoader(sc *config.Config, validateParametersSet bool, runLabels map[string]string) Loader {
	return &ObjectLoader{
		serverContext: sc,
		groups:        GroupSequence,
		interpolator:  interpolation.DefaultInterpolator(),
		interpolateData: map[string]string{
			"project":     sc.Project,
			"environment": sc.Environment,
		},
		validateParametersSet: validateParametersSet,
		runLabels:             runLabels,
	}
}

// ObjectLoader is a type that represents an object loader.
type ObjectLoader struct {
	interpolator          []interpolation.Interpolator
	interpolateData       map[string]string
	groups                []string
	serverContext         *config.Config
	validateParametersSet bool
	runLabels             map[string]string
}

// LoadFromByte parses the provided YAML byte slice `b` and returns an `ObjectSet` and an error.
func (l *ObjectLoader) LoadFromByte(b []byte) (ObjectSet, error) {
	set := ObjectSet{}

	yml, err := loader.ParseYAML(b)
	if err != nil {
		return set, fmt.Errorf("failed to parse yaml: %w", err)
	}

	if len(l.interpolator) != 0 {
		yml, err = interpolation.Interpolate(yml, l.interpolateData, l.validateParametersSet, l.interpolator...)
		if err != nil {
			return set, fmt.Errorf("failed to interpolate: %w", err)
		}
	}

	list, err := l.toObjects(yml)
	if err != nil {
		return set, err
	}

	set.Add(list...)

	return set, nil
}

// LoadFromFiles loads object from file paths and a recursive flag indicating whether to search files recursively.
func (l *ObjectLoader) LoadFromFiles(filePaths []string, recursive bool) (ObjectSet, error) {
	var (
		allPaths []string
		set      = ObjectSet{}
	)

	for _, path := range filePaths {
		files, err := findYAMLFiles(path, recursive)
		if err != nil {
			return set, err
		}

		allPaths = append(allPaths, files...)
	}

	for _, path := range allPaths {
		s, err := l.LoadFromFile(path)
		if err != nil {
			return set, fmt.Errorf("failed to load manifest %s: %w", path, err)
		}

		set.Add(s.All()...)
	}

	return set, nil
}

// LoadFromFile reads the contents of the file at the given path and loads it as a YAML object set.
func (l *ObjectLoader) LoadFromFile(filePath string) (ObjectSet, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return ObjectSet{}, err
	}

	return l.LoadFromByte(b)
}

// toObjects extracts a list of Object from the given map.
func (l *ObjectLoader) toObjects(m map[string]any) ([]Object, error) {
	objs := make([]Object, 0)

	for _, group := range l.groups {
		v, ok := m[group]
		if !ok {
			continue
		}

		list, ok := v.([]any)
		if !ok {
			return nil, fmt.Errorf("invalid %s list: should be array", group)
		}

		for _, item := range list {
			itemMap, ok := item.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid %s: should be map", group)
			}

			obj, err := l.toObject(group, itemMap)
			if err != nil {
				return nil, err
			}

			if obj == nil {
				continue
			}

			objs = append(objs, *obj)
		}
	}

	return objs, nil
}

// toObject converts a map of object properties into an *Object instance.
func (l *ObjectLoader) toObject(group string, obj map[string]any) (*Object, error) {
	if len(obj) == 0 {
		return nil, nil
	}

	name, ok := obj["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name is empty")
	}

	project, err := getName(obj, "project")
	if err != nil {
		return nil, err
	}

	env, err := getName(obj, "environment")
	if err != nil {
		return nil, err
	}

	scope := ObjectScope{
		Project:     project,
		Environment: env,
	}

	switch group {
	case GroupResources:
		if scope.Project == "" {
			scope.Project = l.serverContext.Project
		}

		if scope.Environment == "" {
			scope.Environment = l.serverContext.Environment
		}

		if scope.Project == "" || scope.Environment == "" {
			return nil, fmt.Errorf("need to indicate the project name and environment name for resource %s",
				name)
		}

		obj["project"] = IDName{
			Name: scope.Project,
		}
		obj["environment"] = IDName{
			Name: scope.Environment,
		}

		obj["runLabels"] = l.runLabels

	default:
		return nil, fmt.Errorf("supported group %s", group)
	}

	return &Object{
		ObjectScope: scope,
		IDName: IDName{
			Name: name,
		},
		Group: group,
		Value: obj,
	}, nil
}

var yamlFileSuffix = []string{".yaml", ".yml"}

// findYAMLFiles finds all YAML files within a specified directory.
// If recursive is set to true, it will recursively search for YAML files in subdirectories as well.
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

func getName(obj map[string]any, key string) (string, error) {
	o, ok := obj[key]
	if !ok {
		return "", nil // Key not found is not an error.
	}

	m, ok := o.(map[string]any)
	if !ok {
		return "", fmt.Errorf("%s is invalid", key)
	}

	value, _ := m["name"].(string)

	return value, nil
}
