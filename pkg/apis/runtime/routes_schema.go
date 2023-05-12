package runtime

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/ogen-go/ogen"
	"github.com/ogen-go/ogen/jsonschema"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/utils/strs"
)

func OpenAPI(i *ogen.Info) ogen.Spec {
	sc := *spec
	if i != nil {
		sc.SetInfo(i)
	}

	return sc
}

var spec = ogen.NewSpec().
	SetOpenAPI("3.0.3").
	SetInfo(ogen.NewInfo().
		SetTitle("Restful APIs").
		SetDescription("API to manage resources").
		SetVersion("dev")).
	AddNamedResponses(getSchemaHTTPResponses()...)

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

func getSchemaHTTPResponses() []*ogen.NamedResponse {
	httpc := getSchemaResponseHTTPs()
	resps := make([]*ogen.NamedResponse, 0, len(httpc))

	for i := range httpc {
		c := strconv.Itoa(httpc[i])
		t := http.StatusText(httpc[i])
		resps = append(resps, ogen.NewNamedResponse(
			c,
			ogen.NewResponse().
				SetDescription(t).
				SetJSONContent(ogen.NewSchema().
					AddRequiredProperties(
						ogen.Int().SetDefault([]byte(c)).
							ToProperty("status"),
						ogen.String().SetDefault([]byte(`"`+t+`"`)).
							ToProperty("statusText"),
					).
					AddOptionalProperties(
						ogen.String().ToProperty("message"),
					),
				),
		))
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

func toSchemaPath(path string) *ogen.PathItem {
	paths := strings.Split(path, "/")
	for i := range paths {
		if paths[i] == "" || paths[i][0] != ':' {
			continue
		}
		paths[i] = "{" + paths[i][1:] + "}"
	}
	path = strings.Join(paths, "/")

	if spec.Paths == nil {
		spec.Paths = make(ogen.Paths)
	}

	if _, ok := spec.Paths[path]; !ok {
		spec.Paths[path] = ogen.NewPathItem()
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
) *ogen.Operation {
	r := getRoute(method, path)

	o := ogen.NewOperation().
		SetTags([]string{strs.Camelize(resource)}).
		SetOperationID(handle).
		SetParameters(toSchemaParameters(r, ip)).
		SetRequestBody(toSchemaRequestBody(r, ip)).
		SetSummary(toSchemaSummary(resource, handle, method))
	for _, c := range getSchemaResponseHTTPs() {
		o.AddNamedResponses(spec.RefResponse(strconv.Itoa(c)))
	}

	o.AddNamedResponses(toSchemaResponses(r, op)...)

	return o
}

func toSchemaParameters(r route, ip *InputProfile) []*ogen.Parameter {
	if ip == nil {
		return nil
	}

	props := ip.Flat(ProfileCategoryHeader, ProfileCategoryUri, ProfileCategoryQuery)

	var params []*ogen.Parameter

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

			if r.pathParams.Has(props[i].Name) {
				continue
			}
		}
		param := &ogen.Parameter{
			In:       in,
			Name:     props[i].Name,
			Required: props[i].Required,
			Schema:   toSchemaSchema(ProfileCategoryQuery, &props[i]),
		}
		params = append(params, param)
	}

	return params
}

func toSchemaRequestBody(r route, ip *InputProfile) *ogen.RequestBody {
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

	content := make(map[string]ogen.Media, len(categoryContentTypes))

	for category, contentType := range categoryContentTypes {
		props := ip.Filter(category)
		if len(props) == 0 {
			continue
		}

		if _, exist := content[contentType]; !exist {
			schema := basicSchemas[ip.TypeDescriptor]
			if schema == nil {
				schema = ogen.NewSchema().
					SetType(string(jsonschema.Object))
			}
			content[contentType] = ogen.Media{
				Schema: schema,
			}
		}

		if len(props) == 1 && ip.Type == ProfileTypeArray {
			schema := toSchemaSchema(category, &props[0])
			if schema != nil {
				content[contentType] = ogen.Media{
					Schema: schema,
				}
			}

			continue
		}

		for i := 0; i < len(props); i++ {
			if r.pathParams.Has(props[i].Name) {
				continue
			}

			var add func(...*ogen.Property) *ogen.Schema
			if props[i].Required {
				add = content[contentType].Schema.AddRequiredProperties
			} else {
				add = content[contentType].Schema.AddOptionalProperties
			}

			add(toSchemaProperty(category, &props[i]))
		}
	}

	if len(content) == 0 {
		return nil
	}
	req := ogen.NewRequestBody().
		SetRequired(true).
		SetContent(content)

	return req
}

func toSchemaResponses(r route, op *OutputProfile) []*ogen.NamedResponse {
	c := http.StatusOK
	if r.method == http.MethodPost {
		c = http.StatusCreated
	}

	resp := toSchemaResponse(r, op)
	if resp == nil {
		if r.method == http.MethodPost {
			c = http.StatusNoContent
		}
		resp = ogen.NewResponse().
			SetDescription(http.StatusText(c))
	}

	return []*ogen.NamedResponse{
		ogen.NewNamedResponse(strconv.Itoa(c), resp),
	}
}

func toSchemaResponse(_ route, op *OutputProfile) *ogen.Response {
	if op == nil {
		return nil
	}

	categoryContentTypes := map[ProfileCategory]string{
		ProfileCategoryJson: binding.MIMEJSON,
	}

	content := make(map[string]ogen.Media, 1)

	for category, contentType := range categoryContentTypes {
		props := op.Filter(category)

		if _, exist := content[contentType]; !exist {
			schema := basicSchemas[op.TypeDescriptor]
			if schema == nil {
				schema = ogen.NewSchema().
					SetType(string(jsonschema.Object))
			}
			content[contentType] = ogen.Media{
				Schema: schema,
			}
		}

		if len(props) == 1 && (op.Type == ProfileTypeArray || op.TypeDescriptor == "render.Render") {
			content[contentType] = ogen.Media{
				Schema: toSchemaSchema(category, &props[0]),
			}

			continue
		}

		for i := 0; i < len(props); i++ {
			var add func(...*ogen.Property) *ogen.Schema
			if props[i].Required {
				add = content[contentType].Schema.AddRequiredProperties
			} else {
				add = content[contentType].Schema.AddOptionalProperties
			}

			add(toSchemaProperty(category, &props[i]))
		}
	}

	if op.Page {
		for c := range content {
			media := content[c]
			media.Schema = ogen.NewSchema().
				AddRequiredProperties(
					media.Schema.ToProperty("items"),
					ogen.NewSchema().
						AddRequiredProperties(
							ogen.Int().ToProperty("page"),
							ogen.Int().ToProperty("perPage"),
							ogen.Int().ToProperty("total"),
							ogen.Int().ToProperty("totalPage"),
							ogen.Bool().ToProperty("partial"),
						).
						AddOptionalProperties(
							ogen.Bool().ToProperty("group"),
							ogen.Int().ToProperty("nextPage"),
						).
						ToProperty("pagination"),
				)
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
	resp := ogen.NewResponse().
		SetContent(content)

	return resp
}

func toSchemaProperty(category string, prop *ProfileProperty) *ogen.Property {
	if prop == nil || prop.Name == "" {
		return nil
	}

	return toSchemaSchema(category, prop).ToProperty(prop.Name)
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

var basicSchemas = map[string]*ogen.Schema{
	"bool":                 ogen.Bool(),
	"string":               ogen.String(),
	"int":                  ogen.Int(),
	"int8":                 ogen.Int32(),
	"int16":                ogen.Int32(),
	"int32":                ogen.Int32(),
	"uint":                 ogen.Int32(),
	"uint8":                ogen.Int32(),
	"uint16":               ogen.Int32(),
	"uint32":               ogen.Int32(),
	"int64":                ogen.Int64(),
	"uint64":               ogen.Int64(),
	"float32":              ogen.Float(),
	"float64":              ogen.Double(),
	"time.Time":            ogen.DateTime(),
	"multipart.FileHeader": ogen.Binary(),
	"[]byte":               ogen.Bytes(),
	"uuid.NullUUID":        ogen.UUID(),
	"uuid.UUID":            ogen.UUID(),
	"render.Render":        ogen.Binary(),
	"json.RawMessage":      ogen.NewSchema().SetType(string(jsonschema.Object)),
}

func toSchemaSchema(category string, prop *ProfileProperty) *ogen.Schema {
	if prop == nil {
		return nil
	}

	switch prop.Type {
	default:
		return ogen.NewSchema().
			SetType(string(jsonschema.Object))
	case ProfileTypeBasic:
		return basicSchemas[prop.TypeDescriptor]
	case ProfileTypeArray:
		schema, exist := basicSchemas[prop.TypeDescriptor]
		if exist {
			schema = schema.AsArray()

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				schema.SetMinItems(&items)
				schema.SetMaxItems(&items)
			}

			return schema
		}
	case ProfileTypeObject:
	}

	switch prop.TypeDescriptor {
	case "object":
		schema := ogen.NewSchema().
			SetType(string(jsonschema.Object))
		if prop.Type == ProfileTypeArray {
			schema = schema.AsArray()

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				schema.SetMinItems(&items)
				schema.SetMaxItems(&items)
			}
		}

		return schema
	case "array":
		schema := ogen.NewSchema().
			SetType(string(jsonschema.Object))
		if len(prop.Properties) == 1 {
			schema = toSchemaSchema(category, &prop.Properties[0])
		}

		if prop.Type == ProfileTypeArray {
			schema = schema.AsArray()

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				schema.SetMinItems(&items)
				schema.SetMaxItems(&items)
			}

			return schema
		}
	}

	schemaID := prop.TypeDescriptor
	if category != "" {
		schemaID += "." + category
	}

	if prop.TypeRefer {
		schema := ogen.NewSchema().
			SetRef("#/components/schemas/" + schemaID)
		if prop.Type == ProfileTypeArray {
			schema = schema.AsArray()

			if prop.TypeArrayLength != 0 {
				items := uint64(prop.TypeArrayLength)
				schema.SetMinItems(&items)
				schema.SetMaxItems(&items)
			}
		}

		return schema
	}

	namedSchema := ogen.NewNamedSchema(schemaID,
		ogen.NewSchema().SetType(string(jsonschema.Object)))

	for i := 0; i < len(prop.Properties); i++ {
		var add func(...*ogen.Property) *ogen.Schema
		if prop.Properties[i].Required {
			add = namedSchema.Schema.AddRequiredProperties
		} else {
			add = namedSchema.Schema.AddOptionalProperties
		}

		add(toSchemaProperty(category, &prop.Properties[i]))
	}
	spec.AddNamedSchemas(namedSchema)

	schema := namedSchema.AsLocalRef()
	if prop.Type == ProfileTypeArray {
		schema = schema.AsArray()

		if prop.TypeArrayLength != 0 {
			items := uint64(prop.TypeArrayLength)
			schema.SetMinItems(&items)
			schema.SetMaxItems(&items)
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
