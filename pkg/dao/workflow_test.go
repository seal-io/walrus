package dao

import (
	"reflect"
	"testing"

	"github.com/seal-io/walrus/pkg/dao/types"
)

func TestParseParams(t *testing.T) {
	cases := []struct {
		attributes map[string]any
		params     map[string]string
		config     types.WorkflowVariables
		expected   map[string]any
		wantError  bool
	}{
		{
			attributes: map[string]any{
				"description": "${workflow.var.env1}",
			},
			params: map[string]string{
				"env1": "aaa",
			},
			config: types.WorkflowVariables{
				{
					Name:      "env1",
					Value:     "alice",
					Overwrite: true,
				},
			},
			expected: map[string]any{
				"description": "aaa",
			},
			wantError: false,
		},
		{
			attributes: map[string]any{
				"deepAttr": map[string]any{
					"deepKey": "${workflow.var.replace}",
				},
			},
			params: map[string]string{
				"replace": "newValue",
			},
			config: types.WorkflowVariables{
				{
					Name:      "replace",
					Value:     "oldValue",
					Overwrite: true,
				},
			},
			expected: map[string]any{
				"deepAttr": map[string]any{
					"deepKey": "newValue",
				},
			},
			wantError: false,
		},
		{
			attributes: map[string]any{
				"deepAttr": map[string]any{
					"deepKey": "${workflow.var.replace}",
				},
			},
			params: map[string]string{},
			config: types.WorkflowVariables{
				{
					Name:      "replace",
					Value:     "oldValue",
					Overwrite: true,
				},
			},
			expected: map[string]any{
				"deepAttr": map[string]any{
					"deepKey": "oldValue",
				},
			},
			wantError: false,
		},
		{
			attributes: map[string]any{
				"deepAttr": map[string]any{
					"deepKey": "${workflow.var.replace}",
				},
			},
			params: map[string]string{
				"replace": "newValue",
			},
			config:    types.WorkflowVariables{},
			expected:  nil,
			wantError: true,
		},
	}

	for _, c := range cases {
		actual, err := parseWorkflowVariables(c.attributes, c.params, c.config)
		if c.wantError == true && err == nil {
			t.Errorf("parse params error: expected error, got nil")
		}

		if c.wantError == false && err != nil {
			t.Errorf("parse params error: expected nil, got %v", err)
		}

		if reflect.DeepEqual(actual, c.expected) == false {
			t.Errorf("parse params error: expected %v, got %v", c.expected, actual)
		}
	}
}
