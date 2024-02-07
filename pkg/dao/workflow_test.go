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
				"description": "${workflow.var.var-1}",
			},
			params: map[string]string{
				"var-1": "aaa",
			},
			config: types.WorkflowVariables{
				{
					Name:      "var-1",
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
				"description": "${workflow.var.var-1}",
			},
			config: types.WorkflowVariables{
				{
					Name:      "var-1",
					Value:     "alice",
					Overwrite: false,
				},
			},
			expected: map[string]any{
				"description": "alice",
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

func TestOverwriteWorkflowVariables(t *testing.T) {
	cases := []struct {
		name      string
		vars      map[string]string
		config    types.WorkflowVariables
		expected  map[string]string
		wantError bool
	}{
		{
			name: "overwrite",
			vars: map[string]string{
				"replace": "newValue",
			},
			config: types.WorkflowVariables{
				{
					Name:      "replace",
					Value:     "oldValue",
					Overwrite: true,
				},
			},
			expected: map[string]string{
				"replace": "newValue",
			},
			wantError: false,
		},
		{
			name: "overwrite multi configs",
			vars: map[string]string{
				"replace": "newValue",
			},
			config: types.WorkflowVariables{
				{
					Name:      "replace",
					Value:     "oldValue",
					Overwrite: true,
				},
				{
					Name:      "var1",
					Value:     "var1",
					Overwrite: true,
				},
			},
			expected: map[string]string{
				"replace": "newValue",
				"var1":    "var1",
			},
			wantError: false,
		},
		{
			name: "overwrite multi configs with one not overwrite",
			vars: map[string]string{
				"replace":   "newValue",
				"notConfig": "notConfig",
			},
			config: types.WorkflowVariables{
				{
					Name:      "replace",
					Value:     "oldValue",
					Overwrite: true,
				},
				{
					Name:      "var1",
					Value:     "var1",
					Overwrite: false,
				},
			},
			expected: map[string]string{
				"replace": "newValue",
				"var1":    "var1",
			},
			wantError: true,
		},
		{
			name: "overwrite multi configs with correct overwrite",
			vars: map[string]string{
				"replace": "newValue",
				"var1":    "newVar",
			},
			config: types.WorkflowVariables{
				{
					Name:      "replace",
					Value:     "oldValue",
					Overwrite: true,
				},
				{
					Name:      "var1",
					Value:     "var1",
					Overwrite: true,
				},
				{
					Name:      "var2",
					Value:     "var2",
					Overwrite: false,
				},
			},
			expected: map[string]string{
				"replace": "newValue",
				"var1":    "newVar",
				"var2":    "var2",
			},
			wantError: false,
		},
		{
			name: "not overwrite",
			vars: map[string]string{
				"replace": "newValue",
			},
			config: types.WorkflowVariables{
				{
					Name:      "replace",
					Value:     "oldValue",
					Overwrite: false,
				},
			},
			wantError: true,
		},
		{
			name: "no config",
			vars: map[string]string{
				"replace": "newValue",
			},
			config:    types.WorkflowVariables{},
			wantError: true,
		},
	}

	for _, c := range cases {
		actual, err := OverwriteWorkflowVariables(c.vars, c.config)
		if err != nil && c.wantError == false {
			t.Errorf("overwrite workflow variables error: expected nil, got %v", err)
		}

		if c.wantError == true {
			continue
		}

		if reflect.DeepEqual(actual, c.expected) == false {
			t.Errorf("overwrite workflow variables error: expected %v, got %v", c.expected, actual)
		}
	}
}
