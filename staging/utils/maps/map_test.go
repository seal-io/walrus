package maps

import (
	"reflect"
	"testing"
)

func TestRemoveNulls(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]any
		expected map[string]any
	}{
		{
			name: "Map with no null keys",
			input: map[string]any{
				"foo": "bar",
			},
			expected: map[string]any{
				"foo": "bar",
			},
		},
		{
			name: "Map with null keys",
			input: map[string]any{
				"foo": "bar",
				"baz": nil,
				"qux": []string{},
			},
			expected: map[string]any{
				"foo": "bar",
				"qux": []string{},
			},
		},
		{
			name: "Map with null keys and nested maps",
			input: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"qux": nil,
				},
			},
			expected: map[string]any{
				"foo": "bar",
				"baz": map[string]any{},
			},
		},
		{
			name: "Map with null keys and nested maps with non-null keys",
			input: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"qux": "quux",
				},
			},
			expected: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"qux": "quux",
				},
			},
		},
		{
			name: "Map with null keys and nested maps with null keys",
			input: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"qux":  nil,
					"quux": "quuz",
				},
			},
			expected: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"quux": "quuz",
				},
			},
		},
		{
			name: "slice with null values",
			input: map[string]any{
				"foo": "bar",
				"baz": []map[string]any{
					{
						"qux":  nil,
						"quux": "quuz",
					},
				},
			},
			expected: map[string]any{
				"foo": "bar",
				"baz": []map[string]any{
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
