package schema

import (
	"bytes"
	"sort"

	"gopkg.in/yaml.v2"

	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/utils/json"
)

const (
	defaultGroup = "Basic"
)

var schemaSequence = []string{
	"title",
	"type",
	"description",
	"default",
	"format",
	"enum",
	"example",
	"externalDocs",
	"not",
	"oneOf",
	"anyOf",
	"allOf",
	// Array related.
	"uniqueItems",
	"minItems",
	"maxItems",
	"items",
	// Properties.
	"nullable",
	"readOnly",
	"writeOnly",
	"allowEmptyValue",
	"deprecated",
	// Number related.
	"exclusiveMinimum",
	"exclusiveMaximum",
	"minimum",
	"maximum",
	"multipleOf",
	// String related.
	"minLength",
	"maxLength",
	"pattern",
	// Object related.
	"required",
	"properties",
	"minProperties",
	"maxProperties",
	"additionalProperties",
	"discriminator",
}

// WrapOpenAPI is a wrapper of openapi3.T, used to generated formatted yaml.
type WrapOpenAPI struct {
	OpenAPI    string         `json:"openapi" yaml:"openapi"`
	Info       WrapInfo       `json:"info" yaml:"info"`
	Components WrapComponents `json:"components,omitempty" yaml:"components,omitempty"`
}

// WrapInfo is a wrapper of openapi3.Info and removed field version, used to generated formatted yaml.
type WrapInfo struct {
	Title      string         `json:"title" yaml:"title"`
	Extensions map[string]any `json:",inline" yaml:",inline"`
}

// WrapComponents is a wrapper of openapi3.Components, used to generated formatted yaml.
type WrapComponents struct {
	Schemas map[string]any `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

// FormattedOpenAPI generates formatted openapi yaml.
func FormattedOpenAPI(s types.Schema) ([]byte, error) {
	// Get variables sequence.
	vs := s.VariableSchema()
	seq := openapi.GetOriginalVariablesSequence(vs.Extensions)

	// Expose Variables.
	es := s.Expose()
	injectVariablesExtGroup(&es)

	// Sorted Variables to the same sequence original defined, since the properties is a map,
	// so need to generate sorted map.
	// 1. First convert to the json.
	propsJsonByte, err := json.Marshal(es.VariableSchema().Properties)
	if err != nil {
		return nil, err
	}

	// 2. Then properties convert to yaml.
	propMap := make(map[string]map[string]any)

	err = yaml.Unmarshal(propsJsonByte, &propMap)
	if err != nil {
		return nil, err
	}

	// 3. Generate key sorted properties.
	sortedProps := make(yaml.MapSlice, len(propMap))

	for i, v := range seq {
		sortedSchema := sortWithSequence(schemaSequence, propMap[v])
		sortedProps[i] = yaml.MapItem{
			Key:   v,
			Value: sortedSchema,
		}
	}

	// 4. Generate key sorted variables.
	esv := *es.OpenAPISchema.Components.Schemas["variables"].Value
	sortedVariablesSchema := yaml.MapSlice{
		{
			Key:   "required",
			Value: esv.Required,
		},
		{
			Key:   "type",
			Value: "object",
		},
		{
			Key:   "properties",
			Value: sortedProps,
		},
	}

	w := WrapOpenAPI{
		OpenAPI: es.OpenAPISchema.OpenAPI,
		Info: WrapInfo{
			Title:      es.OpenAPISchema.Info.Title,
			Extensions: es.OpenAPISchema.Info.Extensions,
		},
		Components: WrapComponents{
			Schemas: map[string]any{
				"variables": sortedVariablesSchema,
			},
		},
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)

	err = enc.Encode(w)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// injectVariablesExtGroup default injects group extension for variables.
func injectVariablesExtGroup(s *types.UISchema) {
	// Inject group extension.
	vs := s.VariableSchema()
	for n, v := range vs.Properties {
		if v.Value == nil || v.Value.IsEmpty() {
			continue
		}

		if gp := openapi.GetUIGroup(v.Value.Extensions); gp == "" {
			vs.Properties[n].Value.Extensions = openapi.NewExt(vs.Properties[n].Value.Extensions).
				SetUIGroup(defaultGroup).
				Export()
		}
	}
}

func sortWithSequence(seq []string, m map[string]any) yaml.MapSlice {
	seqKeys := make(map[string]int, len(seq))
	for i, v := range seq {
		seqKeys[v] = i
	}

	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		si, ok1 := seqKeys[keys[i]]
		sj, ok2 := seqKeys[keys[j]]

		switch {
		case ok1 && ok2:
			return si < sj
		case ok1:
			return true
		case ok2:
			return false
		default:
			return keys[i] < keys[j]
		}
	})

	result := make(yaml.MapSlice, len(keys))
	for i, v := range keys {
		result[i] = yaml.MapItem{
			Key:   v,
			Value: m[v],
		}
	}

	return result
}
