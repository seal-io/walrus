package schema

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "Simple valid variables",
			input: "testdata/simple",
		},
		{
			name:  "variables with context",
			input: "testdata/with_context",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Generate(GenerateOption{
				Dir: tc.input,
			})
			assert.NoError(t, err)

			actual, err := os.ReadFile(filepath.Join(tc.input, "schema.yaml"))
			assert.NoError(t, err)

			expected, err := os.ReadFile(filepath.Join(tc.input, "expected.yaml"))
			assert.NoError(t, err)

			assert.Equal(t, expected, actual)
		})
	}
}
