package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/seal-io/walrus/pkg/terraform/block"
)

func TestCreateConfigToBytes(t *testing.T) {
	testCases := []struct {
		name     string
		option   CreateOptions
		expected []byte
	}{
		{
			name:     "test create config to bytes with empty option",
			option:   CreateOptions{},
			expected: nil,
		},
		{
			name: "test create config to bytes with attributes",
			option: CreateOptions{
				Attributes: map[string]any{
					"var1":    "ami-0c55b159cbfafe1f0",
					"secret1": "password",
				},
			},
			expected: []byte(`secret1 = "password"
var1    = "ami-0c55b159cbfafe1f0"
`),
		},
		{
			name: "test create config to bytes with outputs",
			option: CreateOptions{
				OutputOptions: []Output{
					{
						ResourceName: "test-resource",
						Name:         "test-output",
					},
					{
						ResourceName: "test-resource",
						Name:         "test-output-sensitive",
						Sensitive:    true,
					},
				},
			},
			expected: []byte(`output "test-resource_test-output" {
  sensitive = false
  value     = module.test-resource.test-output
}

output "test-resource_test-output-sensitive" {
  sensitive = true
  value     = module.test-resource.test-output-sensitive
}

`),
		},
	}

	for _, tt := range testCases {
		got, err := CreateConfigToBytes(tt.option)
		if err != nil {
			assert.Errorf(t, err, "unexpected error: %v", err)
		}

		if !assert.Equal(t, string(tt.expected), string(got)) {
			assert.Errorf(t, err, "name: %s", tt.name)
		}
	}
}

func TestToModuleBlock(t *testing.T) {
	testCases := []struct {
		Name         string
		ModuleConfig *ModuleConfig
		Expected     block.Block
	}{
		{
			Name: "Template with no attributes",
			ModuleConfig: &ModuleConfig{
				Name:       "test1",
				Attributes: map[string]any{},
			},
			Expected: block.Block{
				Type:       block.TypeModule,
				Labels:     []string{"test1"},
				Attributes: map[string]any{},
			},
		},
		{
			Name: "Template with attributes",
			ModuleConfig: &ModuleConfig{
				Name: "test2",
				Attributes: map[string]any{
					"test": "test",
				},
			},
			Expected: block.Block{
				Type:   block.TypeModule,
				Labels: []string{"test2"},
				Attributes: map[string]any{
					"test": "test",
				},
			},
		},
		{
			Name: "Template with null attributes",
			ModuleConfig: &ModuleConfig{
				Name: "test3",
				Attributes: map[string]any{
					"test": nil,
				},
			},
			Expected: block.Block{
				Type:       block.TypeModule,
				Labels:     []string{"test3"},
				Attributes: map[string]any{},
			},
		},
		{
			Name: "Template with nested attributes and null keys",
			ModuleConfig: &ModuleConfig{
				Name: "test4",
				Attributes: map[string]any{
					"test": map[string]any{
						"test": "test",
						"foo":  nil,
					},
					"foo": nil,
					"blob": []any{
						map[string]any{
							"test": "test",
							"foo":  nil,
						},
						nil,
					},
					"clob": []map[string]any{
						{
							"test": "test",
							"foo":  nil,
						},
						nil,
						{
							"blob": []any{
								map[string]any{
									"test": "test",
									"foo":  nil,
								},
								nil,
							},
						},
					},
				},
			},
			Expected: block.Block{
				Type:   block.TypeModule,
				Labels: []string{"test4"},
				Attributes: map[string]any{
					"test": map[string]any{
						"test": "test",
					},
					"blob": []any{
						map[string]any{
							"test": "test",
						},
					},
					"clob": []map[string]any{
						{
							"test": "test",
						},
						{
							"blob": []any{
								map[string]any{
									"test": "test",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			moduleBlock, err := ToModuleBlock(tc.ModuleConfig)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if moduleBlock.Labels[0] != tc.Expected.Labels[0] {
				t.Errorf("expected block label %s, got %s", tc.Expected.Labels[0], moduleBlock.Labels[0])
			}

			if reflect.DeepEqual(moduleBlock.Attributes, tc.Expected.Attributes) {
				t.Errorf("expected block attributes %v, got %v",
					tc.Expected.Attributes, moduleBlock.Attributes)
			}
		})
	}
}

func TestTypeExprTokens(t *testing.T) {
	cases := []struct {
		name     string
		ctyType  cty.Type
		expected string
		wantErr  bool
	}{
		{
			name:     "test get cty type with string",
			ctyType:  cty.String,
			expected: "string",
			wantErr:  false,
		},
		{
			name:     "test get cty type with number",
			ctyType:  cty.Number,
			expected: "number",
			wantErr:  false,
		},
		{
			name:     "test get cty type with bool",
			ctyType:  cty.Bool,
			expected: "bool",
			wantErr:  false,
		},
		{
			name:     "test get cty type with list",
			ctyType:  cty.List(cty.String),
			expected: "list(string)",
			wantErr:  false,
		},
		{
			name:     "test get cty type with map",
			ctyType:  cty.Map(cty.String),
			expected: "map(string)",
			wantErr:  false,
		},
		{
			name:     "test get cty type with object",
			ctyType:  cty.Object(map[string]cty.Type{"test": cty.String}),
			expected: "{\n  test = string\n}",
			wantErr:  false,
		},
		{
			name:     "test get cty type with tuple",
			ctyType:  cty.Tuple([]cty.Type{cty.String}),
			expected: "[string]",
			wantErr:  false,
		},
		{
			name:     "test get cty type with set",
			ctyType:  cty.Set(cty.String),
			expected: "set(string)",
			wantErr:  false,
		},
		{
			name:     "test get cty type with dynamic",
			ctyType:  cty.DynamicPseudoType,
			expected: "any",
			wantErr:  false,
		},
		{
			name:    "test get cty type with nil",
			ctyType: cty.NilType,
			wantErr: true,
		},
		{
			name:    "test error with nil set",
			ctyType: cty.Set(cty.NilType),
			wantErr: true,
		},
		{
			name:    "test error with nil object",
			ctyType: cty.Object(map[string]cty.Type{"test": cty.NilType}),
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := typeExprTokens(tt.ctyType)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %s", err)
			}

			if !assert.Equal(t, tt.expected, string(got.Bytes())) {
				t.Errorf("expected %s, got %s", tt.expected, string(got.Bytes()))
			}
		})
	}
}
