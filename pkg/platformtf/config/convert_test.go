package config

import (
	"reflect"
	"testing"

	"github.com/seal-io/seal/pkg/dao/model"
)

func TestToModuleBlock(t *testing.T) {
	type input struct {
		Module     *model.Module
		Attributes map[string]interface{}
	}

	testCases := []struct {
		Name     string
		Input    input
		Expected Block
	}{
		{
			Name: "Module with no attributes",
			Input: input{
				Module: &model.Module{
					ID: "test",
				},
				Attributes: map[string]interface{}{},
			},
			Expected: Block{
				Type:       BlockTypeModule,
				Labels:     []string{"test"},
				Attributes: map[string]interface{}{},
			},
		},
		{
			Name: "Module with attributes",
			Input: input{
				Module: &model.Module{
					ID: "test",
				},
				Attributes: map[string]interface{}{
					"test": "test",
				},
			},
			Expected: Block{
				Type:   BlockTypeModule,
				Labels: []string{"test"},
				Attributes: map[string]interface{}{
					"test": "test",
				},
			},
		},
		{
			Name: "Module with null attributes",
			Input: input{
				Module: &model.Module{
					ID: "test",
				},
				Attributes: map[string]interface{}{
					"test": nil,
				},
			},
			Expected: Block{
				Type:       BlockTypeModule,
				Labels:     []string{"test"},
				Attributes: map[string]interface{}{},
			},
		},
		{
			Name: "Module with nested attributes and null keys",
			Input: input{
				Module: &model.Module{
					ID: "test",
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
			Expected: Block{
				Type:   BlockTypeModule,
				Labels: []string{"test"},
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
			block := ToModuleBlock(tc.Input.Module, tc.Input.Attributes)
			if block.Labels[0] != tc.Expected.Labels[0] {
				t.Errorf("expected block label %s, got %s", tc.Expected.Labels[0], block.Labels[0])
			}
			if reflect.DeepEqual(block.Attributes, tc.Expected.Attributes) {
				t.Errorf("expected block attributes %v, got %v", tc.Expected.Attributes, block.Attributes)
			}
		})
	}
}
