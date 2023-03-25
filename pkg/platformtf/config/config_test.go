package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/platformtf/block"
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
				Attributes: map[string]interface{}{
					"var1":    "ami-0c55b159cbfafe1f0",
					"secret1": "password",
				},
			},
			expected: []byte(`secret1 = "password"
var1    = "ami-0c55b159cbfafe1f0"
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
			Name: "Module with no attributes",
			ModuleConfig: &ModuleConfig{
				Name: "test1",
				ModuleVersion: &model.ModuleVersion{
					ModuleID: "test",
					Version:  "0.0.0",
				},
				Attributes: map[string]interface{}{},
			},
			Expected: block.Block{
				Type:       block.TypeModule,
				Labels:     []string{"test1"},
				Attributes: map[string]interface{}{},
			},
		},
		{
			Name: "Module with attributes",
			ModuleConfig: &ModuleConfig{
				Name: "test2",
				ModuleVersion: &model.ModuleVersion{
					ModuleID: "test",
					Version:  "0.0.0",
				},
				Attributes: map[string]interface{}{
					"test": "test",
				},
			},
			Expected: block.Block{
				Type:   block.TypeModule,
				Labels: []string{"test2"},
				Attributes: map[string]interface{}{
					"test": "test",
				},
			},
		},
		{
			Name: "Module with null attributes",
			ModuleConfig: &ModuleConfig{
				Name: "test3",
				ModuleVersion: &model.ModuleVersion{
					ModuleID: "test",
					Version:  "0.0.0",
				},
				Attributes: map[string]interface{}{
					"test": nil,
				},
			},
			Expected: block.Block{
				Type:       block.TypeModule,
				Labels:     []string{"test3"},
				Attributes: map[string]interface{}{},
			},
		},
		{
			Name: "Module with nested attributes and null keys",
			ModuleConfig: &ModuleConfig{
				Name: "test4",
				ModuleVersion: &model.ModuleVersion{
					ModuleID: "test",
					Version:  "0.0.0",
				},
				Attributes: map[string]interface{}{
					"test": map[string]interface{}{
						"test": "test",
						"foo":  nil,
					},
					"foo": nil,
					"blob": []interface{}{
						map[string]interface{}{
							"test": "test",
							"foo":  nil,
						},
						nil,
					},
					"clob": []map[string]interface{}{
						{
							"test": "test",
							"foo":  nil,
						},
						nil,
						{
							"blob": []interface{}{
								map[string]interface{}{
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
				Attributes: map[string]interface{}{
					"test": map[string]interface{}{
						"test": "test",
					},
					"blob": []interface{}{
						map[string]interface{}{
							"test": "test",
						},
					},
					"clob": []map[string]interface{}{
						{
							"test": "test",
						},
						{
							"blob": []interface{}{
								map[string]interface{}{
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
				t.Errorf("expected block attributes %v, got %v", tc.Expected.Attributes, moduleBlock.Attributes)
			}
		})
	}
}
