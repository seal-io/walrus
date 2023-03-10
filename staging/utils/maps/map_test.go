package maps

import (
	"reflect"
	"testing"
)

func TestRemoveNulls(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "Map with no null keys",
			input: map[string]interface{}{
				"foo": "bar",
			},
			expected: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			name: "Map with null keys",
			input: map[string]interface{}{
				"foo": "bar",
				"baz": nil,
				"qux": []string{},
			},
			expected: map[string]interface{}{
				"foo": "bar",
				"qux": []string{},
			},
		},
		{
			name: "Map with null keys and nested maps",
			input: map[string]interface{}{
				"foo": "bar",
				"baz": map[string]interface{}{
					"qux": nil,
				},
			},
			expected: map[string]interface{}{
				"foo": "bar",
				"baz": map[string]interface{}{},
			},
		},
		{
			name: "Map with null keys and nested maps with non-null keys",
			input: map[string]interface{}{
				"foo": "bar",
				"baz": map[string]interface{}{
					"qux": "quux",
				},
			},
			expected: map[string]interface{}{
				"foo": "bar",
				"baz": map[string]interface{}{
					"qux": "quux",
				},
			},
		},
		{
			name: "Map with null keys and nested maps with null keys",
			input: map[string]interface{}{
				"foo": "bar",
				"baz": map[string]interface{}{
					"qux":  nil,
					"quux": "quuz",
				},
			},
			expected: map[string]interface{}{
				"foo": "bar",
				"baz": map[string]interface{}{
					"quux": "quuz",
				},
			},
		},
		{
			name: "slice with null values",
			input: map[string]interface{}{
				"foo": "bar",
				"baz": []map[string]interface{}{
					{
						"qux":  nil,
						"quux": "quuz",
					},
				},
			},
			expected: map[string]interface{}{
				"foo": "bar",
				"baz": []map[string]interface{}{
					{
						"quux": "quuz",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := RemoveNullsCopy(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %#v, got %#v", tc.expected, actual)
			}
		})
	}
}
