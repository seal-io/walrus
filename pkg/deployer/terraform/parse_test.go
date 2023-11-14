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
	type (
		input = struct {
			attributes      map[string]any
			variableNames   []string
			resourceOutputs []string
			replaced        bool
		}
		output = struct {
			variableNames []string
			outputNames   []string
		}
	)

	testCases := []struct {
		name     string
		input    input
		expected output
	}{
		{
			name: "no reference",
			input: input{
				attributes: map[string]any{
					"foo": "bar",
				},
				variableNames:   []string{},
				resourceOutputs: []string{},
				replaced:        false,
			},
			expected: output{
				variableNames: []string{},
				outputNames:   []string{},
			},
		},
		{
			name: "parse var reference",
			input: input{
				attributes: map[string]any{
					"foo": "${var.foo}",
				},
				variableNames:   []string{},
				resourceOutputs: []string{},
				replaced:        false,
			},
			expected: output{
				variableNames: []string{
					"foo",
				},
				outputNames: []string{},
			},
		},
		{
			name: "parse resource reference",
			input: input{
				attributes: map[string]any{
					"foo": "${resource.foo.bar}",
				},
				variableNames:   []string{},
				resourceOutputs: []string{},
				replaced:        false,
			},
			expected: output{
				variableNames: []string{},
				outputNames: []string{
					"resource_foo_bar",
				},
			},
		},
		{
			name: "parse service reference",
			input: input{
				attributes: map[string]any{
					"foo": "${service.foo.bar}",
				},
				variableNames:   []string{},
				resourceOutputs: []string{},
				replaced:        false,
			},
			expected: output{
				variableNames: []string{},
				outputNames: []string{
					"service_foo_bar",
				},
			},
		},
		{
			name: "parse combined",
			input: input{
				attributes: map[string]any{
					"foo": "${var.foo}",
					"bar": "${service.foo1.bar}-${resource.foo2.bar}",
				},
				variableNames:   []string{},
				resourceOutputs: []string{},
				replaced:        false,
			},
			expected: output{
				variableNames: []string{
					"foo",
				},
				outputNames: []string{
					"service_foo1_bar",
					"resource_foo2_bar",
				},
			},
		},
	}

	for _, tc := range testCases {
		actualVariableNames, actualOutputNames := parseAttributeReplace(
			tc.input.attributes,
			tc.input.variableNames,
			tc.input.resourceOutputs,
			tc.input.replaced,
		)

		assert.Equal(
			t,
			tc.expected.variableNames,
			actualVariableNames,
			fmt.Sprintf("unexpected result in test case: %s", tc.name),
		)
		assert.Equal(
			t,
			tc.expected.outputNames,
			actualOutputNames,
			fmt.Sprintf("unexpected result in test case: %s", tc.name),
		)
	}
}
