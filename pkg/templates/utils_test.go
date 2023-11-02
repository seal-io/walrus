package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplateNameByPath(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "versioned template path",
			input:    "templates/foo/0.0.1",
			expected: "foo",
		},
		{
			name:     "non-versioned template path",
			input:    "templates/foo",
			expected: "foo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := GetTemplateNameByPath(tc.input)
			assert.Equal(t, tc.expected, actualOutput)
		})
	}
}
