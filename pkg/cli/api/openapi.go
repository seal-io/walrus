package api

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"

	"github.com/seal-io/walrus/pkg/apis/runtime/openapi"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/utils/strs"
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
	api.Operations = aggregateOperations(operations)

	return api, nil
}

func aggregateOperations(operations []Operation) []Operation {
	var (
		agg    = make(map[string][]Operation)
		finals = make([]Operation, 0)
	)

	for _, op := range operations {
		key := op.Group + "-" + op.Name
		agg[key] = append(agg[key], op)
	}

	for i, ops := range agg {
		if len(ops) == 1 {
			finals = append(finals, ops[0])
			continue
		}

		sorted := agg[i]
		sort.SliceStable(sorted, func(x, y int) bool {
			return len(sorted[x].PathParams) < len(sorted[y].PathParams)
		})

		var (
			tmpl  string
			first = sorted[0]
			last  = sorted[len(sorted)-1]
		)

		// Set path params from flag.
		diff := paramDiff(first.PathParams, last.PathParams, "")
		for pi := range last.PathParams {
			if slices.Contains(diff, last.PathParams[pi].Name) {
				last.PathParams[pi].DataFrom = DataFromFlag
			}
		}

		// Generate uri template.
		for z := len(sorted) - 1; z >= 0; z-- {
			df := paramDiff(first.PathParams, sorted[z].PathParams, `.`)

			switch {
			case z == 0:
				// Last one.
				tmpl += fmt.Sprintf(`{{ else }}%s{{ end }}`, sorted[z].URITemplate)
			case z == len(sorted)-1:
				// First one has the most parameters.
				tmpl += fmt.Sprintf(`{{ if and %s }}%s`, strings.Join(df, " "), sorted[z].URITemplate)
			default:
				tmpl += fmt.Sprintf(`{{ else if and %s }}%s`, strings.Join(df, " "), sorted[z].URITemplate)
			}
		}

		last.URITemplate = tmpl
		last.URIParams = diff
		last.Long = last.Short

		finals = append(finals, last)
	}

	return finals
}

func paramDiff(ps1, ps2 []*Param, prefix string) []string {
	var diff []string

	for _, p2 := range ps2 {
		var exist bool

		for _, p1 := range ps1 {
			if p1.Name == p2.Name {
				exist = true
				break
			}
		}

		if !exist {
			n := prefix + p2.Name
			diff = append(diff, n)
		}
	}

	return diff
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
		if p != nil && isIgnore(p.Extensions) {
			continue
		}

		param := toParam(allParams[i], basePath)

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

	tag := ""
	if len(op.Tags) > 0 {
		tag = op.Tags[0]
	}

	dep := ""
	if op.Deprecated {
		dep = "Deprecated"
	}

	name := deGroupedName(tag, op.OperationID)
	if override := getExt(op.Extensions, openapi.ExtCliOperationName, ""); override != "" {
		name = override
	}

	var formats []string
	if efs := getExt(op.Extensions, openapi.ExtCliOutputFormat, ""); efs != "" {
		formats = strings.Split(efs, ",")
	}

	return Operation{
		Name:          name,
		Group:         tag,
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
		Formats:       formats,
	}
}

// toParam generate param from OpenAPI parameter.
func toParam(p *openapi3.Parameter, uri string) *Param {
	var (
		typ = "string"
		des string
		def any
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

	switch p.In {
	case "path":
		param.DataFrom = DataFromArg
	case "query":
		param.DataFrom = DataFromFlag
	case "header":
		param.DataFrom = DataFromFlag
	}

	for _, ijf := range config.InjectFields {
		// Set context param if it's in inject fields and not the id in the uri.
		if ijf == p.Name && !strings.HasSuffix(uri, fmt.Sprintf("{%s}", ijf)) {
			param.DataFrom = DataFromContextAndArg
		}
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

			if isIgnore(ps.Extensions) {
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
func schemaType(s *openapi3.Schema) (string, string, any) {
	var (
		typ     = s.Type
		extType string
		des     = s.Description
		def     = s.Default
	)

	if len(s.Extensions) != 0 {
		tp, ok := s.Extensions[openapi.ExtCliSchemaTypeName]
		if ok {
			extType = tp.(string)
		}
	}

	switch {
	case extType != "":
		typ = extType
	case s.Type == openapi3.TypeArray:
		typ = ValueTypeArrayObject

		if s.Items != nil && s.Items.Value != nil {
			typ = fmt.Sprintf("array[%s]", s.Items.Value.Type)
			des = s.Items.Value.Description
			def = s.Items.Value.Default
		}
	case s.Type == openapi3.TypeObject:
		// Only id.
		if _, ok := s.Properties["id"]; ok && len(s.Properties) == 1 {
			typ = ValueTypeObjectID
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
	return getExt(ext, openapi.ExtCliIgnore, false)
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
