//go:build ginx

package runtime

import (
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/seal-io/walrus/pkg/apis/runtime/openapi"
	"github.com/seal-io/walrus/utils/json"
)

type (
	//nolint:staticcheck
	X struct {
		BytesValue  []byte           `json:"bytesValue,cli-ignore"`
		Int64Value  int64            `json:"int64Value"`
		JSONValue   json.RawMessage  `json:"jsonValue,omitempty"`
		UUIDValue   uuid.UUID        `json:"uuidValue,omitempty"`
		BoolValue   bool             `json:"boolValue,default=true"`
		StringValue string           `json:"stringValue,default="`
		IntValue    int              `json:"intValue,default=100"`
		AnyValue    any              `json:"anyValue,omitempty"`
		ArrayValue  [4]float64       `json:"arrayValue,omitempty"`
		MapValue    map[string]int32 `json:"mapValue,omitempty"`
		SetValue    map[int]struct{} `json:"setValue,omitempty"`

		StringPath    string   `path:"stringPath"`
		StringsPath   []string `path:"stringsPath"` // Not supported in path parameter.
		StringHeader  string   `header:"stringHeader"`
		StringsHeader []string `header:"stringsHeader,omitempty"`
		StringQuery   string   `query:"stringQuery"`
		StringsQuery  []string `query:"stringsQuery,omitempty"`

		StringForm string                `form:"stringForm"`
		FileForm   *multipart.FileHeader `form:"fileForm,omitempty"`
	}

	Y struct {
		Int64Value int64 `json:"int64Value,omitempty"`

		StringPath    string   `path:"stringPath"`
		StringsPath   []string `path:"stringsPath"` // Not supported in path parameter.
		StringHeader  string   `header:"stringHeader"`
		StringsHeader []string `header:"stringsHeader,omitempty"`
		StringQuery   string   `query:"stringQuery"`
		StringsQuery  []string `query:"stringsQuery,omitempty"`
	}

	Z map[string]int

	L int

	Xs []X

	A = X

	B struct {
		A `json:",inline" header:",inline" path:",inline" query:",inline" form:",inline"`
	}

	C struct {
		Xs `json:",inline" header:",inline" path:",inline" query:",inline" form:",inline"`
	}

	CWithoutFormTag struct {
		Xs `json:",inline" header:",inline" path:",inline" query:",inline"`
	}

	D struct {
		Z `json:",inline" header:",inline" path:",inline" query:",inline" form:",inline"`
	}

	E struct {
		L `json:",inline" header:",inline" path:",inline" query:",inline" form:",inline"`
	}

	F struct {
		*F `header:",inline" path:",inline" query:",inline"`

		YPointer *Y  `json:"yPointer"`
		YSlice   []Y `json:"ySlice,omitempty"`
		FRefer   *F  `json:"fRefer"`

		StringPath    string   `path:"stringPath"`
		StringsPath   []string `path:"stringsPath"` // Not supported in path parameter.
		StringHeader  string   `header:"stringHeader"`
		StringsHeader []string `header:"stringsHeader,omitempty"`
		StringQuery   string   `query:"stringQuery"`
		StringsQuery  []string `query:"stringsQuery,omitempty"`
	}
)

func (x X) SetStream(_ RequestUnidiStream) {}

func Test_getOperationSummaryAndDescription(t *testing.T) {
	testCases := []struct {
		name     string
		given    Route
		expected string
	}{
		{
			name: "standard post",
			given: Route{
				RouteProfile: RouteProfile{
					ResourceProfile: ResourceProfile{
						Kinds: []string{"Project"},
					},
					Method: http.MethodPost,
				},
			},
			expected: "Create a project.",
		},
		{
			name: "standard collection get with one prerequisite kinds",
			given: Route{
				RouteProfile: RouteProfile{
					ResourceProfile: ResourceProfile{
						Kinds: []string{"Project", "Environment"},
					},
					Method:     http.MethodGet,
					Collection: true,
				},
			},
			expected: "Get environments of a project.",
		},
		{
			name: "standard update with two prerequisite kinds",
			given: Route{
				RouteProfile: RouteProfile{
					ResourceProfile: ResourceProfile{
						Kinds: []string{"Project", "Environment", "Resource"},
					},
					Method:     http.MethodPut,
					Collection: true,
				},
			},
			expected: "Update resources of an environment that belongs to a project.",
		},
		{
			name: "standard collection delete with three prerequisite kinds",
			given: Route{
				RouteProfile: RouteProfile{
					ResourceProfile: ResourceProfile{
						Kinds: []string{"Project", "Environment", "Resource", "ResourceRevision"},
					},
					Method:     http.MethodDelete,
					Collection: true,
				},
			},
			expected: "Delete resource revisions that belongs to an environment of a project.",
		},
		{
			name: "long description",
			given: Route{
				RouteProfile: RouteProfile{
					ResourceProfile: ResourceProfile{
						Kinds: []string{
							"Galaxy",
							"Universal",
							"Word",
							"System",
							"Organization",
							"Project",
							"Environment",
							"Resource",
							"Component",
						},
					},
					Method:     http.MethodGet,
					Custom:     true,
					CustomName: "Log",
				},
			},
			expected: "Log for a component of a resource that belongs to an environment of a project under an organization of a system below a word of a universal of a galaxy.",
		},
		{
			name: "skip same prefix",
			given: Route{
				RouteProfile: RouteProfile{
					ResourceProfile: ResourceProfile{
						Kinds: []string{
							"Galaxy",
							"Universal",
							"Word",
							"System",
							"Organization",
							"Project",
							"Environment",
							"Resource",
							"ResourceComponent",
						},
					},
					Method:     http.MethodGet,
					Custom:     true,
					CustomName: "Log",
				},
			},
			expected: "Log for a resource component that belongs to an environment of a project under an organization of a system below a word of a universal of a galaxy.",
		},
		{
			name: "custom get",
			given: Route{
				RouteProfile: RouteProfile{
					ResourceProfile: ResourceProfile{
						Kinds: []string{"Dashboard"},
					},
					Method:     http.MethodGet,
					Custom:     true,
					CustomName: "GetNumber",
				},
			},
			expected: "Get number for a dashboard.",
		},
		{
			name: "custom post",
			given: Route{
				RouteProfile: RouteProfile{
					ResourceProfile: ResourceProfile{
						Kinds: []string{"Connector"},
					},
					Method:     http.MethodPost,
					Custom:     true,
					CustomName: "SyncCost",
				},
			},
			expected: "Sync cost for a connector.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, actual := getOperationSummaryAndDescription(&tc.given)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_getOperationParameters(t *testing.T) {
	testCases := []struct {
		name     string
		given    Route
		expected openapi3.Parameters
	}{
		{
			name: "general struct",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodGet,
					Path:   "/v1/tests/:stringPath",
				},
				RequestType:       reflect.TypeOf(A{}),
				RequestAttributes: RequestWithUnidiStream,
			},
			expected: func() (refs openapi3.Parameters) {
				pv := openapi3.NewQueryParameter("watch").
					WithSchema(openapi3.NewBoolSchema())
				pv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				refs = append(refs,
					&openapi3.ParameterRef{
						Value: openapi3.NewPathParameter("stringPath").
							WithSchema(openapi3.NewStringSchema()),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewHeaderParameter("stringHeader").
							WithSchema(openapi3.NewStringSchema()).
							WithRequired(true),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewHeaderParameter("stringsHeader").
							WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewQueryParameter("stringQuery").
							WithSchema(openapi3.NewStringSchema()).
							WithRequired(true),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewQueryParameter("stringsQuery").
							WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
					},
					&openapi3.ParameterRef{
						Value: pv,
					},
				)
				return
			}(),
		},
		{
			name: "inline struct",
			given: Route{
				RouteProfile: RouteProfile{
					Path: "/v1/tests/:stringPath",
				},
				RequestType: reflect.TypeOf(B{}),
			},
			expected: func() (refs openapi3.Parameters) {
				refs = append(refs,
					&openapi3.ParameterRef{
						Value: openapi3.NewPathParameter("stringPath").
							WithSchema(openapi3.NewStringSchema()),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewHeaderParameter("stringHeader").
							WithSchema(openapi3.NewStringSchema()).
							WithRequired(true),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewHeaderParameter("stringsHeader").
							WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewQueryParameter("stringQuery").
							WithSchema(openapi3.NewStringSchema()).
							WithRequired(true),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewQueryParameter("stringsQuery").
							WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
					},
				)

				return
			}(),
		},
		{
			name: "inline slice",
			given: Route{
				RouteProfile: RouteProfile{
					Path: "/v1/tests/:stringPath",
				},
				RequestType: reflect.TypeOf(C{}),
			},
			expected: func() (refs openapi3.Parameters) {
				return
			}(),
		},
		{
			name: "inline map",
			given: Route{
				RouteProfile: RouteProfile{
					Path: "/v1/tests/:stringPath",
				},
				RequestType: reflect.TypeOf(D{}),
			},
			expected: func() (refs openapi3.Parameters) {
				return
			}(),
		},
		{
			name: "inline basic value",
			given: Route{
				RouteProfile: RouteProfile{
					Path: "/v1/tests/:stringPath",
				},
				RequestType: reflect.TypeOf(E{}),
			},
			expected: func() (refs openapi3.Parameters) {
				return
			}(),
		},
		{
			name: "inline circular dependency struct",
			given: Route{
				RouteProfile: RouteProfile{
					Path: "/v1/tests/:stringPath",
				},
				RequestType: reflect.TypeOf(F{}),
			},
			expected: func() (refs openapi3.Parameters) {
				refs = append(refs,
					&openapi3.ParameterRef{
						Value: openapi3.NewPathParameter("stringPath").
							WithSchema(openapi3.NewStringSchema()),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewHeaderParameter("stringHeader").
							WithSchema(openapi3.NewStringSchema()).
							WithRequired(true),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewHeaderParameter("stringsHeader").
							WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewQueryParameter("stringQuery").
							WithSchema(openapi3.NewStringSchema()).
							WithRequired(true),
					},
					&openapi3.ParameterRef{
						Value: openapi3.NewQueryParameter("stringsQuery").
							WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
					},
				)

				return
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, _ := getOperationParameters(&tc.given)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_getOperationRequestBody(t *testing.T) {
	testCases := []struct {
		name     string
		given    Route
		expected *openapi3.RequestBodyRef
	}{
		{
			name: "create general struct",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodPost,
					Path:   "/v1/tests",
				},
				RequestType: reflect.TypeOf(A{}),
			},
			expected: func() *openapi3.RequestBodyRef {
				formSchema := openapi3.NewObjectSchema().
					WithProperty("stringForm", openapi3.NewStringSchema()).
					WithProperty("fileForm", openapi3.NewStringSchema().
						WithFormat("binary"))

				formSchema.Required = []string{"stringForm"}

				jsonSchema := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().WithFormat("int32"))
				mv.Extensions = map[string]any{openapi.ExtCliSchemaTypeName: "map[string]int32"}

				jsonSchema.
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				jsonSchema.Required = []string{"bytesValue", "int64Value"}

				return &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: map[string]*openapi3.MediaType{
							binding.MIMEMultipartPOSTForm: openapi3.NewMediaType().
								WithSchema(formSchema),
							binding.MIMEJSON: openapi3.NewMediaType().
								WithSchema(jsonSchema),
						},
						Required: true,
					},
				}
			}(),
		},
		{
			name: "get general struct",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodGet,
					Path:   "/v1/tests/:stringPath",
				},
				RequestType: reflect.TypeOf(A{}),
			},
			expected: nil,
		},
		{
			name: "update inline struct",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodPut,
					Path:   "/v1/tests/:stringPath",
				},
				RequestType: reflect.TypeOf(B{}),
			},
			expected: func() *openapi3.RequestBodyRef {
				jsonSchema := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().
						WithFormat("int32"))
				mv.Extensions = map[string]any{
					openapi.ExtCliSchemaTypeName: "map[string]int32",
				}

				jsonSchema.
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				jsonSchema.Required = []string{"bytesValue", "int64Value"}

				return &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: map[string]*openapi3.MediaType{
							binding.MIMEJSON: openapi3.NewMediaType().
								WithSchema(jsonSchema),
						},
						Required: true,
					},
				}
			}(),
		},
		{
			name: "delete inline struct",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodDelete,
					Path:   "/v1/tests/:stringPath",
				},
				RequestType: reflect.TypeOf(B{}),
			},
			expected: func() *openapi3.RequestBodyRef {
				jsonSchema := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().
						WithFormat("int32"))
				mv.Extensions = map[string]any{
					openapi.ExtCliSchemaTypeName: "map[string]int32",
				}

				jsonSchema.
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				jsonSchema.Required = []string{"bytesValue", "int64Value"}

				return &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: map[string]*openapi3.MediaType{
							binding.MIMEJSON: openapi3.NewMediaType().
								WithSchema(jsonSchema),
						},
						Required: true,
					},
				}
			}(),
		},
		{
			name: "create inline slice",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodPost,
					Path:   "/v1/tests",
				},
				RequestType: reflect.TypeOf(C{}),
			},
			expected: func() *openapi3.RequestBodyRef {
				jsonSchema := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().
						WithFormat("int32"))
				mv.Extensions = map[string]any{
					openapi.ExtCliSchemaTypeName: "map[string]int32",
				}

				jsonSchema.
					WithProperty("bytesValue", bv).
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				jsonSchema.Required = []string{"bytesValue", "int64Value"}

				return &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: map[string]*openapi3.MediaType{
							binding.MIMEJSON: openapi3.NewMediaType().
								WithSchema(openapi3.NewArraySchema().
									WithItems(jsonSchema)),
						},
						Required: true,
					},
				}
			}(),
		},
		{
			name: "create inline slice without form tag",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodPost,
					Path:   "/v1/tests",
				},
				RequestType: reflect.TypeOf(CWithoutFormTag{}),
			},
			expected: func() *openapi3.RequestBodyRef {
				jsonSchema := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().
						WithFormat("int32"))
				mv.Extensions = map[string]any{
					openapi.ExtCliSchemaTypeName: "map[string]int32",
				}

				jsonSchema.
					WithProperty("bytesValue", bv).
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				jsonSchema.Required = []string{"bytesValue", "int64Value"}

				return &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: map[string]*openapi3.MediaType{
							binding.MIMEJSON: openapi3.NewMediaType().
								WithSchema(openapi3.NewArraySchema().
									WithItems(jsonSchema)),
						},
						Required: true,
					},
				}
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getOperationRequestBody(&tc.given)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_getOperationHTTPResponses(t *testing.T) {
	testCases := []struct {
		name     string
		given    Route
		expected openapi3.Responses
	}{
		{
			name: "get with generate struct response",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodGet,
					Path:   "/v1/tests/:stringPath",
				},
				ResponseType: reflect.TypeOf(A{}),
			},
			expected: func() openapi3.Responses {
				c := http.StatusOK

				s := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().
						WithFormat("int32"))
				mv.Extensions = map[string]any{
					openapi.ExtCliSchemaTypeName: "map[string]int32",
				}

				s.
					WithProperty("bytesValue", bv).
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				s.Required = []string{"bytesValue", "int64Value"}

				resps := openapi3.Responses{
					strconv.Itoa(c): {
						Value: openapi3.NewResponse().
							WithDescription(http.StatusText(c)).
							WithContent(map[string]*openapi3.MediaType{
								binding.MIMEJSON: {
									Schema: s.NewRef(),
								},
							}),
					},
				}

				return referErrorResponses(resps)
			}(),
		},
		{
			name: "create with inline struct response",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodPost,
					Path:   "/v1/tests",
				},
				ResponseType: reflect.TypeOf(B{}),
			},
			expected: func() openapi3.Responses {
				c := http.StatusCreated

				s := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().
						WithFormat("int32"))
				mv.Extensions = map[string]any{
					openapi.ExtCliSchemaTypeName: "map[string]int32",
				}

				s.
					WithProperty("bytesValue", bv).
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				s.Required = []string{"bytesValue", "int64Value"}

				resps := openapi3.Responses{
					strconv.Itoa(c): {
						Value: openapi3.NewResponse().
							WithDescription(http.StatusText(c)).
							WithContent(map[string]*openapi3.MediaType{
								binding.MIMEJSON: {
									Schema: s.NewRef(),
								},
							}),
					},
				}

				return referErrorResponses(resps)
			}(),
		},
		{
			name: "get with inline slice response",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodGet,
					Path:   "/v1/tests",
				},
				ResponseType: reflect.TypeOf(C{}),
			},
			expected: func() openapi3.Responses {
				c := http.StatusOK

				s := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().
						WithFormat("int32"))
				mv.Extensions = map[string]any{
					openapi.ExtCliSchemaTypeName: "map[string]int32",
				}

				s.
					WithProperty("bytesValue", bv).
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				s.Required = []string{"bytesValue", "int64Value"}

				resps := openapi3.Responses{
					strconv.Itoa(c): {
						Value: openapi3.NewResponse().
							WithDescription(http.StatusText(c)).
							WithContent(map[string]*openapi3.MediaType{
								binding.MIMEJSON: {
									Schema: openapi3.NewArraySchema().
										WithItems(s).
										NewRef(),
								},
							}),
					},
				}

				return referErrorResponses(resps)
			}(),
		},
		{
			name: "get with inline map response",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodGet,
					Path:   "/v1/tests",
				},
				ResponseType: reflect.TypeOf(D{}),
			},
			expected: func() openapi3.Responses {
				c := http.StatusOK

				s := openapi3.NewObjectSchema()

				s.WithAdditionalProperties(openapi3.NewIntegerSchema())

				resps := openapi3.Responses{
					strconv.Itoa(c): {
						Value: openapi3.NewResponse().
							WithDescription(http.StatusText(c)).
							WithContent(map[string]*openapi3.MediaType{
								binding.MIMEJSON: {
									Schema: s.NewRef(),
								},
							}),
					},
				}

				return referErrorResponses(resps)
			}(),
		},
		{
			name: "get with inline slice paging response",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodGet,
					Path:   "/v1/tests",
				},
				ResponseType:       reflect.TypeOf(C{}),
				ResponseAttributes: ResponseWithPage,
			},
			expected: func() openapi3.Responses {
				c := http.StatusOK

				s := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().
						WithFormat("int32"))
				mv.Extensions = map[string]any{
					openapi.ExtCliSchemaTypeName: "map[string]int32",
				}

				s.
					WithProperty("bytesValue", bv).
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				s.Required = []string{"bytesValue", "int64Value"}

				resps := openapi3.Responses{
					strconv.Itoa(c): {
						Value: openapi3.NewResponse().
							WithDescription(http.StatusText(c)).
							WithContent(map[string]*openapi3.MediaType{
								binding.MIMEJSON: {
									Schema: openapi3.NewObjectSchema().
										WithProperty("type", openapi3.NewStringSchema()).
										WithProperty("items", openapi3.NewArraySchema().
											WithItems(s)).
										WithProperty("pagination",
											openapi3.NewObjectSchema().
												WithProperty("page", openapi3.NewIntegerSchema()).
												WithProperty("perPage", openapi3.NewIntegerSchema()).
												WithProperty("total", openapi3.NewIntegerSchema()).
												WithProperty("totalPage", openapi3.NewIntegerSchema()).
												WithProperty("partial", openapi3.NewIntegerSchema()).
												WithProperty("nextPage", openapi3.NewIntegerSchema())).
										NewRef(),
								},
							}),
					},
				}

				return referErrorResponses(resps)
			}(),
		},
		{
			name: "get with binary",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodGet,
					Path:   "/v1/tests",
				},
				ResponseType: reflect.TypeOf(&render.JSON{}),
			},
			expected: func() openapi3.Responses {
				c := http.StatusOK

				s := openapi3.NewStringSchema().WithFormat("binary")

				resps := openapi3.Responses{
					strconv.Itoa(c): {
						Value: openapi3.NewResponse().
							WithDescription(http.StatusText(c)).
							WithContent(map[string]*openapi3.MediaType{
								"application/octet-stream": {
									Schema: s.NewRef(),
								},
							}),
					},
				}

				return referErrorResponses(resps)
			}(),
		},
		{
			name: "get with bytes",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodGet,
					Path:   "/v1/tests",
				},
				ResponseType: reflect.TypeOf([]byte{}),
			},
			expected: func() openapi3.Responses {
				c := http.StatusOK

				s := openapi3.NewBytesSchema()

				resps := openapi3.Responses{
					strconv.Itoa(c): {
						Value: openapi3.NewResponse().
							WithDescription(http.StatusText(c)).
							WithContent(map[string]*openapi3.MediaType{
								"application/octet-stream": {
									Schema: s.NewRef(),
								},
							}),
					},
				}

				return referErrorResponses(resps)
			}(),
		},
		{
			name: "create without response",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodPost,
					Path:   "/v1/tests",
				},
				ResponseType: nil,
			},
			expected: func() openapi3.Responses {
				c := http.StatusAccepted

				resps := openapi3.Responses{
					strconv.Itoa(c): {
						Value: openapi3.NewResponse().
							WithDescription(http.StatusText(c)).
							WithContent(map[string]*openapi3.MediaType{
								binding.MIMEJSON: {
									Schema: openapi3.NewObjectSchema().
										WithProperty("status", openapi3.NewIntegerSchema().
											WithDefault(c)).
										WithProperty("statusText", openapi3.NewStringSchema().
											WithDefault(http.StatusText(c))).
										WithProperty("message", openapi3.NewStringSchema()).
										NewRef(),
								},
							}),
					},
				}

				return referErrorResponses(resps)
			}(),
		},
		{
			name: "update without response",
			given: Route{
				RouteProfile: RouteProfile{
					Method: http.MethodPut,
					Path:   "/v1/tests",
				},
				ResponseType: nil,
			},
			expected: func() openapi3.Responses {
				c := http.StatusAccepted

				resps := openapi3.Responses{
					strconv.Itoa(c): {
						Value: openapi3.NewResponse().
							WithDescription(http.StatusText(c)).
							WithContent(map[string]*openapi3.MediaType{
								binding.MIMEJSON: {
									Schema: openapi3.NewObjectSchema().
										WithProperty("status", openapi3.NewIntegerSchema().
											WithDefault(c)).
										WithProperty("statusText", openapi3.NewStringSchema().
											WithDefault(http.StatusText(c))).
										WithProperty("message", openapi3.NewStringSchema()).
										NewRef(),
								},
							}),
					},
				}

				return referErrorResponses(resps)
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getOperationHTTPResponses(&tc.given)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_getSchemaOfGoType(t *testing.T) {
	type input struct {
		typ      reflect.Type
		category string
	}

	testCases := []struct {
		name     string
		given    input
		expected *openapi3.SchemaRef
	}{
		{
			name: "general struct",
			given: input{
				typ:      reflect.TypeOf(A{}),
				category: "json",
			},
			expected: func() *openapi3.SchemaRef {
				s := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().
						WithFormat("int32"))
				mv.Extensions = map[string]any{
					openapi.ExtCliSchemaTypeName: "map[string]int32",
				}

				s.
					WithProperty("bytesValue", bv).
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				s.Required = []string{"bytesValue", "int64Value"}

				return s.NewRef()
			}(),
		},
		{
			name: "inline struct",
			given: input{
				typ:      reflect.TypeOf(B{}),
				category: "json",
			},
			expected: func() *openapi3.SchemaRef {
				s := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().
						WithFormat("int32"))
				mv.Extensions = map[string]any{
					openapi.ExtCliSchemaTypeName: "map[string]int32",
				}

				s.
					WithProperty("bytesValue", bv).
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				s.Required = []string{"bytesValue", "int64Value"}

				return s.NewRef()
			}(),
		},
		{
			name: "inline slice",
			given: input{
				typ:      reflect.TypeOf(C{}),
				category: "json",
			},
			expected: func() *openapi3.SchemaRef {
				s := openapi3.NewObjectSchema()

				bv := openapi3.NewBytesSchema()
				bv.Extensions = map[string]any{openapi.ExtCliIgnore: true}

				mv := openapi3.NewObjectSchema().
					WithAdditionalProperties(openapi3.NewIntegerSchema().
						WithFormat("int32"))
				mv.Extensions = map[string]any{
					openapi.ExtCliSchemaTypeName: "map[string]int32",
				}

				s.
					WithProperty("bytesValue", bv).
					WithProperty("bytesValue", bv).
					WithProperty("int64Value", openapi3.NewIntegerSchema().
						WithFormat("int64")).
					WithProperty("jsonValue", openapi3.NewBytesSchema()).
					WithProperty("uuidValue", openapi3.NewUUIDSchema()).
					WithProperty("boolValue", openapi3.NewBoolSchema().
						WithDefault(true)).
					WithProperty("stringValue", openapi3.NewStringSchema().
						WithDefault("")).
					WithProperty("intValue", openapi3.NewIntegerSchema().
						WithDefault(int64(100))).
					WithProperty("anyValue", openapi3.NewObjectSchema()).
					WithProperty("arrayValue", openapi3.NewArraySchema().
						WithItems(openapi3.NewFloat64Schema().WithFormat("double")).
						WithMinItems(4).
						WithMaxItems(4)).
					WithProperty("mapValue", mv).
					WithProperty("setValue", openapi3.NewObjectSchema())

				s.Required = []string{"bytesValue", "int64Value"}

				return openapi3.NewArraySchema().WithItems(s).NewRef()
			}(),
		},
		{
			name: "inline map",
			given: input{
				typ:      reflect.TypeOf(D{}),
				category: "json",
			},
			expected: func() *openapi3.SchemaRef {
				s := openapi3.NewObjectSchema()

				s.WithAdditionalProperties(openapi3.NewIntegerSchema())

				return s.NewRef()
			}(),
		},
		{
			name: "inline basic value",
			given: input{
				typ:      reflect.TypeOf(E{}),
				category: "json",
			},
			expected: nil,
		},
		{
			name: "circular dependency struct",
			given: input{
				typ:      reflect.TypeOf(F{}),
				category: "json",
			},
			expected: func() *openapi3.SchemaRef {
				s := openapi3.NewObjectSchema()

				s.
					WithProperty("yPointer",
						openapi3.NewObjectSchema().
							WithProperty("int64Value",
								openapi3.NewIntegerSchema().WithFormat("int64")))

				yslice := openapi3.NewArraySchema()
				yslice.Items = openapi3.NewSchemaRef("#/components/schemas/Ginx.Y.json", nil)
				s.WithPropertyRef("ySlice", yslice.NewRef())

				s.WithPropertyRef("fRefer", openapi3.NewSchemaRef("#/components/schemas/Ginx.F.json", nil))

				s.Required = []string{"yPointer", "fRefer"}

				return s.NewRef()
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getSchemaOfGoType([]string{"Ginx"}, tc.given.typ, tc.given.category, nil)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
