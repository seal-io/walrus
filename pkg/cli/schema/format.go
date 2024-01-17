package schema

import (
	"bytes"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/utils/json"
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
func FormattedOpenAPI(originSchema, fileSchema *types.TemplateVersionSchema) ([]byte, error) {
	// 1. Get variables sequence from original schema.
	vs := originSchema.VariableSchema()
	oseq := openapi.GetExtOriginal(vs.Extensions).VariablesSequence

	// 2. Get merged schema.
	var (
		merged *openapi3.Schema
		err    error
	)

	es := originSchema.Expose(openapi.WalrusContextVariableName)
	if es.IsEmpty() {
		return nil, nil
	}

	if fileSchema == nil || fileSchema.IsEmpty() {
		merged = es.VariableSchema()
	} else {
		fe := fileSchema.Expose()

		merged, err = openapi.UnionSchema(es.VariableSchema(), fe.VariableSchema())
		if err != nil {
			return nil, err
		}
	}

	// 3. Generate merged sequence.
	seq := make([]string, 0)

	for _, v := range oseq {
		if v == openapi.WalrusContextVariableName {
			continue
		}

		seq = append(seq, v)
	}

	exist := sets.NewString(seq...)
	for n := range merged.Properties {
		if !exist.Has(n) {
			seq = append(seq, n)
		}
	}

	// 4. Sorted Variables to the same sequence original defined, since the properties is a map,
	// so need to generate sorted map.
	// 4.1 First convert to the json.
	propsJsonByte, err := json.Marshal(merged.Properties)
	if err != nil {
		return nil, err
	}

	// 4.2 Then properties convert to yaml.
	propMap := make(map[string]map[string]any)

	err = yaml.Unmarshal(propsJsonByte, &propMap)
	if err != nil {
		return nil, err
	}

	// 4.3 Generate key sorted properties.
	sortedProps := make(yaml.MapSlice, len(propMap))

	for i, v := range seq {
		sortedSchema := sortWithSequence(schemaSequence, propMap[v])
		sortedProps[i] = yaml.MapItem{
			Key:   v,
			Value: sortedSchema,
		}
	}

	// 5. Generate final schema.
	w := WrapOpenAPI{
		OpenAPI: es.OpenAPISchema.OpenAPI,
		Info: WrapInfo{
			Title:      es.OpenAPISchema.Info.Title,
			Extensions: es.OpenAPISchema.Info.Extensions,
		},
		Components: WrapComponents{
			Schemas: map[string]any{
				types.VariableSchemaKey: genVariable(*merged, sortedProps),
			},
		},
	}

	// 6. Convert to yaml.
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)

	err = enc.Encode(w)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func genVariable(esv openapi3.Schema, props yaml.MapSlice) yaml.MapSlice {
	sortedVariablesSchema := yaml.MapSlice{
		{
			Key:   "type",
			Value: "object",
		},
	}

	if len(esv.Required) != 0 {
		sortedVariablesSchema = append(sortedVariablesSchema, yaml.MapItem{
			Key:   "required",
			Value: esv.Required,
		})
	}

	if len(props) != 0 {
		sortedVariablesSchema = append(sortedVariablesSchema, yaml.MapItem{
			Key:   "properties",
			Value: props,
		})
	}

	extUI := openapi.GetExtUI(esv.Extensions)
	if !extUI.IsEmpty() {
		sortedVariablesSchema = append(sortedVariablesSchema, yaml.MapItem{
			Key:   openapi.ExtUIKey,
			Value: extUI,
		})
	}

	return sortedVariablesSchema
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
