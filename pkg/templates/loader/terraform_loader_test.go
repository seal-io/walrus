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
	Title:   "OpenAPI schema for template dev-template",
	Version: "dev",
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
													Title:       "first",
													Description: "The first output.",
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.DynamicPseudoType).
														SetOriginalValueExpression([]byte("1")).
														Export(),
												},
											},
											"second": {
												Value: &openapi3.Schema{
													Title:       "second",
													Description: "The second output.",
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.DynamicPseudoType).
														SetOriginalValueExpression([]byte("2")).
														Export(),
													WriteOnly: true,
												},
											},
										},
										Extensions: openapi.NewExt(nil).
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
													Title:   "foo",
													Type:    openapi3.TypeString,
													Default: "bar",
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.String).
														Export(),
												},
											},
										},
										Extensions: openapi.NewExt(nil).
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
													Title:       "foo",
													Type:        "string",
													Description: "description of foo.",
													Default:     "bar",
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.String).
														Export(),
												},
											},
										},
										Extensions: openapi.NewExt(nil).
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
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.String).
														SetUIGroup("Test Group").
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
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.String).
														SetUIGroup("Test Group").
														SetUIShowIf("foo=F1").
														Export(),
												},
											},
											"thee": {
												Value: &openapi3.Schema{
													Title:   "thee",
													Type:    openapi3.TypeString,
													Default: "foo",
													Enum: []any{
														"F1",
														"F2",
														"F3",
													},
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.String).
														SetUIGroup("Test Group").
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
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.Number).
														SetUIGroup("Basic").
														Export(),
												},
											},
											"subgroup1_1": {
												Value: &openapi3.Schema{
													Title:   "Subgroup1_1 Label",
													Type:    openapi3.TypeString,
													Default: "subgroup1_1",
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.String).
														SetUIGroup("Test Subgroup/Subgroup 1").
														Export(),
												},
											},
											"subgroup1_2": {
												Value: &openapi3.Schema{
													Title:   "Subgroup1_2 Label",
													Type:    openapi3.TypeString,
													Default: "subgroup1_2",
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.String).
														SetUIGroup("Test Subgroup/Subgroup 1").
														Export(),
												},
											},
											"subgroup2_1": {
												Value: &openapi3.Schema{
													Title:   "Subgroup2_1 Label",
													Type:    openapi3.TypeString,
													Default: "subgroup2_1",
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.String).
														SetUIGroup("Test Subgroup/Subgroup 2").
														SetUIWidget("Input").
														Export(),
												},
											},
											"subgroup2_1_hidden": {
												Value: &openapi3.Schema{
													Title:   "subgroup2_1_hidden",
													Type:    openapi3.TypeString,
													Default: "",
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.String).
														SetUIGroup("Test Subgroup/Subgroup 2").
														SetUIHidden().
														Export(),
												},
											},
										},
										Extensions: map[string]any{},
									},
								},
								"outputs": {
									Value: &openapi3.Schema{
										Type: "object",
										Properties: map[string]*openapi3.SchemaRef{
											"first": {
												Value: &openapi3.Schema{
													Title:       "first",
													Description: "The first output.",
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.DynamicPseudoType).
														SetOriginalValueExpression([]byte("null_resource.test.id")).
														Export(),
												},
											},
											"second": {
												Value: &openapi3.Schema{
													Title:       "second",
													Description: "The second output.",
													WriteOnly:   true,
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.DynamicPseudoType).
														SetOriginalValueExpression([]byte(`"some value"`)).
														Export(),
												},
											},
										},
										Extensions: openapi.NewExt(nil).
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
													Title: "any",
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.DynamicPseudoType).
														Export(),
												},
											},
											"any_map": {
												Value: &openapi3.Schema{
													Title:      "any_map",
													Type:       "object",
													Properties: map[string]*openapi3.SchemaRef{},
													AdditionalProperties: openapi3.AdditionalProperties{
														Schema: &openapi3.SchemaRef{
															Value: &openapi3.Schema{
																Extensions: openapi.NewExt(nil).
																	SetOriginalType(cty.DynamicPseudoType).
																	Export(),
															},
														},
													},
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.Map(cty.DynamicPseudoType)).
														Export(),
												},
											},
											"string_map": {
												Value: &openapi3.Schema{
													Title: "string_map",
													Type:  "object",
													Default: map[string]any{
														"a": "a",
														"b": "1",
														"c": "true",
													},
													Properties: map[string]*openapi3.SchemaRef{},
													AdditionalProperties: openapi3.AdditionalProperties{
														Schema: &openapi3.SchemaRef{
															Value: &openapi3.Schema{
																Type: "string",
																Extensions: openapi.NewExt(nil).
																	SetOriginalType(cty.String).
																	Export(),
															},
														},
													},
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.Map(cty.String)).
														Export(),
												},
											},
											"string_slice": {
												Value: &openapi3.Schema{
													Title:   "string_slice",
													Type:    "array",
													Default: []any{"x", "y", "z"},
													Items: &openapi3.SchemaRef{
														Value: &openapi3.Schema{
															Type: "string",
															Extensions: openapi.NewExt(nil).
																SetOriginalType(cty.String).
																Export(),
														},
													},
													Extensions: openapi.NewExt(nil).
														SetOriginalType(cty.List(cty.String)).
														Export(),
												},
											},
											"object": {
												Value: &openapi3.Schema{
													Title: "object",
													Type:  "object",
													Default: map[string]any{
														"a": "a",
														"b": float64(1),
														"c": true,
													},
													Properties: map[string]*openapi3.SchemaRef{
														"a": {
															Value: &openapi3.Schema{
																Title: "a",
																Type:  "string",
																Extensions: openapi.NewExt(nil).
																	SetOriginalType(cty.String).
																	Export(),
															},
														},
														"b": {
															Value: &openapi3.Schema{
																Title: "b",
																Type:  "number",
																Extensions: openapi.NewExt(nil).
																	SetOriginalType(cty.Number).
																	Export(),
															},
														},
														"c": {
															Value: &openapi3.Schema{
																Title: "c",
																Type:  "boolean",
																Extensions: openapi.NewExt(nil).
																	SetOriginalType(cty.Bool).
																	Export(),
															},
														},
													},
													Required: []string{"a", "b", "c"},
													Extensions: openapi.NewExt(nil).
														SetOriginalType(
															cty.Object(
																map[string]cty.Type{
																	"a": cty.String,
																	"b": cty.Number,
																	"c": cty.Bool,
																},
															)).
														Export(),
												},
											},
											"object_nested": {
												Value: &openapi3.Schema{
													Title: "object_nested",
													Type:  "object",
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
																Title: "a",
																Type:  "string",
																Extensions: openapi.NewExt(nil).
																	SetOriginalType(cty.String).
																	Export(),
															},
														},
														"b": {
															Value: &openapi3.Schema{
																Title: "b",
																Type:  "array",
																Items: &openapi3.SchemaRef{
																	Value: &openapi3.Schema{
																		Type: "object",
																		Properties: map[string]*openapi3.SchemaRef{
																			"c": {
																				Value: &openapi3.Schema{
																					Title: "c",
																					Type:  "boolean",
																					Extensions: openapi.NewExt(nil).
																						SetOriginalType(cty.Bool).
																						Export(),
																				},
																			},
																		},
																		Required: []string{"c"},
																		Extensions: openapi.NewExt(nil).SetOriginalType(
																			cty.Object(map[string]cty.Type{
																				"c": cty.Bool,
																			})).
																			Export(),
																	},
																},
																Extensions: openapi.NewExt(nil).SetOriginalType(
																	cty.List(
																		cty.Object(map[string]cty.Type{
																			"c": cty.Bool,
																		}))).
																	Export(),
															},
														},
													},
													Extensions: openapi.NewExt(nil).
														SetOriginalType(
															cty.Object(map[string]cty.Type{
																"a": cty.String,
																"b": cty.List(
																	cty.Object(map[string]cty.Type{
																		"c": cty.Bool,
																	})),
															})).
														Export(),
												},
											},
											"list_object": {
												Value: &openapi3.Schema{
													Title: "list_object",
													Type:  "array",
													Items: &openapi3.SchemaRef{
														Value: &openapi3.Schema{
															Type: "object",
															Properties: map[string]*openapi3.SchemaRef{
																"a": {
																	Value: &openapi3.Schema{
																		Title: "a",
																		Type:  "string",
																		Extensions: openapi.NewExt(nil).
																			SetOriginalType(cty.String).
																			Export(),
																	},
																},
																"b": {
																	Value: &openapi3.Schema{
																		Title: "b",
																		Type:  "number",
																		Extensions: openapi.NewExt(nil).
																			SetOriginalType(cty.Number).
																			Export(),
																	},
																},
																"c": {
																	Value: &openapi3.Schema{
																		Title: "c",
																		Type:  "boolean",
																		Extensions: openapi.NewExt(nil).
																			SetOriginalType(cty.Bool).
																			Export(),
																	},
																},
															},
															Required: []string{"a", "b", "c"},
															Extensions: openapi.NewExt(nil).SetOriginalType(
																cty.Object(
																	map[string]cty.Type{
																		"a": cty.String,
																		"b": cty.Number,
																		"c": cty.Bool,
																	})).
																Export(),
														},
													},
													Extensions: openapi.NewExt(nil).
														SetOriginalType(
															cty.List(
																cty.Object(
																	map[string]cty.Type{
																		"a": cty.String,
																		"b": cty.Number,
																		"c": cty.Bool,
																	}))).
														Export(),
												},
											},
											"tuple": {
												Value: &openapi3.Schema{
													Title:     "tuple",
													Type:      "array",
													MaxLength: &length,
													MinLength: 3,
													Items: &openapi3.SchemaRef{
														Value: &openapi3.Schema{
															OneOf: openapi3.SchemaRefs{
																{
																	Value: &openapi3.Schema{
																		Type: "string",
																		Extensions: openapi.NewExt(nil).
																			SetOriginalType(cty.String).
																			Export(),
																	},
																},
																{
																	Value: &openapi3.Schema{
																		Type: "boolean",
																		Extensions: openapi.NewExt(nil).
																			SetOriginalType(cty.Bool).
																			Export(),
																	},
																},
																{
																	Value: &openapi3.Schema{
																		Type: "number",
																		Extensions: openapi.NewExt(nil).
																			SetOriginalType(cty.Number).
																			Export(),
																	},
																},
															},
														},
													},
													Extensions: openapi.NewExt(nil).
														SetOriginalType(
															cty.Tuple(
																[]cty.Type{
																	cty.String, cty.Bool, cty.Number,
																})).
														Export(),
												},
											},
										},
										Required: []string{
											"any",
											"list_object",
											"tuple",
										},
										Extensions: openapi.NewExt(nil).
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

			actualOutput, actualError := loader.Load(tc.input, "dev-template", "dev")
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
