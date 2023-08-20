package templates

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/pointer"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/utils/json"
)

func TestGetTemplateNameByPath(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "versioned template path",
			input:    "templates/foo/0.0.1",
			expected: "foo",
		},
		{
			name:     "non-versioned template path",
			input:    "templates/foo",
			expected: "foo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := GetTemplateNameByPath(tc.input)
			assert.Equal(t, tc.expected, actualOutput)
		})
	}
}

func TestGetVersionedSource(t *testing.T) {
	testCases := []struct {
		name    string
		source  string
		version string
		result  string
	}{
		{
			name:    "github with subdirectory",
			source:  "github.com/foo/bar//module1",
			version: "0.0.1",
			result:  "github.com/foo/bar//module1/0.0.1",
		},
		{
			name:    "github root",
			source:  "github.com/foo/bar",
			version: "0.0.1",
			result:  "github.com/foo/bar//0.0.1",
		},
		{
			name:    "github with ref",
			source:  "github.com/foo/bar//module1?ref=dev",
			version: "0.0.1",
			result:  "github.com/foo/bar//module1/0.0.1?ref=dev",
		},
		{
			name:    "generic git",
			source:  "git::https://github.com/foo/bar.git",
			version: "0.0.1",
			result:  "git::https://github.com/foo/bar.git//0.0.1",
		},
		{
			name:    "generic git with subdirectory",
			source:  "git::https://github.com/foo/bar.git//module1",
			version: "0.0.1",
			result:  "git::https://github.com/foo/bar.git//module1/0.0.1",
		},
		{
			name:    "generic git with ref",
			source:  "git::https://github.com/foo/bar.git//module1?ref=dev",
			version: "0.0.1",
			result:  "git::https://github.com/foo/bar.git//module1/0.0.1?ref=dev",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := getVersionedSource(tc.source, tc.version)
			assert.Equal(t, tc.result, actualOutput)
		})
	}
}

func TestLoadTerraformSchema(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput *types.TemplateSchema
		expectedError  bool
	}{
		{
			name:           "Invalid",
			input:          "testdata/invalid",
			expectedOutput: &types.TemplateSchema{},
			expectedError:  true,
		},
		{
			name:           "Invalid variable type",
			input:          "testdata/invalid_variable_type",
			expectedOutput: &types.TemplateSchema{},
			expectedError:  true,
		},
		{
			name:  "With README.md",
			input: "testdata/with_readme",
			expectedOutput: &types.TemplateSchema{
				Readme: "# test readme",
			},
			expectedError: false,
		},
		{
			name:  "With output",
			input: "testdata/with_output",
			expectedOutput: &types.TemplateSchema{
				Outputs: property.Schemas{
					property.AnySchema("first", nil).
						WithDescription("The first output.").
						WithValue([]byte("1")),
					property.AnySchema("second", nil).
						WithDescription("The second output.").
						WithSensitive().
						WithValue([]byte("2")),
				},
			},
			expectedError: false,
		},
		{
			name:  "With variable",
			input: "testdata/with_variable",
			expectedOutput: &types.TemplateSchema{
				Variables: property.Schemas{
					property.StringSchema("foo", pointer.String("bar")),
				},
			},
			expectedError: false,
		},
		{
			name:  "With description",
			input: "testdata/with_description",
			expectedOutput: &types.TemplateSchema{
				Variables: property.Schemas{
					property.StringSchema("foo", pointer.String("bar")).
						WithDescription("description of foo."),
				},
			},
			expectedError: false,
		},
		{
			name:  "Full schema",
			input: "testdata/full_schema",
			expectedOutput: &types.TemplateSchema{
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
				Outputs: property.Schemas{
					property.AnySchema("first", nil).
						WithDescription("The first output.").
						WithValue([]byte("null_resource.test.id")),
					property.AnySchema("second", nil).
						WithDescription("The second output.").
						WithSensitive().
						WithValue([]byte(`"some value"`)),
				},
				Variables: property.Schemas{
					property.StringSchema("foo", pointer.String("foo")).
						WithGroup("Test Group").
						WithLabel("Foo Label").
						WithOptions(json.MustMarshal("F1"), json.MustMarshal("F2"), json.MustMarshal("F3")),
					property.StringSchema("bar", pointer.String("bar")).
						WithDescription("description of bar.").
						WithGroup("Test Group").
						WithLabel("Bar Label").
						WithOptions(json.MustMarshal("B1"), json.MustMarshal("B2"), json.MustMarshal("B3")).
						WithShowIf("foo=F1"),
					property.StringSchema("thee", pointer.String("thee")),
					property.IntSchema("number_options_var", pointer.Int(1)).
						WithOptions(json.MustMarshal(1), json.MustMarshal(2), json.MustMarshal(3)),
					property.StringSchema("subgroup1_1", pointer.String("subgroup1_1")).
						WithGroup("Test Subgroup/Subgroup 1").
						WithLabel("Subgroup1_1 Label"),
					property.StringSchema("subgroup1_2", pointer.String("subgroup1_2")).
						WithGroup("Test Subgroup/Subgroup 1").
						WithLabel("Subgroup1_2 Label"),
					property.StringSchema("subgroup2_1", pointer.String("subgroup2_1")).
						WithGroup("Test Subgroup/Subgroup 2").
						WithLabel("Subgroup2_1 Label"),
					property.StringSchema("subgroup2_1_hidden", pointer.String("")).
						WithGroup("Test Subgroup/Subgroup 2").
						WithHidden(),
				},
			},
			expectedError: false,
		},
		{
			name:  "Complex variable",
			input: "testdata/complex_variable",
			expectedOutput: &types.TemplateSchema{
				Variables: property.Schemas{
					property.AnySchema("any", nil).
						WithRequired(),
					property.MapSchema("any_map",
						map[string]any(nil)),
					property.MapSchema("string_map",
						map[string]string{
							"a": "a",
							"b": "1",
							"c": "true",
						}),
					property.SliceSchema("string_slice",
						[]*string{
							pointer.String("x"),
							pointer.String("y"),
							pointer.String("z"),
						}),
					property.ObjectSchema("object",
						&struct {
							A string `cty:"a" json:"a"`
							B int    `cty:"b" json:"b"`
							C bool   `cty:"c" json:"c"`
						}{
							A: "a",
							B: 1,
							C: true,
						}),
					property.ObjectSchema("object_nested",
						&struct {
							A string `cty:"a" json:"a"`
							B []struct {
								C bool `cty:"c" json:"c"`
							} `cty:"b" json:"b"`
						}{
							A: "a",
							B: []struct {
								C bool `cty:"c" json:"c"`
							}{
								{
									C: true,
								},
							},
						}),
					property.SliceSchema("list_object",
						[]struct {
							A string `cty:"a" json:"a"`
							B int    `cty:"b" json:"b"`
							C bool   `cty:"c" json:"c"`
						}(nil)).
						WithRequired(),
					// TODO(thxCode): provide a tuple schema builder?
					property.Schema{
						Name: "tuple",
						Type: cty.Tuple([]cty.Type{
							cty.String,
							cty.Bool,
							cty.Number,
						}),
						Required: true,
					},
				},
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput, actualError := loadTerraformTemplateSchema(tc.input)

			if tc.expectedError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}

			if actualOutput == nil {
				return
			}

			assert.Equal(t, tc.expectedOutput, actualOutput)
		})
	}
}

func TestLoadTerraformTemplateVersions(t *testing.T) {
	testCases := []struct {
		name           string
		input          *model.Template
		expectedOutput []*model.TemplateVersion
		expectedError  bool
	}{
		{
			name: "versioned-templates",
			input: &model.Template{
				ID:     "mock-id",
				Source: "testdata/versioned_template",
			},
			expectedOutput: []*model.TemplateVersion{
				{
					TemplateID: "mock-id",
					Version:    "0.0.1",
					Schema: &types.TemplateSchema{
						Readme: "# Version 0.0.1",
					},
				},
				{
					TemplateID: "mock-id",
					Version:    "0.0.2",
					Schema: &types.TemplateSchema{
						Readme: "# Version 0.0.2",
					},
				},
				{
					TemplateID: "mock-id",
					Version:    "100.0.0",
					Schema: &types.TemplateSchema{
						Readme: "# Version 100.0.0",
					},
				},
			},
		},
	}

	// Absolute path is required for getter.Get(dest, local_source).
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get pwd: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.Source = filepath.Join(pwd, tc.input.Source)
			actualOutput, actualError := loadTerraformTemplateVersions(tc.input)

			if tc.expectedError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}

			assert.Equal(t, len(tc.expectedOutput), len(actualOutput))

			for i, v := range tc.expectedOutput {
				actualV := actualOutput[i]

				assert.Equal(t, actualV.TemplateID, v.TemplateID)
				assert.Equal(t, actualV.Version, v.Version)

				assert.Equal(t, actualV.Schema, v.Schema)
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
