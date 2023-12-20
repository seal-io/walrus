package interpolation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/compose-spec/compose-go/interpolation"
	"github.com/compose-spec/compose-go/template"
)

var substitutionFileDecode = regexp.MustCompile(`\${file\(["']?([^"'\)]+)["']?\)}`)

var (
	delimiter            = "\\$"
	substitutionNamed    = "[_a-z][_a-z0-9]*"
	substitutionBraced   = "[_a-z][_a-z0-9]*(?::?[-+?](.*))?"
	substitutionVariable = "{var\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?}"
	substitutionResource = "{res\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?\\.[_a-z][_a-z0-9]*}"
	substitutionService  = "{svc\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?\\.[_a-z][_a-z0-9]*}"

	substitutionFile = `{file\(["']?([a-zA-Z]:)?(/|\\)?([A-Za-z0-9]+(/|\\))*[A-Za-z0-9]+(\.[a-zA-Z]{1,5}){0,1}["']?\)}`
)

var (
	groupEscaped  = "escaped"
	groupNamed    = "named"
	groupBraced   = "braced"
	groupVariable = "variable"
	groupResource = "resource"
	groupService  = "service"
	groupFile     = "file"
	groupInvalid  = "invalid"
)

// patternString is adapted from:
// https://github.com/compose-spec/compose-go/blob/81e1e9036e66e9afcdbecffea6470ff7edcffef8/template/template.go#L38
var patternString = fmt.Sprintf(
	"%s(?i:(?P<%s>%s)|(?P<%s>%s)|(?P<%s>%s)|(?P<%s>%s)|(?P<%s>%s)|(?P<%s>%s)|{(?:(?P<%s>%s)}|(?P<%s>)))",

	delimiter,
	groupEscaped, delimiter,
	groupNamed, substitutionNamed,
	groupVariable, substitutionVariable,
	groupResource, substitutionResource,
	groupService, substitutionService,
	groupFile, substitutionFile,
	groupBraced, substitutionBraced,
	groupInvalid,
)

var defaultPattern = regexp.MustCompile(patternString)

// Interpolate interpolates the given yaml object.
func Interpolate(yml map[string]any, data map[string]string, validateParameterAllSet bool, interps ...Interpolator) (
	map[string]any, error,
) {
	if len(interps) == 0 {
		interps = DefaultInterpolator()
	}

	var err error
	for _, v := range interps {
		yml, err = v.Func(yml, data, validateParameterAllSet)
		if err != nil {
			return nil, fmt.Errorf("failed to interpolate with %s: %w", v.Name, err)
		}
	}

	return yml, nil
}

// DefaultInterpolator returns the default interpolates.
func DefaultInterpolator() []Interpolator {
	return []Interpolator{
		ContextAndEnvironmentVariableInterpolator(),
	}
}

// Interpolator represents a type used for interpolating values in a YAML object.
type Interpolator struct {
	Name string
	Func InterpolateFunc
}

// InterpolateFunc is a function that interpolates values in a YAML object.
type InterpolateFunc = func(ymlObj map[string]any, data map[string]string, validateParameterAllSet bool) (
	map[string]any, error)

// ContextAndEnvironmentVariableInterpolator interpolates the context and environment variable.
func ContextAndEnvironmentVariableInterpolator() Interpolator {
	substitute := func(tmpl string, mapping template.Mapping) (string, error) {
		return template.SubstituteWithOptions(
			tmpl,
			mapping,
			template.WithPattern(defaultPattern),
			template.WithReplacementFunction(func(s string, m template.Mapping, c *template.Config) (string, error) {
				switch {
				case strings.HasPrefix(s, "${var.") ||
					strings.HasPrefix(s, "${svc.") ||
					strings.HasPrefix(s, "${res."):
					return s, nil
				case strings.HasPrefix(s, "${file("):
					matches := substitutionFileDecode.FindAllStringSubmatch(s, -1)
					if len(matches) == 0 || len(matches[0]) < 2 {
						return "", fmt.Errorf("invalid file substitution: %s", s)
					}

					file := matches[0][1]
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
				default:
					return template.DefaultReplacementFunc(s, m, c)
				}
			}),
		)
	}

	return Interpolator{
		Name: "context-and-environment-variable-interpolator",
		Func: func(ymlObj map[string]any, data map[string]string, validateParametersSet bool) (map[string]any, error) {
			return interpolation.Interpolate(
				ymlObj,
				interpolation.Options{
					LookupValue: func(key string) (string, bool) {
						switch key {
						case "Project":
							return data["project"], true
						case "Environment":
							return data["environment"], true
						default:
							val, exist := os.LookupEnv(key)
							if validateParametersSet {
								if !exist || val == "" {
									panic(fmt.Errorf("the value of  %s is not set", key))
								}
							}
							return val, exist
						}
					},
					Substitute: substitute,
				})
		},
	}
}
