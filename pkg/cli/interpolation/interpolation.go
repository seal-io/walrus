package interpolation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/compose-spec/compose-go/interpolation"
	"github.com/compose-spec/compose-go/template"

	"github.com/seal-io/walrus/pkg/cli/config"
)

const (
	APIGroupVariables = "variables"
)

var variableFilePattern = regexp.MustCompile(`@([^@\s]+)`)

// Interpolate interpolates the given yaml object.
func Interpolate(c config.ScopeContext, yml map[string]any, interps ...Interpolator) (map[string]any, error) {
	if len(interps) == 0 {
		interps = DefaultInterpolator(c)
	}

	var err error
	for _, v := range interps {
		yml, err = v.Func(yml)
		if err != nil {
			return nil, fmt.Errorf("failed to interpolate with %s: %w", v.Name, err)
		}
	}

	return yml, nil
}

// DefaultInterpolator returns the default interpolators.
func DefaultInterpolator(c config.ScopeContext) []Interpolator {
	return []Interpolator{
		ContextAndEnvironmentVariableInterpolator(c),
		VariableFileInterpolator(c),
	}
}

// Interpolator represents a type used for interpolating values in a YAML object.
type Interpolator struct {
	Name string
	Func InterpolateFunc
}

// InterpolateFunc is a function that interpolates values in a YAML object.
type InterpolateFunc = func(ymlObj map[string]any) (map[string]any, error)

// ContextAndEnvironmentVariableInterpolator interpolates the context and environment variable.
func ContextAndEnvironmentVariableInterpolator(c config.ScopeContext) Interpolator {
	return Interpolator{
		Name: "context-and-environment-variable-interpolator",
		Func: func(ymlObj map[string]any) (map[string]any, error) {
			return interpolation.Interpolate(ymlObj, interpolation.Options{
				LookupValue: func(key string) (string, bool) {
					switch key {
					case "Project":
						return c.Project, true
					case "Environment":
						return c.Environment, true
					default:
						return os.LookupEnv(key)
					}
				},
			})
		},
	}
}

// VariableFileInterpolator interpolates variables from files.
func VariableFileInterpolator(_ config.ScopeContext) Interpolator {
	substitute := func(tmpl string, mapping template.Mapping) (string, error) {
		return template.SubstituteWithOptions(
			tmpl,
			mapping,
			template.WithPattern(variableFilePattern),
			template.WithReplacementFunction(func(s string, m template.Mapping, c *template.Config) (string, error) {
				file := strings.TrimPrefix(s, "@")
				if !filepath.IsAbs(file) {
					dir, err := os.Getwd()
					if err != nil {
						return "", err
					}

					file = filepath.Join(dir, file)
				}

				b, err := os.ReadFile(file)
				if err != nil {
					return "", err
				}

				return string(b), nil
			}),
		)
	}

	return Interpolator{
		Name: "variable-file-interpolator",
		Func: func(ymlObj map[string]any) (map[string]any, error) {
			opts := interpolation.Options{
				Substitute: substitute,
			}

			var (
				variables = make(map[string]any)
				ok        bool
			)

			variables[APIGroupVariables], ok = ymlObj[APIGroupVariables]
			if !ok {
				return ymlObj, nil
			}

			interp, err := interpolation.Interpolate(variables, opts)
			if err != nil {
				return nil, err
			}

			ymlObj[APIGroupVariables] = interp[APIGroupVariables]
			if err != nil {
				return nil, err
			}

			return ymlObj, nil
		},
	}
}
