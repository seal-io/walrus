package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"golang.org/x/exp/slices"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/apis/runtime/openapi"
	"github.com/seal-io/walrus/utils/hash"
	"github.com/seal-io/walrus/utils/strs"
	"github.com/seal-io/walrus/utils/version"
)

var openAPISchemas = &openapi3.T{
	OpenAPI: "3.0.3",
	Info: &openapi3.Info{
		Title:       "Restful APIs",
		Description: "Restful APIs to access Seal.",
		Version:     version.Version,
	},
	Security: getSecurityRequirements(),
	Components: &openapi3.Components{
		SecuritySchemes: getSecuritySchemes(),
		Responses:       getErrorResponses(),
		Schemas:         openapi3.Schemas{},
	},
}

const (
	parameterInPath   = openapi3.ParameterInPath
	parameterInHeader = openapi3.ParameterInHeader
	parameterInQuery  = openapi3.ParameterInQuery

	contentInForm = "form"
	contentInJSON = "json"
)

func schemeRoute(bp string, r *Route) error {
	// Get operation.
	op, err := getOperationSchema(r)
	if err != nil {
		return err
	}

	// Extend operation.
	extendOperationSchema(r, op)

	// Normalize path.
	ss := strings.Split(path.Join(bp, r.Path), "/")
	for i := range ss {
		if ss[i] == "" || ss[i][0] != ':' {
			continue
		}
		ss[i] = "{" + ss[i][1:] + "}"
	}

	// Bind operation.
	openAPISchemas.AddOperation(strs.Join("/", ss...), r.Method, op)

	return nil
}

func getOperationSchema(r *Route) (*openapi3.Operation, error) {
	op := &openapi3.Operation{
		OperationID: hash.SumStrings(r.Method, r.Path),
	}

	// Tags.
	if len(r.Kinds) != 0 {
		op.Tags = []string{strs.Pluralize(r.Kinds[len(r.Kinds)-1])}
	} else {
		op.Tags = []string{"None Resources"}
	}

	// Summary and description.
	op.Summary, op.Description = getOperationSummaryAndDescription(r)

	// Parameters.
	parameters, err := getOperationParameters(r)
	if err != nil {
		return nil, fmt.Errorf("failed to scheme parameters: %w", err)
	}

	if r.RequestAttributes.HasAll(RequestWithBidiStream) {
		// Append Websocket header parameters.
		parameters = append(parameters, getOperationWebsocketAdditionalParameters()...)
	}

	op.Parameters = parameters

	// Responses.
	if r.RequestAttributes.HasAll(RequestWithBidiStream) {
		op.Responses = getOperationWebsocketResponses()
		op.Description = "[!!! Websocket Connection !!!] " + op.Description
	} else {
		op.RequestBody = getOperationRequestBody(r)
		op.Responses = getOperationHTTPResponses(r)
	}

	// Securities.
	if openAPISchemas.Security != nil {
		op.Security = &openAPISchemas.Security
	}

	return op, nil
}

func getOperationSummaryAndDescription(r *Route) (summary, description string) {
	var sb strings.Builder

	if r.Custom {
		sb.WriteString(strs.Decamelize(r.CustomName, true))
	} else {
		switch r.Method {
		case http.MethodPost:
			sb.WriteString("create ")
		case http.MethodPut:
			sb.WriteString("update ")
		case http.MethodDelete:
			sb.WriteString("delete ")
		default:
			sb.WriteString("get ")
		}
	}

	if len(r.Kinds) == 0 {
		s := sb.String() + "."

		bs := strs.ToBytes(&s)
		bs[0] = bs[0] + byte('A') - byte('a')
		s = strs.FromBytes(&bs)

		return s, s
	}

	if r.Custom {
		sb.WriteString(" for ")
	}

	subject := strs.Decamelize(r.Kinds[len(r.Kinds)-1], true)
	if r.Collection {
		subject = strs.Pluralize(subject)
	}

	sb.WriteString(strs.Article(subject))

	// NB(thxCode): mark the summary of the route.
	r.Summary = sb.String()

	summary = sb.String() + "."
	summaryBs := strs.ToBytes(&summary)
	summaryBs[0] = summaryBs[0] + byte('A') - byte('a')
	summary = strs.FromBytes(&summaryBs)

	// For example,
	// if the method is POST, and the kinds are ["H", "I", "X", "Y", "Z", "L", "M", "N"]
	// then description will be something like below,
	// "Create N of M that belongs to L under Z of Y below X of I of H.".
	preps := []string{" that belongs to ", " under ", " below "}
	for l, i := len(r.Kinds)-2, len(r.Kinds)-2; i >= 0; i-- {
		if (l-i)%2 == 1 && (l-i)/2 < len(preps) {
			sb.WriteString(preps[(l-i)/2])
		} else {
			if strings.HasPrefix(r.Kinds[i+1], r.Kinds[i]) {
				continue
			}

			sb.WriteString(" of ")
		}

		sb.WriteString(strs.Article(strs.Decamelize(r.Kinds[i], true)))
	}

	// NB(thxCode): mark the description of the route.
	r.Description = sb.String()

	description = sb.String() + "."
	descriptionBs := strs.ToBytes(&description)
	descriptionBs[0] = descriptionBs[0] + byte('A') - byte('a')
	description = strs.FromBytes(&descriptionBs)

	return summary, description
}

const (
	stringTypeFormatBinary = "binary"
	stringTypeFormatByte   = "byte"
)

func getOperationParameters(r *Route) (openapi3.Parameters, error) {
	// Get knowledge of path parameters.
	pathValidParamIndex := make(map[string]int)

	for i, sg := range strings.Split(r.Path, "/") {
		if sg == "" || sg[0] != ':' {
			continue
		}

		if _, ok := pathValidParamIndex[sg[1:]]; !ok {
			pathValidParamIndex[sg[1:]] = i
		}
	}

	var refs openapi3.Parameters

	// Parse.
	for _, category := range []string{
		parameterInPath,
		parameterInHeader,
		parameterInQuery,
	} {
		schemaRef := getSchemaOfGoType(r.Kinds, r.RequestType, category, nil)

		if schemaRef == nil ||
			schemaRef.Value == nil || len(schemaRef.Value.Properties) == 0 {
			continue
		}

		schemaRefRequiredSet := sets.New[string](schemaRef.Value.Required...)

		for pn, p := range flattenSchemas(schemaRef.Value.Properties) {
			if category == parameterInPath {
				_, ok := pathValidParamIndex[pn]
				if !ok || p.Value == nil || p.Value.Type == openapi3.TypeArray {
					continue
				}
			}

			if p.Value != nil && p.Value.Type == openapi3.TypeArray {
				pvi := p.Value.Items
				if pvi.Ref != "" ||
					pvi.Value.Type == openapi3.TypeObject ||
					pvi.Value.Format == stringTypeFormatBinary ||
					pvi.Value.Format == stringTypeFormatByte {
					continue
				}
			}

			refs = append(refs, &openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       category,
					Name:     pn,
					Schema:   p,
					Required: category == parameterInPath || schemaRefRequiredSet.Has(pn),
				},
			})

			// NB(thxCode): mark the binder of the route.
			switch category {
			case parameterInHeader:
				r.RequestAttributes.With(RequestWithBindingHeader)
			case parameterInQuery:
				r.RequestAttributes.With(RequestWithBindingQuery)
			case parameterInPath:
				r.RequestAttributes.With(RequestWithBindingPath)
			}
		}

		if category == parameterInPath &&
			len(refs) != len(pathValidParamIndex) {
			return nil, errors.New("invalid path parameters")
		}
	}

	// Sort.
	sort.SliceStable(refs, func(i, j int) bool {
		ri, rj := refs[i].Value, refs[j].Value

		if ri.In != rj.In {
			switch {
			case ri.In == parameterInPath:
				return true
			case ri.In == parameterInHeader && rj.In != parameterInPath:
				return true
			case rj.In == parameterInQuery:
				return true
			}

			return false
		}

		if ri.In == parameterInPath {
			return pathValidParamIndex[ri.Name] < pathValidParamIndex[rj.Name]
		}

		return ri.Required || !rj.Required
	})

	// Append.
	if r.RequestAttributes.HasAll(RequestWithUnidiStream) {
		refs = append(refs, &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				In:       parameterInQuery,
				Name:     "watch",
				Schema:   openapi3.NewBoolSchema().NewRef(),
				Required: false,
				Extensions: map[string]any{
					openapi.ExtCliIgnore: true,
				},
			},
		})
	}

	return refs, nil
}

func getOperationWebsocketAdditionalParameters() openapi3.Parameters {
	return openapi3.Parameters{
		&openapi3.ParameterRef{
			Value: openapi3.NewHeaderParameter("Connection").
				WithRequired(true),
		},
		&openapi3.ParameterRef{
			Value: openapi3.NewHeaderParameter("Upgrade").
				WithRequired(true),
		},
		&openapi3.ParameterRef{
			Value: openapi3.NewHeaderParameter("Sec-WebSocket-Key").
				WithRequired(true),
		},
		&openapi3.ParameterRef{
			Value: openapi3.NewHeaderParameter("Sec-WebSocket-Version").
				WithRequired(true),
		},
	}
}

func getOperationRequestBody(r *Route) *openapi3.RequestBodyRef {
	if r.Method == http.MethodGet {
		return nil
	}

	categoryContentTypes := map[string]string{
		contentInForm: binding.MIMEMultipartPOSTForm,
		contentInJSON: binding.MIMEJSON,
	}
	if r.Method != http.MethodPost {
		delete(categoryContentTypes, contentInForm)
	}

	requestBody := &openapi3.RequestBody{
		Content:  make(map[string]*openapi3.MediaType, len(categoryContentTypes)),
		Required: true,
	}

	for category, contentType := range categoryContentTypes {
		schemaRef := getSchemaOfGoType(r.Kinds, r.RequestType, category, nil)

		if schemaRef == nil {
			continue
		}

		if category == contentInForm &&
			(schemaRef.Value != nil && schemaRef.Value.Type == openapi3.TypeArray) {
			continue
		}

		requestBody.Content[contentType] = openapi3.NewMediaType().
			WithSchemaRef(schemaRef)

		// NB(thxCode): mark the binder of the route.
		switch contentType {
		case binding.MIMEMultipartPOSTForm:
			r.RequestAttributes.With(RequestWithBindingForm)
		case binding.MIMEJSON:
			r.RequestAttributes.With(RequestWithBindingJSON)
		}
	}

	if len(requestBody.Content) == 0 {
		return nil
	}

	return &openapi3.RequestBodyRef{
		Value: requestBody,
	}
}

func getOperationHTTPResponses(r *Route) openapi3.Responses {
	schemaRef := getSchemaOfGoType(r.Kinds, r.ResponseType, contentInJSON, nil)

	contentType := binding.MIMEJSON

	c := http.StatusOK
	if !r.Custom && r.Method == http.MethodPost {
		c = http.StatusCreated
	}

	switch {
	case schemaRef == nil:
		if !r.Custom && r.Method != http.MethodGet {
			c = http.StatusAccepted
		}

		// No response.
		schemaRef = openapi3.NewObjectSchema().
			WithProperty("status",
				openapi3.NewIntegerSchema().WithDefault(c)).
			WithProperty("statusText",
				openapi3.NewStringSchema().WithDefault(http.StatusText(c))).
			WithProperty("message",
				openapi3.NewStringSchema()).
			NewRef()

	case schemaRef.Value != nil:
		switch {
		case !r.Custom && r.Method == http.MethodPost && r.Collection:
			// Response of collection creation.
			schemaRef = openapi3.NewObjectSchema().
				WithProperty("items", schemaRef.Value).
				NewRef()
		case schemaRef.Value.Type == openapi3.TypeArray && r.ResponseAttributes.HasAll(ResponseWithPage):
			// Response in pagination or stream.
			schemaRef = openapi3.NewObjectSchema().
				WithProperty("type",
					openapi3.NewStringSchema()).
				WithProperty("items",
					schemaRef.Value).
				WithProperty("pagination",
					openapi3.NewObjectSchema().
						WithProperty("page",
							openapi3.NewIntegerSchema()).
						WithProperty("perPage",
							openapi3.NewIntegerSchema()).
						WithProperty("total",
							openapi3.NewIntegerSchema()).
						WithProperty("totalPage",
							openapi3.NewIntegerSchema()).
						WithProperty("partial",
							openapi3.NewIntegerSchema()).
						WithProperty("nextPage",
							openapi3.NewIntegerSchema())).
				NewRef()
		}
	}

	if schemaRef.Value != nil && schemaRef.Value.Type != openapi3.TypeArray &&
		schemaRef.Value.Type != openapi3.TypeObject {
		// Response in bytes.
		contentType = "application/octet-stream"
	}

	resps := openapi3.Responses{
		strconv.Itoa(c): {
			Value: openapi3.NewResponse().
				WithDescription(http.StatusText(c)).
				WithContent(map[string]*openapi3.MediaType{
					contentType: {
						Schema: schemaRef,
					},
				}),
		},
	}

	return referErrorResponses(resps)
}

func getOperationWebsocketResponses() openapi3.Responses {
	_101 := openapi3.NewResponse().
		WithDescription("Switching Protocols")
	_101.Headers = openapi3.Headers{
		"Connection": {
			Value: &openapi3.Header{Parameter: openapi3.Parameter{Required: true}},
		},
		"Upgrade": {
			Value: &openapi3.Header{Parameter: openapi3.Parameter{Required: true}},
		},
		"Sec-WebSocket-Accept": {
			Value: &openapi3.Header{Parameter: openapi3.Parameter{Required: true}},
		},
	}

	resps := openapi3.Responses{
		"101": {Value: _101},
	}

	return referErrorResponses(resps)
}

func getSecurityRequirements() openapi3.SecurityRequirements {
	requires := openapi3.SecurityRequirements{}

	for provider := range getSecuritySchemes() {
		requires.With(openapi3.NewSecurityRequirement().
			Authenticate(provider))
	}

	return requires
}

func getSecuritySchemes() openapi3.SecuritySchemes {
	schemes := openapi3.SecuritySchemes{}

	schemes["BasicAuth"] = &openapi3.SecuritySchemeRef{
		Value: openapi3.NewSecurityScheme().
			WithType("http").
			WithIn("header").
			WithScheme("basic").
			WithDescription("Basic Authentication, in form of base64(<username>:<password>), " +
				"the password must be a valid Seal API token."),
	}

	schemes["BearerAuth"] = &openapi3.SecuritySchemeRef{
		Value: openapi3.NewSecurityScheme().
			WithType("http").
			WithIn("header").
			WithScheme("bearer").
			WithDescription("Bearer Authentication, the token must be a valid Seal API token."),
	}

	return schemes
}

func getErrorResponses() openapi3.Responses {
	httpc := getErrorResponseStatus()
	resps := openapi3.Responses{}

	for _, c := range httpc {
		resps[strconv.Itoa(c)] = &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription(http.StatusText(c)).
				WithContent(openapi3.NewContentWithJSONSchema(
					openapi3.NewObjectSchema().
						WithProperty("status", openapi3.NewIntegerSchema().
							WithDefault(c)).
						WithProperty("statusText", openapi3.NewStringSchema().
							WithDefault(http.StatusText(c))).
						WithProperty("message", openapi3.NewStringSchema())),
				),
		}
	}

	return resps
}

func referErrorResponses(resps openapi3.Responses) openapi3.Responses {
	for _, s := range getErrorResponseStatus() {
		k := strconv.Itoa(s)
		resps[k] = &openapi3.ResponseRef{
			Ref: "#/components/responses/" + k,
		}
	}

	return resps
}

func getErrorResponseStatus() []int {
	return []int{
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusUnsupportedMediaType,
		http.StatusUnprocessableEntity,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
	}
}

func getSchemaOfGoType(
	kinds []string,
	typ reflect.Type,
	category string,
	visited sets.Set[string],
) *openapi3.SchemaRef {
	if typ == nil {
		return nil
	}

	switch typ.Kind() {
	case reflect.Bool:
		return openapi3.NewBoolSchema().NewRef()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := openapi3.NewIntegerSchema()

		switch typ.Kind() {
		case reflect.Uint16, reflect.Int32:
			s.WithFormat("int32")
		case reflect.Uint32, reflect.Int64:
			s.WithFormat("int64")
		}

		return s.NewRef()
	case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		s := openapi3.NewFloat64Schema()

		switch typ.Kind() {
		case reflect.Float32:
			s.WithFormat("float")
		case reflect.Float64:
			s.WithFormat("double")
		}

		return s.NewRef()
	case reflect.String:
		return openapi3.NewStringSchema().NewRef()
	case reflect.Pointer:
		return getSchemaOfGoType(kinds, typ.Elem(), category, visited)
	case reflect.Interface:
		ts := typ.String()

		switch ts {
		case "render.Render":
			return openapi3.NewStringSchema().WithFormat(stringTypeFormatBinary).NewRef()
		case "json.RawMessage":
			return openapi3.NewBytesSchema().NewRef()
		}

		return openapi3.NewObjectSchema().NewRef()
	case reflect.Array, reflect.Slice:
		ts := typ.String()

		switch ts {
		case "[]uint8", "json.RawMessage":
			return openapi3.NewBytesSchema().NewRef()
		case "uuid.UUID":
			return openapi3.NewUUIDSchema().NewRef()
		}

		var s *openapi3.Schema

		if tes := typ.Elem().String(); tes == "byte" || tes == "uint8" {
			s = openapi3.NewBytesSchema()
		} else {
			s = openapi3.NewArraySchema()

			isr := getSchemaOfGoType(kinds, typ.Elem(), category, visited)
			if isr != nil {
				if isr.Ref != "" {
					s.Items = isr
				} else {
					s.WithItems(isr.Value)
				}
			} else {
				s.WithItems(openapi3.NewObjectSchema())
			}
		}

		if typ.Kind() == reflect.Array {
			s.WithMaxItems(int64(typ.Len())).
				WithMinItems(int64(typ.Len()))
		}

		return s.NewRef()
	case reflect.Map:
		s := openapi3.NewObjectSchema()

		if typ.Key().Kind() == reflect.String {
			apsr := getSchemaOfGoType(kinds, typ.Elem(), category, visited)
			if apsr != nil {
				if apsr.Ref != "" {
					s.AdditionalProperties = openapi3.AdditionalProperties{Schema: apsr}
				} else {
					s.WithAdditionalProperties(apsr.Value)
				}
			}

			s.Extensions = map[string]any{
				openapi.ExtCliSchemaTypeName: "map[string]" + typ.Elem().String(),
			}
		}

		return s.NewRef()
	case reflect.Struct:
		ts := typ.String()

		switch ts {
		case "time.Time":
			return openapi3.NewDateTimeSchema().NewRef()
		case "multipart.FileHeader":
			return openapi3.NewStringSchema().WithFormat(stringTypeFormatBinary).NewRef()
		case "uuid.NullUUID":
			return openapi3.NewUUIDSchema().NewRef()
		}

		if typ.ConvertibleTo(reflect.TypeOf((*render.Render)(nil)).Elem()) {
			return openapi3.NewStringSchema().WithFormat(stringTypeFormatBinary).NewRef()
		}

		if visited == nil {
			visited = sets.Set[string]{}
		}

		id := strs.Join(".", typ.Name(), category)
		if len(kinds) != 0 {
			id = strs.Join(".", append(kinds[:len(kinds):len(kinds)], id)...)
		}

		if visited.Has(id) {
			return openapi3.NewSchemaRef("#/components/schemas/"+id, nil)
		}

		visited.Insert(id)

		s := openapi3.NewObjectSchema()

		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if f.PkgPath != "" && !f.Anonymous {
				continue
			}

			fa := getCategorizedAttrsOfGoField(f, category)
			if fa == nil {
				continue
			}

			if fa.Name() == "" && !fa.IsInline() {
				continue
			}

			fsr := getSchemaOfGoType(kinds, f.Type, category, visited)
			if fsr == nil {
				continue
			}

			if fsr.Ref != "" {
				s.WithPropertyRef(fa.Name(), fsr)
			} else {
				fs := fsr.Value

				if fa.IsInline() {
					// Exclusive definition if embedded an array.
					if fs.Type == openapi3.TypeArray {
						s = fs
						break
					}

					// Exclusive definition if embedded a byte array.
					if fs.Type == openapi3.TypeString &&
						slices.Contains([]string{stringTypeFormatBinary, stringTypeFormatByte}, fs.Format) {
						s = fs
						break
					}

					// Otherwise, merge everything.
					for pn, p := range fs.Properties {
						if pn == "" {
							continue
						}

						s.WithProperty(pn, p.Value)
					}

					if sv := fs.AdditionalProperties.Schema; sv != nil && sv.Value != nil {
						s.WithAdditionalProperties(sv.Value)
					}

					s.Required = append(s.Required, fs.Required...)

					continue
				}

				for extK, extV := range fa.Extensions() {
					if fs.Extensions == nil {
						fs.Extensions = make(map[string]any)
					}

					fs.Extensions[extK] = extV
				}
				fs.Default = fa.Default()
				s.WithProperty(fa.Name(), fs)
			}

			if fa.IsRequired() {
				s.Required = append(s.Required, fa.Name())
			}
		}

		if len(s.Properties) == 0 &&
			s.AdditionalProperties.Schema == nil &&
			s.Items == nil &&
			!(s.Type == openapi3.TypeString &&
				slices.Contains([]string{stringTypeFormatBinary, stringTypeFormatByte}, s.Format)) {
			return nil
		}

		sr := s.NewRef()

		// Scheme.
		switch category {
		case contentInForm, contentInJSON:
			openAPISchemas.Components.Schemas[id] = sr
		}

		return sr
	}

	return nil
}

func getCategorizedAttrsOfGoField(f reflect.StructField, category string) attributes {
	tag := f.Tag.Get(category)
	if tag == "" || tag == "-" {
		return nil
	}

	attr := attributes{}

	tgs := strings.SplitN(tag, ",", 2)

	attr["$type"] = f.Type.String()
	attr["$name"] = strings.TrimSpace(tgs[0])

	if len(tgs) == 1 {
		return attr
	}

	tgs[1] = strings.TrimSpace(tgs[1])
	if tgs[1] == "" {
		return attr
	}

	tgvs := strings.Split(tgs[1], ",")
	for i := range tgvs {
		tgvss := strings.SplitN(tgvs[i], "=", 2)

		if len(tgvss) == 1 {
			attr[strings.TrimSpace(tgvss[0])] = ""
			continue
		}

		attr[strings.TrimSpace(tgvss[0])] = strings.TrimSpace(tgvss[1])
	}

	return attr
}

type attributes map[string]string

func (attrs attributes) Name() string {
	return attrs["$name"]
}

func (attrs attributes) IsInline() bool {
	_, exist := attrs["inline"]
	return exist
}

func (attrs attributes) IsRequired() bool {
	if _, exist := attrs["omitempty"]; exist {
		return false
	}

	if _, exist := attrs["default"]; exist {
		return false
	}

	return true
}

func (attrs attributes) Default() any {
	v, exist := attrs["default"]
	if !exist {
		return nil
	}

	switch attrs["$type"] {
	case "bool":
		r, err := strconv.ParseBool(v)
		if err != nil {
			return nil
		}

		return r
	case "int", "int8", "int16", "int32", "int64":
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil
		}

		return r
	case "uint", "uint8", "uint16", "uint32", "uint64":
		r, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return r
		}
	case "float32", "float64":
		r, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return r
		}
	case "string":
		return v
	}

	return nil
}

func (attrs attributes) Extensions() (exts map[string]any) {
	for n := range attrs {
		if !strings.HasPrefix(n, "cli") {
			continue
		}

		if exts == nil {
			exts = map[string]any{}
		}
		exts[fmt.Sprintf("x-%s", n)] = true
	}

	return exts
}

func flattenSchemas(i openapi3.Schemas) (o openapi3.Schemas) {
	o = openapi3.Schemas{}

	for pn, p := range i {
		if pn == "" || p.Value == nil {
			continue
		}

		if p.Value.Type == openapi3.TypeObject {
			for pn, p := range flattenSchemas(p.Value.Properties) {
				o[pn] = p
			}

			continue
		}

		o[pn] = p
	}

	return
}
