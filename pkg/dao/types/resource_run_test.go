package types

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/seal-io/walrus/utils/json"
)

func TestPatchLeaf(t *testing.T) {
	mask := "<sensitive>"

	testCases := []struct {
		name           string
		value          any
		sensitiveValue any
		merge          bool
		expected       any
	}{
		{
			name: "sensitive value is nil",
			value: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			sensitiveValue: nil,
			expected: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "top level sensitive value",
			value: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			sensitiveValue: map[string]any{
				"key1": true,
			},
			expected: map[string]any{
				"key1": mask,
				"key2": "value2",
			},
		},
		{
			name: "nested sensitive value",
			value: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"key3": "value3",
					"key4": "value4",
				},
			},
			sensitiveValue: map[string]any{
				"key2": map[string]any{
					"key3": true,
				},
			},
			expected: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"key3": mask,
					"key4": "value4",
				},
			},
		},
		{
			name: "sensitive slice",
			value: map[string]any{
				"key1": "value1",
				"key2": []any{"value1", "value2"},
			},
			sensitiveValue: map[string]any{
				"key2": []any{false, true},
			},
			expected: map[string]any{
				"key1": "value1",
				"key2": []any{"value1", mask},
			},
		},
		{
			name: "nested sensitive slice",
			value: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"key3": "value3",
					"key4": []any{"value1", "value2", "value3"},
				},
			},
			sensitiveValue: map[string]any{
				"key2": map[string]any{
					"key4": []any{false, true},
				},
			},
			expected: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"key3": "value3",
					"key4": []any{"value1", mask, "value3"},
				},
			},
		},
		{
			name: "test merge array",
			value: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"key3": "value3",
					"key4": []any{"value1", "value2", "value3"},
				},
			},
			sensitiveValue: map[string]any{
				"key2": map[string]any{
					"key4": []any{false, true, false, true, "value5"},
				},
			},
			merge: true,
			expected: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"key3": "value3",
					"key4": []any{"value1", mask, "value3", mask, "value5"},
				},
			},
		},
		{
			name: "test merge map",
			value: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"key3": "value3",
					"key4": []any{"value1", "value2", "value3"},
				},
			},
			sensitiveValue: map[string]any{
				"key2": map[string]any{
					"key4": []any{false, true, false, true, "value5"},
					"key5": "value5",
				},
			},
			merge: true,
			expected: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"key3": "value3",
					"key4": []any{"value1", mask, "value3", mask, "value5"},
					"key5": "value5",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := patchLeaf(tc.value, tc.sensitiveValue, mask, tc.merge)
			actualJSON, _ := json.Marshal(actual)
			expectedJSON, _ := json.Marshal(tc.expected)
			if !assert.JSONEq(t, string(expectedJSON), string(actualJSON)) {
				t.Errorf("expected: %s, got: %s", string(expectedJSON), string(actualJSON))
			}
		})
	}
}

func TestPlanUnmarshalJSON(t *testing.T) {
	cases := []struct {
		filePath         string
		expectedFilePath string
	}{
		{
			filePath:         "./testdata/resourcerun/create.json",
			expectedFilePath: "./testdata/resourcerun/create-expected.json",
		},
		{
			filePath:         "./testdata/resourcerun/update.json",
			expectedFilePath: "./testdata/resourcerun/update-expected.json",
		},
		{
			filePath:         "./testdata/resourcerun/delete.json",
			expectedFilePath: "./testdata/resourcerun/delete-expected.json",
		},
	}

	for _, c := range cases {
		planData, err := os.ReadFile(c.filePath)
		if err != nil {
			t.Fatalf("failed to read plan file: %v", err)
		}

		var plan Plan

		err = json.Unmarshal(planData, &plan)
		if err != nil {
			t.Fatalf("failed to unmarshal plan: %v", err)
		}

		actual, err := json.Marshal(plan)
		if err != nil {
			t.Fatalf("failed to marshal plan: %v", err)
		}

		expectedData, err := os.ReadFile(c.expectedFilePath)
		if err != nil {
			t.Fatalf("failed to read expected plan file: %v", err)
		}
		if !assert.JSONEq(t, string(expectedData), string(actual)) {
			t.Errorf("expected: %s, got: %s", string(expectedData), string(actual))
		}
	}
}
