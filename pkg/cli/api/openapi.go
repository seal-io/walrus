package api

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/seal-io/seal/utils/strs"
)

// OpenAPI Extensions.
const (
	// ExtOperationName define the extension key to set the CLI operation name.
	ExtOperationName = "x-cli-operation-name"

	// ExtSchemaTypeName define the extension key to set the CLI operation params schema type.
	ExtSchemaTypeName = "x-cli-schema-type"

	// ExtIgnore define the extension key to ignore an operation.
	ExtIgnore = "x-cli-ignore"
)

const (
	// JsonMediaType is support request body media type.
	JsonMediaType = "application/json"
)

// LoadOpenAPI load OpenAPI schema from response body and generate API.
func LoadOpenAPI(resp *http.Response) (*API, error) {
	data, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	t := openapi3.T{}

	err = t.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	if !isSupportOpenAPI(t.OpenAPI) {
		return nil, fmt.Errorf("unsupported OpenAPI version")
	}

	var (
		api        = &API{}
		operations []Operation
	)

	if t.Paths == nil {
		return api, nil
	}

	basePath, err := t.Servers.BasePath()
	if err != nil {
		return nil, err
	}

	for subPath, pathItem := range t.Paths {
		if isIgnore(pathItem.Extensions) {
			continue
		}

		opPath := path.Join(basePath, subPath)

		for method, operation := range pathItem.Operations() {
			if isIgnore(operation.Extensions) {
				continue
			}

			op := toOperation(method, opPath, pathItem, operation, t.Components)
			operations = append(operations, op)
		}
	}

	api.Short = t.Info.Title
	api.Long = t.Info.Description
	api.Operations = operations

	return api, nil
}

// toOperation generate operation from OpenAPI operation schema.
func toOperation(
	method, basePath string,
	pathItem *openapi3.PathItem,
	op *openapi3.Operation,
	comps *openapi3.Components,
) Operation {
	var (
		allParams = make([]*openapi3.Parameter, len(op.Parameters))
		seen      = make(map[string]struct{})
	)

	for i, p := range op.Parameters {
		allParams[i] = p.Value
		seen[p.Value.Name] = struct{}{}
	}

	for _, p := range pathItem.Parameters {
		if _, ok := seen[p.Value.Name]; !ok {
			allParams = append(allParams, p.Value)
		}
	}

	var (
		pathParams   []*Param
		queryParams  []*Param
		headerParams []*Param
	)

	for i, p := range allParams {
		param := toParam(allParams[i])

		switch p.In {
		case "path":
			pathParams = append(pathParams, param)
		case "query":
			queryParams = append(queryParams, param)
		case "header":
			headerParams = append(headerParams, param)
		}
	}

	md := strings.ToUpper(method)
	mediaType, bodyParams := toBodyParams(op.RequestBody, comps)

	var (
		group = ""
		tag   = ""
	)

	if len(op.Tags) > 0 {
		tag = op.Tags[0]
		group = strings.ToLower(strs.Singularize(tag))
	}

	dep := ""
	if op.Deprecated {
		dep = "Deprecated"
	}

	name := deGroupedName(tag, op.OperationID)
	if override := getExt(op.Extensions, ExtOperationName, ""); override != "" {
		name = override
	}

	return Operation{
		Name:          name,
		Group:         group,
		Short:         op.Summary,
		Long:          op.Description,
		Method:        md,
		URITemplate:   basePath,
		PathParams:    pathParams,
		QueryParams:   queryParams,
		HeaderParams:  headerParams,
		BodyParams:    bodyParams,
		BodyMediaType: mediaType,
		Deprecated:    dep,
	}
}

// toParam generate param from OpenAPI parameter.
func toParam(p *openapi3.Parameter) *Param {
	var (
		typ = "string"
		des string
		def interface{}
	)

	if p.Schema != nil && p.Schema.Value != nil {
		typ, des, def = schemaType(p.Schema.Value)
	}

	param := &Param{
		Type:        typ,
		Name:        p.Name,
		Description: des,
		Style:       openapi3.SerializationSimple,
		Default:     def,
	}

	if p.Style != "" {
		param.Style = p.Style
	}

	if p.Explode != nil {
		param.Explode = *p.Explode
	}

	return param
}

// toBodyParams generate body params from OpenAPI request body.
func toBodyParams(bodyRef *openapi3.RequestBodyRef, comps *openapi3.Components) (string, *BodyParams) {
	if bodyRef == nil || bodyRef.Value == nil {
		return "", nil
	}

	mt := bodyRef.Value.GetMediaType(JsonMediaType)
	if mt == nil || mt.Schema == nil || mt.Schema.Value == nil {
		return "", nil
	}

	propToBodyParams := func(s *openapi3.Schema, comps *openapi3.Components) []*BodyParam {
		var params []*BodyParam

		for n := range s.Properties {
			ps := propSchema(s.Properties[n], comps)
			if ps == nil {
				continue
			}

			typ, des, def := schemaType(ps)
			bp := &BodyParam{
				Name:        n,
				Type:        typ,
				Description: des,
				Default:     def,
			}

			params = append(params, bp)
		}

		return params
	}

	bps := BodyParams{}

	// Request body support array and object.
	switch mt.Schema.Value.Type {
	case openapi3.TypeArray:
		it := mt.Schema.Value.Items
		if it == nil {
			return "", nil
		}

		s := propSchema(it, comps)
		if s == nil {
			return "", nil
		}

		bps.Type = openapi3.TypeArray
		bps.Params = propToBodyParams(s, comps)

	case openapi3.TypeObject:
		s := propSchema(mt.Schema, comps)
		if s == nil {
			return "", nil
		}

		bps.Type = openapi3.TypeObject
		bps.Params = propToBodyParams(s, comps)
	}

	return JsonMediaType, &bps
}

// schemaType get schema type, description, default from OpenAPI schema.
func schemaType(s *openapi3.Schema) (string, string, interface{}) {
	var (
		typ     = s.Type
		extType string
		des     = s.Description
		def     = s.Default
	)

	if len(s.Extensions) != 0 {
		tp, ok := s.Extensions[ExtSchemaTypeName]
		if ok {
			extType = tp.(string)
		}
	}

	switch {
	case extType != "":
		typ = extType
	case s.Type == "array":
		typ = "array[object]"

		if s.Items != nil && s.Items.Value != nil {
			typ = fmt.Sprintf("array[%s]", s.Items.Value.Type)
			des = s.Items.Value.Description
			def = s.Items.Value.Default
		}
	case s.Type == "object":
		// Only id.
		if _, ok := s.Properties["id"]; ok && len(s.Properties) == 1 {
			typ = "objectID"
		}
	}

	return typ, des, def
}

// propSchema get schema from schema reference.
func propSchema(prop *openapi3.SchemaRef, comps *openapi3.Components) *openapi3.Schema {
	if prop.Value != nil {
		return prop.Value
	}

	if comps == nil || len(comps.Schemas) == 0 {
		return nil
	}

	arr := strings.Split(prop.Ref, "/")
	if len(arr) < 1 {
		return nil
	}

	name := arr[len(arr)-1]

	sr, ok := comps.Schemas[name]
	if !ok {
		return nil
	}

	return sr.Value
}

// isSupportOpenAPI check OpenAPI version.
func isSupportOpenAPI(v string) bool {
	vs := strings.Split(v, "")
	if len(vs) < 1 {
		return false
	}

	return vs[0] == "3"
}

// isIgnore check whether it include ignore extension.
func isIgnore(ext map[string]any) bool {
	return getExt(ext, ExtIgnore, false)
}

// getExt get extension by key.
func getExt[T any](v map[string]any, key string, def T) T {
	if v != nil {
		if i := v[key]; i != nil {
			if t, ok := i.(T); ok {
				return t
			}
		}
	}

	return def
}

// deGroupedName generate name without group.
func deGroupedName(group, name string) string {
	name = strings.TrimPrefix(name, strs.Pluralize(strs.Camelize(group)))
	name = strings.TrimPrefix(name, strs.Singularize(group))
	name = strings.TrimPrefix(name, ".")
	name = strs.Dasherize(name)
	name = strings.ToLower(name)

	return name
}
