package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

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
						ServiceName: "test-service",
						Name:        "test-output",
					},
					{
						ServiceName: "test-service",
						Name:        "test-output-sensitive",
						Sensitive:   true,
					},
				},
			},
			expected: []byte(`output "test-service_test-output" {
  sensitive = false
  value     = module.test-service.test-output
}

output "test-service_test-output-sensitive" {
  sensitive = true
  value     = module.test-service.test-output-sensitive
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
