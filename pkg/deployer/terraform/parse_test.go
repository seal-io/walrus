package terraform

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToDependOutputMap(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected map[string]string
	}{
		{
			name:     "empty",
			input:    nil,
			expected: map[string]string{},
		},
		{
			name:  "output with service reference",
			input: []string{"service_s1_output1"},
			expected: map[string]string{
				"s1_output1": "service",
			},
		},
		{
			name:  "output with resource reference",
			input: []string{"resource_s1_output1"},
			expected: map[string]string{
				"s1_output1": "resource",
			},
		},
		{
			name: "outputs with both service and resource reference",
			input: []string{
				"service_s1_output1",
				"resource_s2_output2",
			},
			expected: map[string]string{
				"s1_output1": "service",
				"s2_output2": "resource",
			},
		},
	}

	for _, tc := range testCases {
		actual := toDependOutputMap(tc.input)
		assert.Equal(t, tc.expected, actual, fmt.Sprintf("unexpected result in test case: %s", tc.name))
	}
}

func TestParseAttributeReplace(t *testing.T) {
	testCases := []struct {
		name                    string
		attributes              map[string]any
		replaced                bool
		expectedVariableNames   []string
		expectedResourceOutputs []string
		expectedAttributes      map[string]any
		expectedError           bool
	}{
		{
			name: "no reference",
			attributes: map[string]any{
				"foo": "bar",
			},
			replaced:                false,
			expectedVariableNames:   []string{},
			expectedResourceOutputs: []string{},
			expectedAttributes: map[string]any{
				"foo": "bar",
			},
			expectedError: false,
		},
		{
			name: "parse var reference",
			attributes: map[string]any{
				"foo": "${var.foo}",
			},
			replaced:                false,
			expectedVariableNames:   []string{"foo"},
			expectedResourceOutputs: []string{},
			expectedAttributes: map[string]any{
				"foo": "${var.foo}",
			},
			expectedError: false,
		},
		{
			name: "parse resource reference",
			attributes: map[string]any{
				"foo": "${res.foo.bar}",
			},
			replaced:                true,
			expectedVariableNames:   []string{},
			expectedResourceOutputs: []string{"res_foo_bar"},
			expectedAttributes: map[string]any{
				"foo": "${var._walrus_res_foo_bar}",
			},
			expectedError: false,
		},
		{
			name: "parse combined",
			attributes: map[string]any{
				"foo": "${var.foo}",
				"bar": "${svc.foo1.bar}-${res.foo2.bar}",
			},
			replaced:                true,
			expectedVariableNames:   []string{"foo"},
			expectedResourceOutputs: []string{"res_foo2_bar", "svc_foo1_bar"},
			expectedAttributes: map[string]any{
				"foo": "${var._walrus_var_foo}",
				"bar": "${var._walrus_res_foo1_bar}-${var._walrus_res_foo2_bar}",
			},
			expectedError: false,
		},
		{
			name: "parse combined with interpolation",
			attributes: map[string]any{
				"foo":    "${var.foo}",
				"bar":    "${svc.foo1.bar}-${res.foo2.bar}",
				"baz":    "${var.foo}-${svc.foo1.bar}-${res.foo2.bar}",
				"qux":    "${MYSQL_DATABASE}",
				"double": "$${ENV_PORT}",
			},
			replaced:                true,
			expectedVariableNames:   []string{"foo"},
			expectedResourceOutputs: []string{"res_foo2_bar", "svc_foo1_bar"},
			expectedAttributes: map[string]any{
				"foo":    "${var._walrus_var_foo}",
				"bar":    "${var._walrus_res_foo1_bar}-${var._walrus_res_foo2_bar}",
				"baz":    "${var._walrus_var_foo}-${var._walrus_res_foo1_bar}-${var._walrus_res_foo2_bar}",
				"qux":    "$${MYSQL_DATABASE}",
				"double": "$$${ENV_PORT}", // Terraform will replace $$ with $.
			},
			expectedError: false,
		},
		{
			name: "parse array with interpolation",
			attributes: map[string]any{
				"foo": []string{
					"${var.foo}",
					"${svc.foo1.bar}-${res.foo2.bar}",
				},
				"ENV": []string{
					"${ENV_PORT}",
					"${ENV_HOST}",
				},
			},
			replaced:                true,
			expectedVariableNames:   []string{"foo"},
			expectedResourceOutputs: []string{"res_foo2_bar", "svc_foo1_bar"},
			expectedAttributes: map[string]any{
				"foo": []any{
					"${var._walrus_var_foo}",
					"${var._walrus_res_foo1_bar}-${var._walrus_res_foo2_bar}",
				},
				"ENV": []any{
					"$${ENV_PORT}",
					"$${ENV_HOST}",
				},
			},
		},
		{
			name: "parse error attribute",
			attributes: map[string]any{
				"foo":     "${var.foo}",
				"invalid": make(chan int),
				"txt": []any{
					"foo",
					map[string]any{
						"bar": "${var.bar}",
					},
				},
			},
			replaced:      true,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		attrs, actualVariableNames, actualResourceOutputs, err := parseAttributeReplace(
			tc.attributes,
			tc.replaced,
		)

		if tc.expectedError && err != nil {
			continue
		}

		if err != nil {
			t.Errorf("expected error but got nil in test case: %s", tc.name)
		}

		assert.Equal(
			t,
			tc.expectedVariableNames,
			actualVariableNames,
			fmt.Sprintf("unexpected result in test case: %s", tc.name),
		)
		assert.Equal(
			t,
			tc.expectedResourceOutputs,
			actualResourceOutputs,
			fmt.Sprintf("unexpected result in test case: %s", tc.name),
		)
		assert.Equal(
			t,
			tc.expectedAttributes,
			attrs,
			fmt.Sprintf("unexpected result in test case: %s", tc.name),
		)
	}
}
