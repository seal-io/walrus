package generators

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/seal-io/utils/stringx"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/code-generator/third_party/forked/golang/reflect"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	"github.com/seal-io/code-generator/utils"
)

var (
	knownScopes      = sets.New(apiext.NamespaceScoped, apiext.ClusterScoped)
	knownDataTypes   = sets.New("integer", "number", "string", "boolean")
	knownDataFormats = sets.New("int32", "int64", "float", "double", "byte", "binary", "date", "date-time", "password")
)

// reflectType reflects the given package and type into a CRDTypeDefinition,
// according to the given package.
func reflectType(p *types.Package, t *types.Type) *CRDTypeDefinition {
	if p == nil || t == nil {
		return nil
	}

	logger := klog.Background().
		WithName("$").
		WithValues("gen", "crd-gen", "type", t.String())

	// Collect package markers.
	//
	// +groupName=
	// +versionName=
	pm := map[string][]string{}
	collectMarkers(p.Comments, pm)

	if len(pm) == 0 || len(pm["group"]) == 0 || len(pm["version"]) == 0 {
		return nil
	}

	var (
		group    = pm["group"][len(pm["group"])-1]
		version  = pm["version"][len(pm["version"])-1]
		kind     = t.Name.Name
		singular = strings.ToLower(kind)
		plural   = strings.ToLower(stringx.Pluralize(kind))
	)

	crd := apiext.CustomResourceDefinition{
		TypeMeta: meta.TypeMeta{
			APIVersion: apiext.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: meta.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", plural, group),
		},
		Spec: apiext.CustomResourceDefinitionSpec{
			Group: group,
			Names: apiext.CustomResourceDefinitionNames{
				Kind:     kind,
				ListKind: fmt.Sprintf("%sList", kind),
				Plural:   plural,
				Singular: singular,
			},
			Scope: apiext.NamespaceScoped,
			Versions: []apiext.CustomResourceDefinitionVersion{
				{
					Name:    version,
					Served:  true,
					Storage: true,
					Schema:  &apiext.CustomResourceValidation{},
				},
			},
		},
	}

	// Collect type markers.
	//
	// +k8s:crd-gen:resource:scope=,categories=,shortName=,plural=,subResources=
	// +k8s:crd-gen:printcolumn:name=,type=,jsonPath=,description=,format=,priority=
	tm := map[string][]string{}
	collectMarkers(t.SecondClosestCommentLines, tm)
	collectMarkers(t.CommentLines, tm)

	if len(tm) == 0 || len(tm["resource"]) == 0 {
		return nil
	}

	for _, res := range tm["resource"] {
		logger := logger.WithValues("markers", "resource")

		for mk, mv := range utils.ParseMarker(res) {
			switch mk {
			case "scope":
				var v apiext.ResourceScope
				err := json.Unmarshal([]byte(mv), &v)
				switch {
				case err != nil:
					logger.Error(err, "unmarshal scope", "value", mv)
				case !knownScopes.Has(v):
					logger.Error(nil, "invalid scope, select from known scopes", "value", mv)
				default:
					crd.Spec.Scope = v
				}
			case "categories":
				if err := json.Unmarshal([]byte(mv), &crd.Spec.Names.Categories); err != nil {
					logger.Error(nil, "unmarshal categories", "value", mv)
				}
			case "shortName":
				if err := json.Unmarshal([]byte(mv), &crd.Spec.Names.ShortNames); err != nil {
					logger.Error(nil, "unmarshal shortName", "value", mv)
				}
			case "plural":
				if err := json.Unmarshal([]byte(mv), &crd.Spec.Names.Plural); err != nil {
					logger.Error(nil, "unmarshal plural", "value", mv)
				}
			case "subResources":
				var subResources []string
				if err := json.Unmarshal([]byte(mv), &subResources); err != nil {
					logger.Error(nil, "unmarshal subResources", "value", mv)
					continue
				}

				if len(subResources) != 0 {
					crd.Spec.Versions[0].Subresources = &apiext.CustomResourceSubresources{}
				}

				for _, subres := range subResources {
					switch subres {
					case "status":
						crd.Spec.Versions[0].Subresources.Status = &apiext.CustomResourceSubresourceStatus{}
					case "scale":
						crd.Spec.Versions[0].Subresources.Scale = &apiext.CustomResourceSubresourceScale{}
					}
				}
			}
		}
	}
	if crd.Spec.Scope == "" {
		crd.Spec.Scope = apiext.NamespaceScoped
	}
	if crd.Spec.Names.Plural == "" {
		crd.Spec.Names.Plural = plural
	}

	for _, pc := range tm["printcolumn"] {
		logger := logger.WithValues("markers", "printcolumn")

		var pcd apiext.CustomResourceColumnDefinition

		for mk, mv := range utils.ParseMarker(pc) {
			switch mk {
			case "name":
				if err := json.Unmarshal([]byte(mv), &pcd.Name); err != nil {
					logger.Error(err, "unmarshal name", "value", mv)
				}
			case "type":
				var v string
				err := json.Unmarshal([]byte(mv), &v)
				switch {
				case err != nil:
					logger.Error(err, "unmarshal type", "value", mv)
				case !knownDataTypes.Has(v):
					logger.Error(nil, "invalid type, select from known data types", "value", mv)
				default:
					pcd.Type = v
				}
			case "format":
				var v string
				err := json.Unmarshal([]byte(mv), &v)
				switch {
				case err != nil:
					logger.Error(err, "unmarshal format", "value", mv)
				case !knownDataFormats.Has(v):
					logger.Error(nil, "invalid format, select from known formats", "value", mv)
				default:
					pcd.Format = v
				}
			case "jsonPath":
				if err := json.Unmarshal([]byte(mv), &pcd.JSONPath); err != nil {
					logger.Error(err, "unmarshal jsonPath", "value", mv)
				}
			case "priority":
				if err := json.Unmarshal([]byte(mv), &pcd.Priority); err != nil {
					logger.Error(err, "unmarshal priority", "value", mv)
				}
			}
		}

		if pcd.Name == "" || pcd.Type == "" || pcd.JSONPath == "" {
			logger.Error(nil, "invalid print column", "line", crdGenMarker+":printcolumn:"+pc)
			continue
		}

		crd.Spec.Versions[0].AdditionalPrinterColumns = append(crd.Spec.Versions[0].AdditionalPrinterColumns, pcd)
	}

	if props := schemeType(logger, nil, t); props != nil {
		props.Description = strings.Join(tm["comment"], "\n")
		crd.Spec.Versions[0].Schema.OpenAPIV3Schema = props
	}

	return &crd
}

var (
	bytesProps = apiext.JSONSchemaProps{
		Type:   "string",
		Format: "byte",
	}

	knownTypedProps = map[string]apiext.JSONSchemaProps{
		"k8s.io/apimachinery/pkg/runtime.RawExtension": {
			Type:                   "object",
			XPreserveUnknownFields: ptr.To(true),
		},
		"k8s.io/api/core/v1.Protocol": {
			Type:    "string",
			Default: &apiext.JSON{Raw: []byte(`"TCP"`)},
		},
		"k8s.io/apimachinery/pkg/apis/meta/v1.TypeMeta": {
			Type: "object",
			Properties: map[string]apiext.JSONSchemaProps{
				"apiVersion": {
					Type: "string",
				},
				"kind": {
					Type: "string",
				},
			},
		},
		"k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta": {
			Type: "object",
		},
		"k8s.io/apimachinery/pkg/apis/meta/v1.Fields": {
			Type:                 "object",
			AdditionalProperties: &apiext.JSONSchemaPropsOrBool{Allows: true},
		},
		"k8s.io/apimachinery/pkg/apis/meta/v1.FieldsV1": {
			Type:                 "object",
			AdditionalProperties: &apiext.JSONSchemaPropsOrBool{Allows: true},
		},
		"k8s.io/apimachinery/pkg/apis/meta/v1.Time": {
			Type:   "string",
			Format: "date-time",
		},
		"k8s.io/apimachinery/pkg/apis/meta/v1.MicroTime": {
			Type:   "string",
			Format: "date-time",
		},
		"k8s.io/apimachinery/pkg/apis/meta/v1.Duration": {
			Type:   "string",
			Format: "duration",
		},
		"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured.Unstructured": {
			Type: "object",
		},
		"k8s.io/apimachinery/pkg/api/resource.Quantity": {
			XIntOrString: true,
			AnyOf: []apiext.JSONSchemaProps{
				{Type: "integer"},
				{Type: "string"},
			},
			Pattern: "^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$",
		},
		"k8s.io/apimachinery/pkg/types.UID": {
			Type:   "string",
			Format: "uuid",
		},
		"k8s.io/apimachinery/pkg/util/intstr.IntOrString": {
			XIntOrString: true,
			AnyOf: []apiext.JSONSchemaProps{
				{Type: "integer"},
				{Type: "string"},
			},
		},
		"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1.JSON": {
			XPreserveUnknownFields: ptr.To(true),
		},
		"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1.JSON": {
			XPreserveUnknownFields: ptr.To(true),
		},
		"encoding/json.RawMessage": bytesProps,
	}
	knownEmbeddedResourceTypedProps = map[string]apiext.JSONSchemaProps{
		"k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta": {
			Type: "object",
			Properties: map[string]apiext.JSONSchemaProps{
				"name": {
					Type: "string",
				},
				"namespace": {
					Type: "string",
				},
				"annotations": {
					Type: "object",
					AdditionalProperties: &apiext.JSONSchemaPropsOrBool{
						Schema: &apiext.JSONSchemaProps{
							Type: "string",
						},
					},
				},
				"labels": {
					Type: "object",
					AdditionalProperties: &apiext.JSONSchemaPropsOrBool{
						Schema: &apiext.JSONSchemaProps{
							Type: "string",
						},
					},
				},
				"finalizers": {
					Type: "array",
					Items: &apiext.JSONSchemaPropsOrArray{
						Schema: &apiext.JSONSchemaProps{
							Type: "string",
						},
					},
				},
			},
		},
	}

	knownListTypes = sets.New("atomic", "map", "set")
	knownMapTypes  = sets.New("atomic", "granular")
	knownFormats   = sets.New("bsonobjectid", "uri", "email", "hostname", "ipv4", "ipv6",
		"cidr", "mac", "uuid", "uuid3", "uuid4", "uuid5", "isbn", "isbn10", "isbn13", "creditcard",
		"ssn", "hexcolor", "rgbcolor", "byte", "password", "date", "duration", "datetime")
	knownCelErrorReasons = sets.New(
		apiext.FieldValueRequired,
		apiext.FieldValueDuplicate,
		apiext.FieldValueInvalid,
		apiext.FieldValueForbidden)
)

// schemeType reflects the given type into a JSONSchemaProps .
func schemeType(logger klog.Logger, visited map[*types.Type]struct{}, t *types.Type) (props *apiext.JSONSchemaProps) {
	if t == nil {
		return nil
	}

	if visited == nil {
		visited = map[*types.Type]struct{}{}
	}

	if _, found := visited[t]; found {
		return &apiext.JSONSchemaProps{
			Ref: ptr.To(refer(t)),
		}
	}

	switch t.Kind {
	case types.Pointer:
		props = schemeType(logger, visited, t.Elem)
		if props != nil {
			props.Nullable = true
		}
	case types.Alias:
		if r, found := knownTypedProps[t.String()]; found {
			props = &r
		} else {
			props = schemeType(logger, visited, t.Underlying)
		}
	case types.Map:
		if t.Key != types.String {
			// Fallback to byte slices if map key is not string.
			props = ptr.To(bytesProps)
			logger.Error(nil, "invalid map key type, must be string, fallback to byte slice",
				"type", t.Key.String())
		} else if r := schemeType(logger, visited, t.Elem); r != nil {
			props = &apiext.JSONSchemaProps{
				Type: "object",
				AdditionalProperties: &apiext.JSONSchemaPropsOrBool{
					Allows: true,
					Schema: r,
				},
			}
		}

		if props != nil {
			props.Nullable = true
		}
	case types.Interface:
		props = &apiext.JSONSchemaProps{
			Type:                   "object",
			XPreserveUnknownFields: ptr.To(true),
		}
	case types.Array:
		props = schemeType(logger, visited, t.Elem)

		if props != nil {
			props.MinItems = ptr.To(t.Len)
			props.MaxItems = ptr.To(t.Len)
		}
	case types.Slice:
		if t.Elem == types.Byte && t.Len == 0 {
			props = ptr.To(bytesProps)
		} else if r := schemeType(logger, visited, t.Elem); r != nil {
			props = &apiext.JSONSchemaProps{
				Type: "array",
				Items: &apiext.JSONSchemaPropsOrArray{
					Schema: r,
				},
			}
		}

		if props != nil {
			props.Nullable = true

			if t.Len > 0 {
				props.MinItems = ptr.To(t.Len)
				props.MaxItems = ptr.To(t.Len)
			}
		}
	case types.Struct:
		if r, found := knownTypedProps[t.String()]; found {
			props = &r
		} else {
			props = &apiext.JSONSchemaProps{
				Type:       "object",
				Properties: map[string]apiext.JSONSchemaProps{},
			}

			for _, mem := range t.Members {
				var (
					name   = stringx.CamelizeDownFirst(mem.Name)
					inline bool
					hidden bool
				)
				if tg := parseStructTags(mem.Tags).Get("json"); tg.Available() {
					c, tg, s := tg.Next()
					if c != "" {
						name = c
					}

					for ; !s; c, tg, s = tg.Next() {
						switch c {
						case "inline":
							inline = true
						case "-":
							hidden = true
						}
					}
				}

				if hidden {
					continue
				}

				if t.Name.Package != "" && !t.IsAnonymousStruct() {
					visited[t] = struct{}{}
				}

				logger := logger.WithName(name)

				subProps := schemeType(logger, visited, mem.Type)
				if subProps == nil {
					continue
				}

				if inline {
					for k := range subProps.Properties {
						props.Properties[k] = subProps.Properties[k]
					}
					continue
				}

				if name != "status" {
					props.Required = append(props.Required, name)
				}

				// Collect member markers.
				//
				// +k8s:validation:default=
				// +k8s:validation:example=
				// +k8s:validation:enum=
				// +k8s:validation:maximum=
				// +k8s:validation:minimum=
				// +k8s:validation:exclusiveMaximum
				// +k8s:validation:exclusiveMinimum
				// +k8s:validation:multipleOf=
				// +k8s:validation:maxItems=
				// +k8s:validation:minItems=
				// +k8s:validation:uniqueItems
				// +k8s:validation:maxProperties=
				// +k8s:validation:minProperties=
				// +k8s:validation:maxLength=
				// +k8s:validation:minLength=
				// +k8s:validation:format=
				// +k8s:validation:pattern=
				// +k8s:validation:preserveUnknownFields
				// +k8s:validation:embeddedResource
				// +k8s:validation:cel[?]:rule=
				// +k8s:validation:cel[?]:rule>
				// +k8s:validation:cel[?]:message=
				// +k8s:validation:cel[?]:message>
				// +k8s:validation:cel[?]:messageExpression=
				// +k8s:validation:cel[?]:messageExpression>
				// +k8s:validation:cel[?]:reason=
				// +k8s:validation:cel[?]:fieldPath=
				// +k8s:validation:cel[?]:optionalOldSelf=
				//
				// +nullable=
				// +optional
				// +listType=
				// +listMapKey=
				// +mapType=
				// +default=
				mm := map[string][]string{}
				collectMarkers(mem.CommentLines, mm)

				for _, val := range mm["validation"] {
					for mk, mv := range utils.ParseMarker(val) {
						switch mk {
						case "default":
							if err := json.Unmarshal([]byte(mv), &subProps.Default); err != nil {
								logger.Error(err, "unmarshal default", "value", mv)
							}
						case "example":
							if err := json.Unmarshal([]byte(mv), &subProps.Example); err != nil {
								logger.Error(err, "unmarshal example", "value", mv)
							}
						case "enum":
							if err := json.Unmarshal([]byte(mv), &subProps.Enum); err != nil {
								logger.Error(err, "unmarshal enum", "value", mv)
							}
						case "maximum":
							var v float64
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal maximum", "value", mv)
							case !isNumber(mem.Type):
								logger.Error(nil, "unsupported maximum, must type with number")
							default:
								subProps.Maximum = ptr.To(v)
							}
						case "minimum":
							var v float64
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal minimum", "value", mv)
							case !isNumber(mem.Type):
								logger.Error(nil, "unsupported minimum, must type with number")
							default:
								subProps.Minimum = ptr.To(v)
							}
						case "exclusiveMaximum":
							if !isNumber(mem.Type) {
								logger.Error(nil, "unsupported exclusiveMaximum, must type with number")
							} else {
								subProps.ExclusiveMaximum = true
							}
						case "exclusiveMinimum":
							if !isNumber(mem.Type) {
								logger.Error(nil, "unsupported exclusiveMinimum, must type with number")
							} else {
								subProps.ExclusiveMinimum = true
							}
						case "multipleOf":
							var v float64
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal multipleOf", "value", mv)
							case v <= 0:
								logger.Error(nil, "invalid multipleOf, must be greater than 0",
									"value", mv)
							case !isNumber(mem.Type):
								logger.Error(nil, "unsupported multipleOf, must type with number")
							default:
								subProps.MultipleOf = ptr.To(v)
							}
						case "maxItems":
							var v int64
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal maxItems", "value", mv)
							case v <= 0:
								logger.Error(nil, "invalid maxItems, must be greater than 0",
									"value", mv)
							case !isSlice(mem.Type):
								logger.Error(nil, "unsupported maxItems, must type with slice")
							default:
								subProps.MaxItems = ptr.To(v)
							}
						case "minItems":
							var v int64
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal minItems", "value", mv)
							case v < 0:
								logger.Error(nil, "invalid minItems, must not be negative",
									"value", mv)
							case !isSlice(mem.Type):
								logger.Error(nil, "unsupported minItems, must type with slice")
							default:
								subProps.MinItems = ptr.To(v)
							}
						case "uniqueItems":
							if !isSlice(mem.Type) {
								logger.Error(nil, "unsupported uniqueItems, must type with slice")
							} else {
								subProps.UniqueItems = true
							}
						case "maxProperties":
							var v int64
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal maxProperties", "value", mv)
							case v <= 0:
								logger.Error(nil, "invalid maxProperties, must be greater than 0",
									"value", mv)
							case !isMap(mem.Type):
								logger.Error(nil, "unsupported maxProperties, must type with map")
							default:
								subProps.MaxProperties = ptr.To(v)
							}
						case "minProperties":
							var v int64
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal minProperties", "value", mv)
							case v < 0:
								logger.Error(nil, "invalid minProperties, must not be negative",
									"value", mv)
							case !isMap(mem.Type):
								logger.Error(nil, "unsupported maxProperties, must type with map")
							default:
								subProps.MinProperties = ptr.To(v)
							}
						case "maxLength":
							var v int64
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal maxLength", "value", mv)
							case v <= 0:
								logger.Error(nil, "invalid maxLength, must be greater than 0",
									"value", mv)
							case !isString(mem.Type):
								logger.Error(nil, "unsupported maxLength, must type with string")
							default:
								subProps.MaxLength = ptr.To(v)
							}
						case "minLength":
							var v int64
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal minLength", "value", mv)
							case v < 0:
								logger.Error(nil, "invalid minLength, must not be negative",
									"value", mv)
							case !isString(mem.Type):
								logger.Error(nil, "unsupported minLength, must type with string")
							default:
								subProps.MinLength = ptr.To(v)
							}
						case "format":
							var v string
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal format", "value", mv)
							case !knownFormats.Has(v):
								logger.Error(nil, "invalid format, select from known formats",
									"value", mv)
							case !isString(mem.Type):
								logger.Error(nil, "unsupported format, must type with string")
							default:
								subProps.Format = v
							}
						case "pattern":
							var v string
							if err := json.Unmarshal([]byte(mv), &v); err != nil {
								logger.Error(err, "unmarshal pattern", "value", mv)
							} else if _, err = regexp.Compile(v); err != nil {
								logger.Error(err, "invalid pattern", "value", mv)
							} else {
								subProps.Pattern = v
							}
						case "preserveUnknownFields":
							if !isStruct(mem.Type) {
								logger.Error(nil, "unsupported preserveUnknownFields, must type with struct")
							} else {
								subProps.XPreserveUnknownFields = ptr.To(true)
							}
						case "embeddedResource":
							if !isStruct(mem.Type) {
								logger.Error(nil, "unsupported embeddedResource, must type with struct")
							} else {
								subProps.XEmbeddedResource = true
								if r, found := knownEmbeddedResourceTypedProps[mem.Type.String()]; found {
									subProps.Properties = r.Properties
								}
							}
						case "nullable":
							if err := json.Unmarshal([]byte(mv), &subProps.Nullable); err != nil {
								logger.Error(err, "unmarshal nullable", "value", mv)
							}
						case "optional":
							props.Required = slices.DeleteFunc(props.Required,
								func(n string) bool {
									return n == name
								})
							if len(props.Required) == 0 {
								props.Required = nil
							}
						case "listType":
							var v string
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal listType", "value", mv)
							case !knownListTypes.Has(v):
								logger.Error(nil, "invalid listType, select from 'atomic', 'map' or 'set'",
									"value", mv)
							case !isSlice(mem.Type):
								logger.Error(nil, "unsupported listType, must type with slice")
							default:
								subProps.XListType = ptr.To(v)
							}
						case "listMapKey":
							var v string
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal listMapKey", "value", mv)
							case v == "":
								logger.Error(nil, "invalid listMapKey, must not be blank")
							case ptr.Deref(subProps.XListType, "") != "map":
								logger.Error(nil, "invalid listMapKey, as listType is not 'map'")
							default:
								subProps.XListMapKeys = append(subProps.XListMapKeys, v)
							}
						case "mapType":
							var v string
							err := json.Unmarshal([]byte(mv), &v)
							switch {
							case err != nil:
								logger.Error(err, "unmarshal mapType", "value", mv)
							case !knownMapTypes.Has(v):
								logger.Error(nil, "invalid mapType, select from 'map' atomic 'granular'")
							case !isMap(mem.Type):
								logger.Error(nil, "unsupported mapType, must type with map")
							default:
								subProps.XMapType = ptr.To(v)
							}
						}
					}
				}

				for _, cel := range mm["cel"] {
					kv := strings.SplitN(strings.TrimSpace(cel), ":", 2)
					if len(kv) != 2 || !strings.HasPrefix(kv[0], "cel[") {
						continue
					}

					idx, err := strconv.ParseInt(kv[0][4:len(kv[0])-1], 10, 64)
					if err != nil {
						logger.Error(err, "invalid cel index", "value", cel)
						continue
					}

					if idx < 0 {
						logger.Error(nil, "invalid cel index, must be non-negative", "value", cel)
						continue
					}

					switch length := int64(len(subProps.XValidations)); {
					case subProps.XValidations == nil && idx == 0:
						// Initialize.
						subProps.XValidations = apiext.ValidationRules{{}}
					case subProps.XValidations != nil && idx == length-1:
						// Update.
					case subProps.XValidations != nil && idx == length:
						// Increase.
						subProps.XValidations = append(subProps.XValidations, apiext.ValidationRule{})
					default:
						// Clean.
						subProps.XValidations = apiext.ValidationRules{}
						logger.Error(nil, "invalid cel index, must be incremental", "value", cel)
						continue
					}

					switch {
					case strings.HasPrefix(kv[1], "rule="):
						var v string
						if err := json.Unmarshal([]byte(kv[1][5:]), &v); err != nil {
							logger.Error(err, "unmarshal CEL rule", "value", cel)
						} else {
							subProps.XValidations[idx].Rule = v
						}
					case strings.HasPrefix(kv[1], "rule>"):
						v := strings.TrimSpace(kv[1][5:])
						if subProps.XValidations[idx].Rule == "" {
							subProps.XValidations[idx].Rule = v
						} else {
							subProps.XValidations[idx].Rule += "\n" + v
						}
					case strings.HasPrefix(kv[1], "message="):
						var v string
						if err := json.Unmarshal([]byte(kv[1][8:]), &v); err != nil {
							logger.Error(err, "unmarshal CEL message", "value", cel)
						} else {
							subProps.XValidations[idx].Message = v
						}
					case strings.HasPrefix(kv[1], "message>"):
						v := strings.TrimSpace(kv[1][8:])
						if subProps.XValidations[idx].Message == "" {
							subProps.XValidations[idx].Message = v
						} else {
							subProps.XValidations[idx].Message += "\n" + v
						}
					case strings.HasPrefix(kv[1], "messageExpression="):
						var v string
						if err := json.Unmarshal([]byte(kv[1][18:]), &v); err != nil {
							logger.Error(err, "unmarshal CEL messageExpression", "value", cel)
						} else {
							subProps.XValidations[idx].MessageExpression = v
						}
					case strings.HasPrefix(kv[1], "messageExpression>"):
						v := strings.TrimSpace(kv[1][18:])
						if subProps.XValidations[idx].MessageExpression == "" {
							subProps.XValidations[idx].MessageExpression = v
						} else {
							subProps.XValidations[idx].MessageExpression += "\n" + v
						}
					case strings.HasPrefix(kv[1], "reason="):
						var v string
						err := json.Unmarshal([]byte(kv[1][7:]), &v)
						switch {
						case err != nil:
							logger.Error(err, "unmarshal CEL reason", "value", cel)
						case !knownCelErrorReasons.Has(apiext.FieldValueErrorReason(v)):
							logger.Error(nil, "invalid CEL reason, select from known reasons", "value", cel)
						default:
							subProps.XValidations[idx].Reason = ptr.To(apiext.FieldValueErrorReason(v))
						}
					case strings.HasPrefix(kv[1], "fieldPath="):
						var v string
						if err := json.Unmarshal([]byte(kv[1][10:]), &v); err != nil {
							logger.Error(err, "unmarshal CEL fieldPath", "value", cel)
						} else {
							subProps.XValidations[idx].FieldPath = v
						}
					case strings.HasPrefix(kv[1], "optionalOldSelf="):
						var v bool
						if err := json.Unmarshal([]byte(kv[1][16:]), &v); err != nil {
							logger.Error(err, "unmarshal CEL optionalOldSelf", "value", cel)
						} else {
							subProps.XValidations[idx].OptionalOldSelf = ptr.To(v)
						}
					}
				}

				subProps.Description = strings.Join(mm["comment"], "\n")
				props.Properties[name] = *subProps
			}
		}
	}

	if t.IsPrimitive() {
		if r, found := knownTypedProps[t.String()]; found {
			props = &r
		} else {
			switch {
			case t == types.Bool:
				props = &apiext.JSONSchemaProps{
					Type: "boolean",
				}
			case t == types.String:
				props = &apiext.JSONSchemaProps{
					Type: "string",
				}
			case types.IsInteger(t):
				props = &apiext.JSONSchemaProps{
					Type: "integer",
				}

				switch t {
				case types.Int32, types.Uint16:
					props.Format = "int32"
				case types.Int64, types.Uint32:
					props.Format = "int64"
				}
			default:
				props = &apiext.JSONSchemaProps{
					Type: "number",
				}

				switch t {
				case types.Float32:
					props.Format = "float"
				case types.Float64:
					props.Format = "double"
				}
			}
		}
	}

	if props == nil {
		// Fallback to byte slices if type is not supported.
		props = ptr.To(bytesProps)
		logger.Info("unsupported type, fallback to byte slice", "type", t.String())
	}

	return props
}

// isNumber returns true if the given type is a number.
func isNumber(t *types.Type) bool {
	if t != nil {
		switch t.Kind {
		case types.Builtin:
			return types.IsInteger(t) || t == types.Float32 || t == types.Float64
		case types.Alias:
			return isNumber(t.Underlying)
		case types.Pointer:
			return isNumber(t.Elem)
		}
	}

	return false
}

// isString returns true if the given type is a string.
func isString(t *types.Type) bool {
	if t != nil {
		switch t.Kind {
		case types.Builtin:
			return t == types.String
		case types.Alias:
			return isString(t.Underlying)
		case types.Pointer:
			return isString(t.Elem)
		}
	}

	return false
}

// isSlice returns true if the given type is a slice.
func isSlice(t *types.Type) bool {
	if t != nil {
		switch t.Kind {
		case types.Slice:
			return true
		case types.Alias:
			return isSlice(t.Underlying)
		case types.Pointer:
			return isSlice(t.Elem)
		}
	}

	return false
}

// isMap returns true if the given type is a map.
func isMap(t *types.Type) bool {
	if t != nil {
		switch t.Kind {
		case types.Map:
			return true
		case types.Alias:
			return isMap(t.Underlying)
		case types.Pointer:
			return isMap(t.Elem)
		}
	}

	return false
}

// isStruct returns true if the given type is a struct.
func isStruct(t *types.Type) bool {
	if t != nil {
		switch t.Kind {
		case types.Struct:
			return true
		case types.Alias:
			return isStruct(t.Underlying)
		case types.Pointer:
			return isStruct(t.Elem)
		}
	}

	return false
}

// refer creates a definition link for the given package and type.
func refer(t *types.Type) string {
	r := t.Name.Name
	if t.Name.Package != "" {
		// Replace `/` with `~1` and `~` with `~0` according to JSONPointer escapes.
		r = strings.ReplaceAll(t.Name.Package, "/", "~1") + "~0" + t.Name.Name
	}
	return "#/definitions/" + r
}

const (
	groupMarker      = "+groupName="
	versionMarker    = "+versionName="
	nullableMarker   = "+nullable" // By default, this is nullable.
	defaultMarker    = "+default="
	optionalMarker   = "+optional"
	listTypeMarker   = "+listType="
	listMapKeyMarker = "+listMapKey="
	mapTypeMarker    = "+mapType="
	crdGenMarker     = "+k8s:crd-gen:"
	validationMarker = "+k8s:validation:" // Compatibility with openapi-gen.
)

// collectMarkers collects markers from the given comments into a map.
func collectMarkers(comments []string, into map[string][]string) {
	for _, line := range comments {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		switch {
		default:
			if !strings.HasPrefix(line, "+") {
				into["comment"] = append(into["comment"], line)
			}
		case strings.HasPrefix(line, groupMarker):
			if v := line[len(groupMarker):]; v != "" {
				into["group"] = append(into["group"], v)
			}
		case strings.HasPrefix(line, versionMarker):
			if v := line[len(versionMarker):]; v != "" {
				into["version"] = append(into["version"], v)
			}
		case strings.HasPrefix(line, crdGenMarker):
			kv := strings.SplitN(line[len(crdGenMarker):], ":", 2)
			if len(kv) == 2 {
				into[kv[0]] = append(into[kv[0]], kv[1])
			} else if len(kv) == 1 {
				into[kv[0]] = append(into[kv[0]], "")
			}
		case strings.HasPrefix(line, nullableMarker):
			if v := line[len(nullableMarker):]; v != "" {
				into["validation"] = append(into["validation"], "nullable"+v)
			} else {
				into["validation"] = append(into["validation"], "nullable=true")
			}
		case strings.HasPrefix(line, defaultMarker):
			if v := line[len(defaultMarker):]; v != "" {
				into["validation"] = append(into["validation"], "default="+v)
			}
		case line == optionalMarker:
			into["validation"] = append(into["validation"], "optional")
		case strings.HasPrefix(line, listTypeMarker):
			if v := line[len(listTypeMarker):]; v != "" {
				into["validation"] = append(into["validation"], "listType="+v)
			}
		case strings.HasPrefix(line, listMapKeyMarker):
			if v := line[len(listMapKeyMarker):]; v != "" {
				into["validation"] = append(into["validation"], "listMapKey="+v)
			}
		case strings.HasPrefix(line, mapTypeMarker):
			if v := line[len(mapTypeMarker):]; v != "" {
				into["validation"] = append(into["validation"], "mapType="+v)
			}
		case strings.HasPrefix(line, validationMarker):
			if v := line[len(validationMarker):]; v != "" {
				if strings.HasPrefix(v, "cel[") {
					into["cel"] = append(into["cel"], v)
					continue
				}
				into["validation"] = append(into["validation"], v)
			}
		}
	}
}

// parseStructTags returns the struct tags of the given type.
func parseStructTags(st string) StructTags {
	p, _ := reflect.ParseStructTags(st)
	return StructTags(p)
}

type (
	StructTags reflect.StructTags
	Tag        string
)

func (tags StructTags) String() string {
	if len(tags) == 0 {
		return ""
	}

	return reflect.StructTags(tags).String()
}

func (tags StructTags) Has(name string) bool {
	if len(tags) == 0 {
		return false
	}

	return reflect.StructTags(tags).Has(name)
}

func (tags StructTags) Get(name string) Tag {
	for _, tag := range tags {
		if tag.Name == name {
			return Tag(tag.Value)
		}
	}

	return ""
}

func (t Tag) String() string {
	return string(t)
}

func (t Tag) Available() bool {
	return string(t) != ""
}

func (t Tag) Next() (current string, remain Tag, stop bool) {
	b, a, _ := strings.Cut(string(t), ",")
	return b, Tag(a), b == "" && a == ""
}
