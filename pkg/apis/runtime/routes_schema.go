package runtime

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin/binding"
	"k8s.io/apimachinery/pkg/util/sets"

	cliapi "github.com/seal-io/seal/pkg/cli/api"
	"github.com/seal-io/seal/utils/slice"
	"github.com/seal-io/seal/utils/strs"
)

func OpenAPI(i *openapi3.Info) openapi3.T {
	sc := *spec
	if i != nil {
		sc.Info = i
	}

	return sc
}

var spec = &openapi3.T{
	OpenAPI: "3.0.3",
	Info: &openapi3.Info{
		Title:       "Restful APIs",
		Description: "API to manage resources",
		Version:     "dev",
	},
	Components: &openapi3.Components{
		Responses: getSchemaHTTPResponses(),
		Schemas:   make(map[string]*openapi3.SchemaRef),
	},
}

func getSchemaResponseHTTPs() []int {
	return []int{
		http.StatusNoContent,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusNotAcceptable,
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusGone,
		http.StatusUnprocessableEntity,
		http.StatusPreconditionRequired,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
	}
}

func getSchemaHTTPResponses() openapi3.Responses {
	httpc := getSchemaResponseHTTPs()
	resps := openapi3.Responses{}

	for i := range httpc {
		c := strconv.Itoa(httpc[i])
		t := http.StatusText(httpc[i])
		resps[c] = &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription(t).
				WithContent(
					openapi3.NewContentWithJSONSchema(
						&openapi3.Schema{
							Properties: map[string]*openapi3.SchemaRef{
								"status":     {Value: openapi3.NewIntegerSchema().WithDefault(c)},
								"statusText": {Value: openapi3.NewStringSchema().WithDefault(`"` + t + `"`)},
								"message":    {Value: openapi3.NewStringSchema()},
							},
							Required: []string{"status", "statusText"},
						},
					),
				),
		}
	}

	return resps
}

func schemeRoute(resource, handle, method, path string, ip *InputProfile, op *OutputProfile) {
	switch method {
	case http.MethodPost:
		toSchemaPath(path).Post = toSchemaOperation(resource, handle, method, path, ip, op)
	case http.MethodDelete:
		toSchemaPath(path).Delete = toSchemaOperation(resource, handle, method, path, ip, op)
	case http.MethodPut:
		toSchemaPath(path).Put = toSchemaOperation(resource, handle, method, path, ip, op)
	case http.MethodGet:
		toSchemaPath(path).Get = toSchemaOperation(resource, handle, method, path, ip, op)
	}
}

func toSchemaPath(path string) *openapi3.PathItem {
	paths := strings.Split(path, "/")
	for i := range paths {
		if paths[i] == "" || paths[i][0] != ':' {
			continue
		}
		paths[i] = "{" + paths[i][1:] + "}"
	}
	path = strings.Join(paths, "/")

	if spec.Paths == nil {
		spec.Paths = make(openapi3.Paths)
	}

	if _, ok := spec.Paths[path]; !ok {
		spec.Paths[path] = &openapi3.PathItem{}
	}

	return spec.Paths[path]
}

func toSchemaOperation(
	resource,
	handle,
	method,
	path string,
	ip *InputProfile,
	op *OutputProfile,
) *openapi3.Operation {
	r := getRoute(method, path)

	o := &openapi3.Operation{
		Tags:        []string{strs.Camelize(resource)},
		OperationID: handle,
		Parameters:  toSchemaParameters(r, ip),
		RequestBody: toSchemaRequestBody(r, ip),
		Summary:     toSchemaSummary(resource, handle, method),
		Extensions:  toSchemaExtension(resource, handle, path),
	}

	resp := getSchemaHTTPResponses()

	for k := range resp {
		resp[k].Ref = "#/components/responses/" + k
	}

	for k, v := range toSchemaResponses(r, op) {
		resp[k] = v
	}

	o.Responses = resp

	return o
}

func toSchemaParameters(r route, ip *InputProfile) []*openapi3.ParameterRef {
	if ip == nil {
		return nil
	}

	props := ip.Flat(ProfileCategoryHeader, ProfileCategoryUri, ProfileCategoryQuery)

	var (
		params         []*openapi3.ParameterRef
		queryParamsSet = sets.NewString()
	)

	for i := 0; i < len(props); i++ {
		var in string

		switch props[i].Category {
		default:
			continue
		case ProfileCategoryHeader:
			in = "header"
		case ProfileCategoryUri:
			in = "path"

			if !r.pathParams.Has(props[i].Name) {
				continue
			}
		case ProfileCategoryQuery:
			in = "query"
			name := props[i].Name

			if r.pathParams.Has(name) {
				continue
			}

			queryParamsSet.Insert(name)
		}
		param := &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				In:       in,
				Name:     props[i].Name,
				Required: props[i].Required,
				Schema:   toSchemaSchema(ProfileCategoryQuery, &props[i]),
			},
		}
		params = append(params, param)
	}

	// TODO: temporary workaround inject project query parameter.
	if injectProjectQueryPath.Has(r.path) {
		for k := range injectProjectQueryParameter {
			if !queryParamsSet.Has(k) {
				params = append(params, injectProjectQueryParameter[k])
			}
		}
	}

	return params
}

func toSchemaRequestBody(r route, ip *InputProfile) *openapi3.RequestBodyRef {
	if r.method == http.MethodGet || ip == nil {
		return nil
	}

	categoryContentTypes := map[ProfileCategory]string{
		ProfileCategoryForm: binding.MIMEMultipartPOSTForm,
		ProfileCategoryJson: binding.MIMEJSON,
	}
	if r.method != http.MethodPost {
		delete(categoryContentTypes, ProfileCategoryForm)
	}

	content := make(map[string]*openapi3.MediaType, len(categoryContentTypes))

	for category, contentType := range categoryContentTypes {
		props := ip.Filter(category)
		if len(props) == 0 {
			continue
		}

		if _, exist := content[contentType]; !exist {
			schema := basicSchemas[ip.TypeDescriptor]
			if schema == nil {
				schema = openapi3.NewObjectSchema()
			}
			content[contentType] = &openapi3.MediaType{
				Schema: schema.NewRef(),
			}
		}

		if len(props) == 1 && ip.Type == ProfileTypeArray {
			schema := toSchemaSchema(category, &props[0])
			if schema != nil {
				content[contentType] = &openapi3.MediaType{
					Schema: schema,
				}
			}

			continue
		}

		content[contentType].Schema.Value.Properties = openapi3.Schemas{}

		for i := 0; i < len(props); i++ {
			if r.pathParams.Has(props[i].Name) {
				continue
			}

			ps := toSchemaProperty(category, &props[i])
			if ps == nil {
				continue
			}

			content[contentType].Schema.Value.Properties[props[i].Name] = ps
			if props[i].Required {
				content[contentType].Schema.Value.Required = append(
					content[contentType].Schema.Value.Required,
					props[i].Name,
				)
			}
		}
	}

	if len(content) == 0 {
		return nil
	}
	req := &openapi3.RequestBodyRef{}
	req.Value = &openapi3.RequestBody{
		Content:  content,
		Required: true,
	}

	return req
}

func toSchemaResponses(r route, op *OutputProfile) map[string]*openapi3.ResponseRef {
	c := http.StatusOK
	if r.method == http.MethodPost {
		c = http.StatusCreated
	}

	resp := toSchemaResponse(r, op)
	if resp == nil {
		if r.method == http.MethodPost {
			c = http.StatusNoContent
		}
		resp = openapi3.NewResponse().WithDescription(http.StatusText(c))
	}

	return map[string]*openapi3.ResponseRef{
		strconv.Itoa(c): {
			Value: resp,
		},
	}
}

func toSchemaResponse(_ route, op *OutputProfile) *openapi3.Response {
	if op == nil {
		return nil
	}

	categoryContentTypes := map[ProfileCategory]string{
		ProfileCategoryJson: binding.MIMEJSON,
	}

	content := make(map[string]*openapi3.MediaType, 1)

	for category, contentType := range categoryContentTypes {
		props := op.Filter(category)

		if _, exist := content[contentType]; !exist {
			schema := basicSchemas[op.TypeDescriptor]
			if schema == nil {
				schema = openapi3.NewObjectSchema()
			}
			content[contentType] = &openapi3.MediaType{
				Schema: schema.NewRef(),
			}
		}

		if len(props) == 1 && (op.Type == ProfileTypeArray || op.TypeDescriptor == "render.Render") {
			content[contentType] = &openapi3.MediaType{
				Schema: toSchemaSchema(category, &props[0]),
			}

			continue
		}

		content[contentType].Schema.Value.Properties = openapi3.Schemas{}

		for i := 0; i < len(props); i++ {
			if props[i].Required {
				content[contentType].Schema.Value.Required = append(
					content[contentType].Schema.Value.Required,
					props[i].Name,
				)
			}

			ps := toSchemaProperty(category, &props[i])
			if ps != nil {
				content[contentType].Schema.Value.Properties[props[i].Name] = ps
			}
		}
	}

	if op.Page {
		for c := range content {
			media := content[c]
			itemSchema := openapi3.Schema{
				Type:  openapi3.TypeArray,
				Items: media.Schema,
			}
			s := openapi3.Schema{
				Required: []string{
					"items",
					"pagination",
				},
				Properties: map[string]*openapi3.SchemaRef{
					"items": itemSchema.NewRef(),
					"pagination": {
						Value: &openapi3.Schema{
							Properties: map[string]*openapi3.SchemaRef{
								"page":      openapi3.NewInt32Schema().NewRef(),
								"perPage":   openapi3.NewInt32Schema().NewRef(),
								"total":     openapi3.NewInt32Schema().NewRef(),
								"totalPage": openapi3.NewInt32Schema().NewRef(),
								"partial":   openapi3.NewBoolSchema().NewRef(),
								"group":     openapi3.NewBoolSchema().NewRef(),
								"nextPage":  openapi3.NewInt32Schema().NewRef(),
							},
							Required: []string{"page", "perPage", "total", "totalPage", "partial"},
						},
					},
				},
			}

			media.Schema = s.NewRef()
			content[c] = media
		}
	}

	if len(content) == 0 {
		return nil
	}

	if op.Type == ProfileTypeBasic && op.TypeDescriptor == "render.Render" {
		content["application/octet-stream"] = content[binding.MIMEJSON]
		delete(content, binding.MIMEJSON)
	}
	resp := openapi3.NewResponse().
		WithContent(content)

	return resp
}

func toSchemaProperty(category string, prop *ProfileProperty) *openapi3.SchemaRef {
	if prop == nil || prop.Name == "" {
		return nil
	}

	return toSchemaSchema(category, prop)
}

func toSchemaSummary(resource, handle, method string) string {
	var summary, subresource, handleName string

	_, handleName, _ = strings.Cut(handle, ".")
	resource = strs.Decamelize(resource, true)

	switch {
	case handleName == createPrefix:
		summary = fmt.Sprintf("Create %s", strs.SingularizeWithArticle(resource))
	case handleName == updatePrefix:
		summary = fmt.Sprintf("Update %s", strs.SingularizeWithArticle(resource))
	case handleName == deletePrefix:
		summary = fmt.Sprintf("Delete %s", strs.SingularizeWithArticle(resource))
	case handleName == getPrefix:
		summary = fmt.Sprintf("Get %s by ID", strs.Singularize(resource))
	case handleName == collectionGetPrefix:
		summary = fmt.Sprintf("Get %s", resource)
	case handleName == collectionCreatePrefix:
		summary = fmt.Sprintf("Create %s", resource)
	case handleName == collectionUpdatePrefix:
		summary = fmt.Sprintf("Update %s", resource)
	case handleName == collectionDeletePrefix:
		summary = fmt.Sprintf("Delete %s", resource)
	case handleName == collectionStreamPrefix:
		summary = fmt.Sprintf("Stream %s", resource)
	case strings.HasPrefix(handleName, getPrefix):
		subresource = strings.TrimPrefix(handleName, getPrefix)
		summary = fmt.Sprintf("Get %s", strs.Pluralize(strs.Decamelize(subresource, true)))
	case strings.HasPrefix(handleName, createPrefix):
		subresource = strings.TrimPrefix(handleName, createPrefix)
		summary = fmt.Sprintf("Create %s", strs.SingularizeWithArticle(subresource))
	case strings.HasPrefix(handleName, updatePrefix):
		subresource = strings.TrimPrefix(handleName, updatePrefix)
		summary = fmt.Sprintf("Update %s", strs.Pluralize(strs.Decamelize(subresource, true)))
	case strings.HasPrefix(handleName, deletePrefix):
		subresource = strings.TrimPrefix(handleName, deletePrefix)
		summary = fmt.Sprintf("Delete %s", strs.Pluralize(strs.Decamelize(subresource, true)))
	case strings.HasPrefix(handleName, streamPrefix):
		subresource = strings.TrimPrefix(handleName, streamPrefix)
		summary = fmt.Sprintf("Stream %s", strs.Pluralize(strs.Decamelize(subresource, true)))
	case strings.HasPrefix(handleName, collectionGetPrefix):
		subresource = strings.TrimPrefix(handleName, collectionGetPrefix)
		summary = fmt.Sprintf("Get %s", strs.Pluralize(strs.Decamelize(subresource, true)))
	case strings.HasPrefix(handleName, collectionCreatePrefix):
		subresource = strings.TrimPrefix(handleName, collectionCreatePrefix)
		summary = fmt.Sprintf("Create %s", strs.Pluralize(strs.Decamelize(subresource, true)))
	case strings.HasPrefix(handleName, collectionUpdatePrefix):
		subresource = strings.TrimPrefix(handleName, collectionUpdatePrefix)
		summary = fmt.Sprintf("Update %s", strs.Pluralize(strs.Decamelize(subresource, true)))
	case strings.HasPrefix(handleName, collectionDeletePrefix):
		subresource = strings.TrimPrefix(handleName, collectionDeletePrefix)
		summary = fmt.Sprintf("Delete %s", strs.Pluralize(strs.Decamelize(subresource, true)))
	case strings.HasPrefix(handleName, collectionStreamPrefix):
		subresource = strings.TrimPrefix(handleName, collectionStreamPrefix)
		summary = fmt.Sprintf("Stream %s", strs.Pluralize(strs.Decamelize(subresource, true)))
	// Conversion over configuration for route handlers:
	// - GET method expects the handle name to be subresource/link name
	// - Other methods expect the handle name to be the descriptive action name.
	case strings.HasPrefix(handleName, routePrefix) && method == http.MethodGet:
		subresource = strings.TrimPrefix(handleName, routePrefix)
		summary = fmt.Sprintf("Get %s", strs.Pluralize(strs.Decamelize(subresource, true)))
	case strings.HasPrefix(handleName, routePrefix) && method != http.MethodGet:
		subresource = strings.TrimPrefix(handleName, routePrefix)
		summary = strs.Capitalize(strs.Decamelize(subresource, true))
	case strings.HasPrefix(handleName, collectionRoutePrefix) && method == http.MethodGet:
		subresource = strings.TrimPrefix(handleName, collectionRoutePrefix)
		summary = fmt.Sprintf("Get %s", strs.Pluralize(strs.Decamelize(subresource, true)))
	case strings.HasPrefix(handleName, collectionRoutePrefix) && method != http.MethodGet:
		subresource = strings.TrimPrefix(handleName, collectionRoutePrefix)
		summary = strs.Capitalize(strs.Decamelize(subresource, true))
	}

	return summary
}

var (
	cliIgnoreResources = []string{
		"subjects",
		"tokens",
		"subjectRoles",
		"dashboards",
		"templateCompletions",
		"costs",
	}
	cliIgnorePaths = []string{
		"/service-revisions/:id/terraform-states",
		"/connectors/:id/repositories",
		"/connectors/:id/repository-branches",
		"/perspectives/_/field-values",
		"/perspectives/_/fields",
		"/services/_/graph",
		"/service-resources/:id/keys",
		"/service-resources/_/graph",
	}
	cliJsonYamlOutputFormatPaths = []string{
		"/service-revisions/:id/diff-latest",
		"/service-revisions/:id/diff-previous",
	}
)

func toSchemaExtension(resource, handle, path string) map[string]any {
	var (
		ext              = make(map[string]any)
		_, handleName, _ = strings.Cut(handle, ".")
	)

	switch {
	case slice.ContainsAll(cliIgnoreResources, resource):
		ext[cliapi.ExtCliIgnore] = true
	case slice.ContainsAll(cliIgnorePaths, path):
		ext[cliapi.ExtCliIgnore] = true
	case handleName == collectionGetPrefix:
		ext[cliapi.ExtCliOperationName] = "list"
	case handleName == collectionCreatePrefix:
		ext[cliapi.ExtCliIgnore] = true
	case handleName == collectionUpdatePrefix:
		ext[cliapi.ExtCliIgnore] = true
	case handleName == collectionDeletePrefix:
		ext[cliapi.ExtCliIgnore] = true
	case handleName == collectionStreamPrefix:
		ext[cliapi.ExtCliIgnore] = true
	case strings.HasPrefix(handleName, streamPrefix):
		ext[cliapi.ExtCliIgnore] = true
	case strings.HasPrefix(handleName, collectionGetPrefix):
		subresource := strings.TrimPrefix(handleName, collectionGetPrefix)
		ext[cliapi.ExtCliOperationName] = fmt.Sprintf("list-%s", strs.Pluralize(strs.Dasherize(subresource)))
	case strings.HasPrefix(handleName, collectionCreatePrefix):
		ext[cliapi.ExtCliIgnore] = true
	case strings.HasPrefix(handleName, collectionUpdatePrefix):
		ext[cliapi.ExtCliIgnore] = true
	case strings.HasPrefix(handleName, collectionDeletePrefix):
		ext[cliapi.ExtCliIgnore] = true
	case strings.HasPrefix(handleName, collectionStreamPrefix):
		ext[cliapi.ExtCliIgnore] = true
	case strings.HasPrefix(handleName, routePrefix):
		ext[cliapi.ExtCliOperationName] = strs.Dasherize(strings.TrimPrefix(handleName, routePrefix))
	case strings.HasPrefix(handleName, collectionRoutePrefix):
		ext[cliapi.ExtCliIgnore] = true
	}

	if slice.ContainsAll(cliJsonYamlOutputFormatPaths, path) {
		ext[cliapi.ExtCliOutputFormat] = "json,yaml"
	}

	return ext
}

var basicSchemas = map[string]*openapi3.Schema{
	"bool":                 openapi3.NewBoolSchema(),
	"string":               openapi3.NewStringSchema(),
	"int":                  openapi3.NewIntegerSchema(),
	"int8":                 openapi3.NewInt32Schema(),
	"int16":                openapi3.NewInt32Schema(),
	"int32":                openapi3.NewInt32Schema(),
	"uint":                 openapi3.NewInt32Schema(),
	"uint8":                openapi3.NewInt32Schema(),
	"uint16":               openapi3.NewInt32Schema(),
	"uint32":               openapi3.NewInt32Schema(),
	"int64":                openapi3.NewInt64Schema(),
	"uint64":               openapi3.NewInt64Schema(),
	"float32":              openapi3.NewFloat64Schema(),
	"float64":              openapi3.NewFloat64Schema(),
	"time.Time":            openapi3.NewDateTimeSchema(),
	"multipart.FileHeader": {Type: "string", Format: "binary"},
	"[]byte":               openapi3.NewBytesSchema(),
	"uuid.NullUUID":        openapi3.NewUUIDSchema(),
	"uuid.UUID":            openapi3.NewUUIDSchema(),
	"render.Render":        {Type: "string", Format: "binary"},
	"json.RawMessage":      openapi3.NewObjectSchema(),
}

var (
	stringToStringSchema = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			AdditionalProperties: openapi3.AdditionalProperties{
				Schema: &openapi3.SchemaRef{
					Value: openapi3.NewStringSchema(),
				},
			},
			Extensions: map[string]any{
				cliapi.ExtCliSchemaTypeName: "map[string]string",
			},
		},
	}
	stringToIntSchema = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			AdditionalProperties: openapi3.AdditionalProperties{
				Schema: &openapi3.SchemaRef{
					Value: openapi3.NewIntegerSchema(),
				},
			},
			Extensions: map[string]any{
				cliapi.ExtCliSchemaTypeName: "map[string]int",
			},
		},
	}
	stringToInt32Schema = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			AdditionalProperties: openapi3.AdditionalProperties{
				Schema: &openapi3.SchemaRef{
					Value: openapi3.NewInt32Schema(),
				},
			},
			Extensions: map[string]any{
				cliapi.ExtCliSchemaTypeName: "map[string]int32",
			},
		},
	}
	stringToInt64Schema = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			AdditionalProperties: openapi3.AdditionalProperties{
				Schema: &openapi3.SchemaRef{
					Value: openapi3.NewInt64Schema(),
				},
			},
			Extensions: map[string]any{
				cliapi.ExtCliSchemaTypeName: "map[string]int64",
			},
		},
	}
)

func toSchemaSchema(category string, prop *ProfileProperty) *openapi3.SchemaRef {
	if prop == nil {
		return nil
	}

	switch prop.Type {
	default:
		return openapi3.NewObjectSchema().NewRef()
	case ProfileTypeBasic:
		return basicSchemas[prop.TypeDescriptor].NewRef()
	case ProfileTypeArray:
		schema, exist := basicSchemas[prop.TypeDescriptor]
		if exist {
			schema = openapi3.NewArraySchema().WithItems(schema)

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				schema.MinLength = items
				schema.MaxLength = &items
			}

			return schema.NewRef()
		}
	case ProfileTypeObject:
	}

	switch prop.TypeDescriptor {
	case "object":
		schema := openapi3.NewObjectSchema()
		if prop.Type == ProfileTypeArray {
			schema = openapi3.NewArraySchema().WithItems(schema)

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				schema.MinLength = items
				schema.MaxLength = &items
			}
		}

		return schema.NewRef()
	case "array":
		schema := openapi3.NewObjectSchema().NewRef()
		if len(prop.Properties) == 1 {
			schema = toSchemaSchema(category, &prop.Properties[0])
		}

		if prop.Type == ProfileTypeArray {
			s := openapi3.Schema{
				Type:  openapi3.TypeArray,
				Items: schema,
			}

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				s.MinLength = items
				s.MaxLength = &items
			}

			return s.NewRef()
		}
	case "map[string]string":
		schema := stringToStringSchema.Value
		if prop.Type == ProfileTypeArray {
			schema = openapi3.NewArraySchema().WithItems(schema)

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				schema.MinLength = items
				schema.MaxLength = &items
			}
		}

		return schema.NewRef()
	case "map[string]int":
		schema := stringToIntSchema.Value
		if prop.Type == ProfileTypeArray {
			schema = openapi3.NewArraySchema().WithItems(schema)

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				schema.MinLength = items
				schema.MaxLength = &items
			}
		}

		return schema.NewRef()
	case "map[string]int32":
		schema := stringToInt32Schema.Value
		if prop.Type == ProfileTypeArray {
			schema = openapi3.NewArraySchema().WithItems(schema)

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				schema.MinLength = items
				schema.MaxLength = &items
			}
		}

		return schema.NewRef()
	case "map[string]int64":
		schema := stringToInt64Schema.Value
		if prop.Type == ProfileTypeArray {
			schema = openapi3.NewArraySchema().WithItems(schema)

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				schema.MinLength = items
				schema.MaxLength = &items
			}
		}

		return schema.NewRef()
	}

	schemaID := prop.TypeDescriptor
	if category != "" {
		schemaID += "." + category
	}

	if prop.TypeRefer {
		schema := openapi3.NewSchemaRef("#/components/schemas/"+schemaID, openapi3.NewSchema())
		if prop.Type == ProfileTypeArray {
			schema = &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:  openapi3.TypeArray,
					Items: schema,
				},
			}

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				schema.Value.MinLength = items
				schema.Value.MaxLength = &items
			}
		}

		return schema
	}

	namedSchema := &openapi3.Schema{
		Properties: map[string]*openapi3.SchemaRef{},
		Type:       openapi3.TypeObject,
	}

	for i := 0; i < len(prop.Properties); i++ {
		if prop.Properties[i].Required {
			namedSchema.Required = append(namedSchema.Required, prop.Properties[i].Name)
		}

		ps := toSchemaProperty(category, &prop.Properties[i])
		if ps == nil {
			continue
		}

		namedSchema.Properties[prop.Properties[i].Name] = ps
	}

	spec.Components.Schemas[schemaID] = &openapi3.SchemaRef{
		Value: namedSchema,
	}

	schema := openapi3.NewSchemaRef("#/components/schemas/"+schemaID, namedSchema)
	if prop.Type == ProfileTypeArray {
		schema = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:  openapi3.TypeArray,
				Items: schema,
			},
		}

		if prop.TypeArrayLength != 0 {
			items := uint64(prop.TypeArrayLength)
			schema.Value.MinLength = items
			schema.Value.MaxLength = &items
		}
	}

	return schema
}

func getRoute(method, path string) route {
	pathParams := sets.Set[string]{}

	for _, sg := range strings.Split(path, "/") {
		if sg == "" || sg[0] != ':' {
			continue
		}

		pathParams.Insert(sg[1:])
	}

	return route{
		method:     method,
		path:       path,
		pathParams: pathParams,
	}
}

type route struct {
	method     string
	path       string
	pathParams sets.Set[string]
}

// injectProjectQueryPath is a temporary workaround to add project id and project name query parameter,
// since some resource need project info to check permission, so project id or project name is required,
// but now we don't include them in io.
// TODO: remove this workaround after io include them.
var injectProjectQueryPath = sets.NewString(
	"/connectors/:id",
	"/connectors/:id/apply-cost-tools",
	"/connectors/:id/sync-cost-data",
	"/environments/:id",
	"/services/:id",
	"/services/:id/clone",
	"/services/:id/upgrade",
)

var injectProjectQueryParameter = map[string]*openapi3.ParameterRef{
	"projectID": {
		Value: &openapi3.Parameter{
			In:       ProfileCategoryQuery,
			Name:     "projectID",
			Required: false,
			Schema:   basicSchemas["string"].NewRef(),
		},
	},
	"projectName": {
		Value: &openapi3.Parameter{
			In:       ProfileCategoryQuery,
			Name:     "projectName",
			Required: false,
			Schema:   basicSchemas["string"].NewRef(),
		},
	},
}
