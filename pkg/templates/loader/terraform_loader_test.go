package loader

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/templates/openapi"
)

var mockInfo = &openapi3.Info{
	Title: "OpenAPI schema for template dev-template",
}

func TestLoadTerraformSchema(t *testing.T) {
	var length uint64 = 3
	testCases := []struct {
		name           string
		input          string
		expectedOutput *types.TemplateVersionSchema
		expectedError  bool
	}{
		{
			name:           "Invalid",
			input:          "testdata/invalid",
			expectedOutput: &types.TemplateVersionSchema{},
			expectedError:  true,
		},
		{
			name:           "Invalid variable type",
			input:          "testdata/invalid_variable_type",
			expectedOutput: &types.TemplateVersionSchema{},
			expectedError:  true,
		},
		{
			name:  "With README.md",
			input: "testdata/with_readme",
			expectedOutput: &types.TemplateVersionSchema{
				TemplateVersionSchemaData: types.TemplateVersionSchemaData{
					Readme: "# test readme",
				},
			},
			expectedError: false,
		},
		{
			name:  "Any variable",
			input: "testdata/any_variable",
			expectedOutput: &types.TemplateVersionSchema{
				Schema: types.Schema{
					OpenAPISchema: &openapi3.T{
						OpenAPI: openapi.OpenAPIVersion,
						Info:    mockInfo,
						Components: &openapi3.Components{
							Schemas: map[string]*openapi3.SchemaRef{
								"variables": {
									Value: &openapi3.Schema{
										Required: []string{
											"list_object_with_any_default",
											"map_object_with_any_default",
											"object_with_any_default",
										},
										Type: openapi3.TypeObject,
										Properties: map[string]*openapi3.SchemaRef{
											"list_any_with_default": {
												Value: &openapi3.Schema{
													Title: "List Any With Default",
													Type:  openapi3.TypeArray,
													Items: &openapi3.SchemaRef{
														Value: &openapi3.Schema{
															Type: openapi3.TypeObject,
															Extensions: openapi.NewExt().
																WithOriginalType(cty.DynamicPseudoType).
																WithUIColSpan(12).
																Export(),
															Properties: map[string]*openapi3.SchemaRef{},
														},
													},
													Default: []any{map[string]any{"name": "default-name"}},
													Extensions: openapi.NewExt().
														WithOriginalType(cty.List(cty.DynamicPseudoType)).
														WithUIGroup("Basic").
														WithUIColSpan(12).
														WithUIOrder(1).
														Export(),
												},
											},
											"map_any_with_default": {
												Value: &openapi3.Schema{
													Title:   "Map Any With Default",
													Type:    openapi3.TypeObject,
													Default: map[string]any{"name": "default-name"},
													AdditionalProperties: openapi3.AdditionalProperties{
														Schema: &openapi3.SchemaRef{
															Value: &openapi3.Schema{
																Type:       openapi3.TypeObject,
																Properties: map[string]*openapi3.SchemaRef{},
																Extensions: openapi.NewExt().
																	WithOriginalType(cty.DynamicPseudoType).
																	WithUIColSpan(12).
																	Export(),
															},
														},
													},
													Properties: map[string]*openapi3.SchemaRef{},
													Extensions: openapi.NewExt().
														WithOriginalType(cty.Map(cty.DynamicPseudoType)).
														WithUIGroup("Basic").
														WithUIColSpan(12).
														WithUIOrder(2).
														Export(),
												},
											},
											"list_map_any_with_default": {
												Value: &openapi3.Schema{
													Title:   "List Map Any With Default",
													Type:    openapi3.TypeArray,
													Default: []any{map[string]any{"name": "default-name"}},
													Items: &openapi3.SchemaRef{
														Value: &openapi3.Schema{
															Type:       openapi3.TypeObject,
															Properties: map[string]*openapi3.SchemaRef{},
															AdditionalProperties: openapi3.AdditionalProperties{
																Schema: &openapi3.SchemaRef{
																	Value: &openapi3.Schema{
																		Type:       openapi3.TypeObject,
																		Properties: map[string]*openapi3.SchemaRef{},
																		Extensions: openapi.NewExt().
																			WithOriginalType(cty.DynamicPseudoType).
																			WithUIColSpan(12).
																			Export(),
																	},
																},
															},
															Extensions: openapi.NewExt().
																WithOriginalType(cty.Map(cty.DynamicPseudoType)).
																WithUIColSpan(12).
																Export(),
														},
													},
													Extensions: openapi.NewExt().
														WithOriginalType(cty.List(cty.Map(cty.DynamicPseudoType))).
														WithUIGroup("Basic").
														WithUIOrder(3).
														WithUIColSpan(12).
														Export(),
												},
											},
											"object_with_any_default": {
												Value: &openapi3.Schema{
													Title: "Object With Any Default",
													Type:  openapi3.TypeObject,
													Properties: map[string]*openapi3.SchemaRef{
														"any_data": {
															Value: &openapi3.Schema{
																Title: "Any Data",
																Type:  openapi3.TypeObject,
																Default: map[string]any{
																	"headers": map[string]any{
																		"X-Forwarded-Proto": "https",
																	},
																	"port": float64(80),
																},
																Properties: map[string]*openapi3.SchemaRef{},
																Extensions: openapi.NewExt().
																	WithOriginalType(cty.DynamicPseudoType).
																	WithUIOrder(1).
																	WithUIColSpan(12).
																	Export(),
															},
														},
													},
													Extensions: openapi.NewExt().
														WithOriginalType(
															cty.ObjectWithOptionalAttrs(
																map[string]cty.Type{
																	"any_data": cty.DynamicPseudoType,
																},
																[]string{"any_data"})).
														WithUIGroup("Basic").
														WithUIOrder(4).
														WithUIColSpan(12).
														Export(),
												},
											},
											"list_object_with_any_default": {
												Value: &openapi3.Schema{
													Title: "List Object With Any Default",
													Type:  openapi3.TypeArray,
													Items: &openapi3.SchemaRef{
														Value: &openapi3.Schema{
															Type: openapi3.TypeObject,
															Properties: map[string]*openapi3.SchemaRef{
																"any_data": {
																	Value: &openapi3.Schema{
																		Title: "Any Data",
																		Type:  openapi3.TypeObject,
																		Default: map[string]any{
																			"headers": map[string]any{
																				"X-Forwarded-Proto": "https",
																			},
																			"port": float64(80),
																		},
																		Properties: map[string]*openapi3.SchemaRef{},
																		Extensions: openapi.NewExt().
																			WithOriginalType(cty.DynamicPseudoType).
																			WithUIOrder(1).
																			WithUIColSpan(12).
																			Export(),
																	},
																},
															},
															Extensions: openapi.NewExt().
																WithOriginalType(
																	cty.ObjectWithOptionalAttrs(
																		map[string]cty.Type{
																			"any_data": cty.DynamicPseudoType,
																		},
																		[]string{"any_data"})).
																WithUIColSpan(12).
																Export(),
														},
													},
													Extensions: openapi.NewExt().
														WithOriginalType(
															cty.List(
																cty.ObjectWithOptionalAttrs(
																	map[string]cty.Type{
																		"any_data": cty.DynamicPseudoType,
																	},
																	[]string{"any_data"}))).
														WithUIGroup("Basic").
														WithUIOrder(5).
														WithUIColSpan(12).
														Export(),
												},
											},
											"map_object_with_any_default": {
												Value: &openapi3.Schema{
													Title:      "Map Object With Any Default",
													Type:       openapi3.TypeObject,
													Properties: map[string]*openapi3.SchemaRef{},
													AdditionalProperties: openapi3.AdditionalProperties{
														Schema: &openapi3.SchemaRef{
															Value: &openapi3.Schema{
																Type: openapi3.TypeObject,
																Properties: map[string]*openapi3.SchemaRef{
																	"any_data": {
																		Value: &openapi3.Schema{
																			Title: "Any Data",
																			Type:  openapi3.TypeObject,
																			Default: map[string]any{
																				"headers": map[string]any{
																					"X-Forwarded-Proto": "https",
																				},
																				"port": float64(80),
																			},
																			Properties: map[string]*openapi3.SchemaRef{},
																			Extensions: openapi.NewExt().
																				WithOriginalType(cty.DynamicPseudoType).
																				WithUIOrder(1).
																				WithUIColSpan(12).
																				Export(),
																		},
																	},
																},
																Extensions: openapi.NewExt().
																	WithOriginalType(
																		cty.ObjectWithOptionalAttrs(
																			map[string]cty.Type{
																				"any_data": cty.DynamicPseudoType,
																			},
																			[]string{"any_data"})).
																	WithUIColSpan(12).
																	Export(),
															},
														},
													},
													Extensions: openapi.NewExt().
														WithOriginalType(
															cty.Map(
																cty.ObjectWithOptionalAttrs(
																	map[string]cty.Type{
																		"any_data": cty.DynamicPseudoType,
																	},
																	[]string{"any_data"}))).
														WithUIGroup("Basic").
														WithUIOrder(6).
														WithUIColSpan(12).
														Export(),
												},
											},
										},
										Extensions: openapi.NewExt().
											WithOriginalVariablesSequence(
												[]string{
													"list_any_with_default",
													"map_any_with_default",
													"list_map_any_with_default",
													"object_with_any_default",
													"list_object_with_any_default",
													"map_object_with_any_default",
												}).
											Export(),
									},
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "With output",
			input: "testdata/with_output",
			expectedOutput: &types.TemplateVersionSchema{
				Schema: types.Schema{
					OpenAPISchema: &openapi3.T{
						OpenAPI: openapi.OpenAPIVersion,
						Info:    mockInfo,
						Components: &openapi3.Components{
							Schemas: map[string]*openapi3.SchemaRef{
								"outputs": {
									Value: &openapi3.Schema{
										Type: "object",
										Properties: map[string]*openapi3.SchemaRef{
											"first": {
												Value: &openapi3.Schema{
													Title:       "First",
													Type:        openapi3.TypeObject,
													Description: "The first output.",
													Extensions: openapi.NewExt().
														WithOriginalType(cty.DynamicPseudoType).
														WithOriginalValueExpression([]byte("1")).
														WithUIOrder(1).
														WithUIColSpan(12).
														Export(),
													Properties: map[string]*openapi3.SchemaRef{},
												},
											},
											"second": {
												Value: &openapi3.Schema{
													Title:       "Second",
													Type:        openapi3.TypeObject,
													Description: "The second output.",
													Extensions: openapi.NewExt().
														WithOriginalType(cty.DynamicPseudoType).
														WithOriginalValueExpression([]byte("2")).
														WithUIOrder(2).
														WithUIColSpan(12).
														Export(),
													WriteOnly:  true,
													Properties: map[string]*openapi3.SchemaRef{},
												},
											},
										},
										Extensions: openapi.NewExt().
											Export(),
									},
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "With variable",
			input: "testdata/with_variable",
			expectedOutput: &types.TemplateVersionSchema{
				Schema: types.Schema{
					OpenAPISchema: &openapi3.T{
						OpenAPI: openapi.OpenAPIVersion,
						Info:    mockInfo,
						Components: &openapi3.Components{
							Schemas: map[string]*openapi3.SchemaRef{
								"variables": {
									Value: &openapi3.Schema{
										Type: openapi3.TypeObject,
										Properties: map[string]*openapi3.SchemaRef{
											"foo": {
												Value: &openapi3.Schema{
													Title:   "Foo",
													Type:    openapi3.TypeString,
													Default: "bar",
													Extensions: openapi.NewExt().
														WithOriginalType(cty.String).
														WithUIGroup("Basic").
														WithUIOrder(1).
														Export(),
												},
											},
										},
										Extensions: openapi.NewExt().
											WithOriginalVariablesSequence([]string{"foo"}).
											WithUIGroupOrder("Basic").
											Export(),
									},
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "With description",
			input: "testdata/with_description",
			expectedOutput: &types.TemplateVersionSchema{
				Schema: types.Schema{
					OpenAPISchema: &openapi3.T{
						OpenAPI: openapi.OpenAPIVersion,
						Info:    mockInfo,
						Components: &openapi3.Components{
							Schemas: map[string]*openapi3.SchemaRef{
								"variables": {
									Value: &openapi3.Schema{
										Type: "object",
										Properties: map[string]*openapi3.SchemaRef{
											"foo": {
												Value: &openapi3.Schema{
													Title:       "Foo",
													Type:        "string",
													Description: "description of foo.",
													Default:     "bar",
													Extensions: openapi.NewExt().
														WithOriginalType(cty.String).
														WithUIGroup("Basic").
														WithUIOrder(1).
														Export(),
												},
											},
										},
										Extensions: openapi.NewExt().
											WithOriginalVariablesSequence([]string{"foo"}).
											WithUIGroupOrder("Basic").
											Export(),
									},
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "With schema.yaml",
			input: "testdata/full_schema",
			expectedOutput: &types.TemplateVersionSchema{
				TemplateVersionSchemaData: types.TemplateVersionSchemaData{
					Readme: "# test readme",
					RequiredProviders: []types.ProviderRequirement{
						{
							Name: "mycloud",
							ProviderRequirement: &tfconfig.ProviderRequirement{
								Source:             "mycorp/mycloud",
								VersionConstraints: []string{"~> 1.0"},
							},
						},
						{
							Name: "null",
							ProviderRequirement: &tfconfig.ProviderRequirement{
								Source: "hashicorp/null",
							},
						},
					},
				},
				Schema: types.Schema{
					OpenAPISchema: &openapi3.T{
						OpenAPI: openapi.OpenAPIVersion,
						Info:    mockInfo,
						Components: &openapi3.Components{
							Schemas: map[string]*openapi3.SchemaRef{
								"variables": {
									Value: &openapi3.Schema{
										Type: "object",
										Properties: map[string]*openapi3.SchemaRef{
											"foo": {
												Value: &openapi3.Schema{
													Title:   "Foo Label",
													Type:    openapi3.TypeString,
													Default: "foo",
													Enum: []any{
														"F1",
														"F2",
														"F3",
													},
													Extensions: openapi.NewExt().
														WithOriginalType(cty.String).
														WithUIGroup("Test Group").
														WithUIOrder(1).
														Export(),
												},
											},
											"bar": {
												Value: &openapi3.Schema{
													Title:       "Bar Label",
													Type:        "string",
													Default:     "bar",
													Description: "description of bar.",
													Enum: []any{
														"B1",
														"B2",
														"B3",
													},
													Extensions: openapi.NewExt().
														WithOriginalType(cty.String).
														WithUIGroup("Test Group").
														WithUIShowIf("foo=F1").
														WithUIOrder(2).
														Export(),
												},
											},
											"thee": {
												Value: &openapi3.Schema{
													Title:   "Thee",
													Type:    openapi3.TypeString,
													Default: "foo",
													Enum: []any{
														"F1",
														"F2",
														"F3",
													},
													Extensions: openapi.NewExt().
														WithOriginalType(cty.String).
														WithUIGroup("Test Group").
														WithUIOrder(3).
														Export(),
												},
											},
											"number_options_var": {
												Value: &openapi3.Schema{
													Title:   "number_options_var",
													Type:    "number",
													Default: float64(1),
													Enum: []any{
														float64(1),
														float64(2),
														float64(3),
													},
													Extensions: openapi.NewExt().
														WithOriginalType(cty.Number).
														WithUIGroup("Basic").
														WithUIOrder(4).
														Export(),
												},
											},
											"subgroup1_1": {
												Value: &openapi3.Schema{
													Title:   "Subgroup1_1 Label",
													Type:    openapi3.TypeString,
													Default: "subgroup1_1",
													Extensions: openapi.NewExt().
														WithOriginalType(cty.String).
														WithUIGroup("Test Subgroup/Subgroup 1").
														WithUIOrder(5).
														Export(),
												},
											},
											"subgroup1_2": {
												Value: &openapi3.Schema{
													Title:   "Subgroup1_2 Label",
													Type:    openapi3.TypeString,
													Default: "subgroup1_2",
													Extensions: openapi.NewExt().
														WithOriginalType(cty.String).
														WithUIGroup("Test Subgroup/Subgroup 1").
														WithUIOrder(6).
														Export(),
												},
											},
											"subgroup2_1": {
												Value: &openapi3.Schema{
													Title:   "Subgroup2_1 Label",
													Type:    openapi3.TypeString,
													Default: "subgroup2_1",
													Extensions: openapi.NewExt().
														WithOriginalType(cty.String).
														WithUIGroup("Test Subgroup/Subgroup 2").
														WithUIWidget("Input").
														WithUIOrder(7).
														Export(),
												},
											},
											"subgroup2_1_hidden": {
												Value: &openapi3.Schema{
													Title:   "subgroup2_1_hidden",
													Type:    openapi3.TypeString,
													Default: "",
													Extensions: openapi.NewExt().
														WithOriginalType(cty.String).
														WithUIGroup("Test Subgroup/Subgroup 2").
														WithUIHidden().
														WithUIOrder(8).
														Export(),
												},
											},
										},
										Extensions: openapi.NewExt().
											WithOriginalVariablesSequence([]string{
												"foo",
												"bar",
												"thee",
												"number_options_var",
												"subgroup1_1",
												"subgroup1_2",
												"subgroup2_1",
												"subgroup2_1_hidden",
											}).
											WithUIGroupOrder(
												"Test Group",
												"Basic",
												"Test Subgroup/Subgroup 1",
												"Test Subgroup/Subgroup 2",
											).
											Export(),
									},
								},
								"outputs": {
									Value: &openapi3.Schema{
										Type: "object",
										Properties: map[string]*openapi3.SchemaRef{
											"first": {
												Value: &openapi3.Schema{
													Title:       "First",
													Type:        openapi3.TypeObject,
													Description: "The first output.",
													Extensions: openapi.NewExt().
														WithOriginalType(cty.DynamicPseudoType).
														WithOriginalValueExpression([]byte("null_resource.test.id")).
														WithUIOrder(1).
														WithUIColSpan(12).
														Export(),
													Properties: map[string]*openapi3.SchemaRef{},
												},
											},
											"second": {
												Value: &openapi3.Schema{
													Title:       "Second",
													Type:        openapi3.TypeObject,
													Description: "The second output.",
													WriteOnly:   true,
													Extensions: openapi.NewExt().
														WithOriginalType(cty.DynamicPseudoType).
														WithOriginalValueExpression([]byte(`"some value"`)).
														WithUIOrder(2).
														WithUIColSpan(12).
														Export(),
													Properties: map[string]*openapi3.SchemaRef{},
												},
											},
										},
										Extensions: openapi.NewExt().
											Export(),
									},
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "Complex variable",
			input: "testdata/complex_variable",
			expectedOutput: &types.TemplateVersionSchema{
				Schema: types.Schema{
					OpenAPISchema: &openapi3.T{
						Info: mockInfo,
						Components: &openapi3.Components{
							Schemas: map[string]*openapi3.SchemaRef{
								"variables": {
									Value: &openapi3.Schema{
										Type: "object",
										Properties: map[string]*openapi3.SchemaRef{
											"any": {
												Value: &openapi3.Schema{
													Title: "Any",
													Type:  openapi3.TypeObject,
													Extensions: openapi.NewExt().
														WithOriginalType(cty.DynamicPseudoType).
														WithUIGroup("Basic").
														WithUIOrder(1).
														WithUIColSpan(12).
														Export(),
													Properties: map[string]*openapi3.SchemaRef{},
												},
											},
											"any_map": {
												Value: &openapi3.Schema{
													Title:      "Any Map",
													Type:       openapi3.TypeObject,
													Properties: map[string]*openapi3.SchemaRef{},
													AdditionalProperties: openapi3.AdditionalProperties{
														Schema: &openapi3.SchemaRef{
															Value: &openapi3.Schema{
																Type: openapi3.TypeObject,
																Extensions: openapi.NewExt().
																	WithOriginalType(cty.DynamicPseudoType).
																	WithUIColSpan(12).
																	Export(),
																Properties: map[string]*openapi3.SchemaRef{},
															},
														},
													},
													Extensions: openapi.NewExt().
														WithOriginalType(cty.Map(cty.DynamicPseudoType)).
														WithUIGroup("Basic").
														WithUIOrder(2).
														WithUIColSpan(12).
														Export(),
												},
											},
											"string_map": {
												Value: &openapi3.Schema{
													Title: "String Map",
													Type:  openapi3.TypeObject,
													Default: map[string]any{
														"a": "a",
														"b": "1",
														"c": "true",
													},
													Properties: map[string]*openapi3.SchemaRef{},
													AdditionalProperties: openapi3.AdditionalProperties{
														Schema: &openapi3.SchemaRef{
															Value: &openapi3.Schema{
																Type: openapi3.TypeString,
																Extensions: openapi.NewExt().
																	WithOriginalType(cty.String).
																	Export(),
															},
														},
													},
													Extensions: openapi.NewExt().
														WithOriginalType(cty.Map(cty.String)).
														WithUIGroup("Basic").
														WithUIOrder(3).
														WithUIColSpan(12).
														Export(),
												},
											},
											"string_slice": {
												Value: &openapi3.Schema{
													Title:   "String Slice",
													Type:    "array",
													Default: []any{"x", "y", "z"},
													Items: &openapi3.SchemaRef{
														Value: &openapi3.Schema{
															Type: openapi3.TypeString,
															Extensions: openapi.NewExt().
																WithOriginalType(cty.String).
																Export(),
														},
													},
													Extensions: openapi.NewExt().
														WithOriginalType(cty.List(cty.String)).
														WithUIGroup("Basic").
														WithUIOrder(4).
														WithUIColSpan(12).
														Export(),
												},
											},
											"object": {
												Value: &openapi3.Schema{
													Title: "Object",
													Type:  openapi3.TypeObject,
													Default: map[string]any{
														"a": "a",
														"b": float64(1),
														"c": true,
													},
													Properties: map[string]*openapi3.SchemaRef{
														"a": {
															Value: &openapi3.Schema{
																Title: "A",
																Type:  openapi3.TypeString,
																Extensions: openapi.NewExt().
																	WithOriginalType(cty.String).
																	WithUIOrder(1).
																	Export(),
															},
														},
														"b": {
															Value: &openapi3.Schema{
																Title: "B",
																Type:  openapi3.TypeNumber,
																Extensions: openapi.NewExt().
																	WithOriginalType(cty.Number).
																	WithUIOrder(2).
																	Export(),
															},
														},
														"c": {
															Value: &openapi3.Schema{
																Title: "C",
																Type:  openapi3.TypeBoolean,
																Extensions: openapi.NewExt().
																	WithOriginalType(cty.Bool).
																	WithUIOrder(3).
																	Export(),
															},
														},
													},
													Required: []string{"a", "b", "c"},
													Extensions: openapi.NewExt().
														WithOriginalType(
															cty.Object(
																map[string]cty.Type{
																	"a": cty.String,
																	"b": cty.Number,
																	"c": cty.Bool,
																},
															)).
														WithUIGroup("Basic").
														WithUIOrder(5).
														WithUIColSpan(12).
														Export(),
												},
											},
											"object_nested": {
												Value: &openapi3.Schema{
													Title: "Object Nested",
													Type:  openapi3.TypeObject,
													Default: map[string]any{
														"a": "a",
														"b": []any{
															map[string]any{
																"c": true,
															},
														},
													},
													Required: []string{"a", "b"},
													Properties: map[string]*openapi3.SchemaRef{
														"a": {
															Value: &openapi3.Schema{
																Title: "A",
																Type:  openapi3.TypeString,
																Extensions: openapi.NewExt().
																	WithOriginalType(cty.String).
																	WithUIOrder(1).
																	Export(),
															},
														},
														"b": {
															Value: &openapi3.Schema{
																Title: "B",
																Type:  openapi3.TypeArray,
																Items: &openapi3.SchemaRef{
																	Value: &openapi3.Schema{
																		Type: openapi3.TypeObject,
																		Properties: map[string]*openapi3.SchemaRef{
																			"c": {
																				Value: &openapi3.Schema{
																					Title: "C",
																					Type:  openapi3.TypeBoolean,
																					Extensions: openapi.NewExt().
																						WithOriginalType(cty.Bool).
																						WithUIOrder(1).
																						Export(),
																				},
																			},
																		},
																		Required: []string{"c"},
																		Extensions: openapi.NewExt().
																			WithOriginalType(
																				cty.Object(map[string]cty.Type{
																					"c": cty.Bool,
																				})).
																			WithUIColSpan(12).
																			Export(),
																	},
																},
																Extensions: openapi.NewExt().WithOriginalType(
																	cty.List(
																		cty.Object(map[string]cty.Type{
																			"c": cty.Bool,
																		}))).
																	WithUIOrder(2).
																	WithUIColSpan(12).
																	Export(),
															},
														},
													},
													Extensions: openapi.NewExt().
														WithOriginalType(
															cty.Object(map[string]cty.Type{
																"a": cty.String,
																"b": cty.List(
																	cty.Object(map[string]cty.Type{
																		"c": cty.Bool,
																	})),
															})).
														WithUIGroup("Basic").
														WithUIOrder(6).
														WithUIColSpan(12).
														Export(),
												},
											},
											"list_object": {
												Value: &openapi3.Schema{
													Title: "List Object",
													Type:  openapi3.TypeArray,
													Items: &openapi3.SchemaRef{
														Value: &openapi3.Schema{
															Type: openapi3.TypeObject,
															Properties: map[string]*openapi3.SchemaRef{
																"a": {
																	Value: &openapi3.Schema{
																		Title: "A",
																		Type:  openapi3.TypeString,
																		Extensions: openapi.NewExt().
																			WithOriginalType(cty.String).
																			WithUIOrder(1).
																			Export(),
																	},
																},
																"b": {
																	Value: &openapi3.Schema{
																		Title: "B",
																		Type:  openapi3.TypeNumber,
																		Extensions: openapi.NewExt().
																			WithOriginalType(cty.Number).
																			WithUIOrder(2).
																			Export(),
																	},
																},
																"c": {
																	Value: &openapi3.Schema{
																		Title: "C",
																		Type:  openapi3.TypeBoolean,
																		Extensions: openapi.NewExt().
																			WithOriginalType(cty.Bool).
																			WithUIOrder(3).
																			Export(),
																	},
																},
															},
															Required: []string{"a", "b", "c"},
															Extensions: openapi.NewExt().
																WithOriginalType(
																	cty.Object(
																		map[string]cty.Type{
																			"a": cty.String,
																			"b": cty.Number,
																			"c": cty.Bool,
																		})).
																WithUIColSpan(12).
																Export(),
														},
													},
													Extensions: openapi.NewExt().
														WithOriginalType(
															cty.List(
																cty.Object(
																	map[string]cty.Type{
																		"a": cty.String,
																		"b": cty.Number,
																		"c": cty.Bool,
																	}))).
														WithUIGroup("Basic").
														WithUIOrder(7).
														WithUIColSpan(12).
														Export(),
												},
											},
											"tuple": {
												Value: &openapi3.Schema{
													Title:     "Tuple",
													Type:      "array",
													MaxLength: &length,
													MinLength: 3,
													Items:     openapi3.NewObjectSchema().NewRef(),
													Extensions: openapi.NewExt().
														WithOriginalType(
															cty.Tuple(
																[]cty.Type{
																	cty.String, cty.Bool, cty.Number,
																})).
														WithUIGroup("Basic").
														WithUIOrder(8).
														WithUIColSpan(12).
														Export(),
												},
											},
										},
										Required: []string{
											"any",
											"list_object",
											"tuple",
										},
										Extensions: openapi.NewExt().
											WithOriginalVariablesSequence([]string{
												"any",
												"any_map",
												"string_map",
												"string_slice",
												"object",
												"object_nested",
												"list_object",
												"tuple",
											}).
											WithUIGroupOrder("Basic").
											Export(),
									},
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loader := NewTerraformLoader()

			actualOutput, actualError := loader.Load(tc.input, "dev-template", ModeSchemaFile)
			if tc.expectedError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}

			if actualOutput == nil {
				return
			}

			assert.Equal(t, tc.expectedOutput.TemplateVersionSchemaData, actualOutput.TemplateVersionSchemaData)

			if tc.expectedOutput.OpenAPISchema != nil && tc.expectedOutput.OpenAPISchema.Components != nil {
				assert.Equal(
					t,
					tc.expectedOutput.OpenAPISchema.Components.Schemas["variables"],
					actualOutput.OpenAPISchema.Components.Schemas["variables"],
				)
				assert.Equal(
					t,
					tc.expectedOutput.OpenAPISchema.Components.Schemas["outputs"],
					actualOutput.OpenAPISchema.Components.Schemas["outputs"],
				)
			}
		})
	}
}

func TestOutputValueExpression(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput map[string][]byte
		expectedError  bool
	}{
		{
			name:  "Get output value expression",
			input: "testdata/with_output/main.tf",
			expectedOutput: map[string][]byte{
				"first":  []byte("1"),
				"second": []byte("2"),
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			files := sets.Set[string]{}
			files.Insert(tc.input)

			actualOutput, actualError := getOutputValues(files)
			if tc.expectedError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}

			assert.Equal(t, tc.expectedOutput, actualOutput)
		})
	}
}
