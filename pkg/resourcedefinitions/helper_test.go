package resourcedefinitions

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

func Test_alignSchemas(t *testing.T) {
	testCases, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}

	loader := openapi3.NewLoader()

	for _, tc := range testCases {
		if !tc.IsDir() {
			continue
		}

		givenPaths, err := filepath.Glob(
			filepath.Join("testdata", tc.Name(), "given.*.yaml"))
		if err != nil {
			t.Fatalf("failed to glob given files: %v", err)
		}

		givens := make([]openapi3.Schemas, 0, len(givenPaths))

		for i := range givenPaths {
			givenPath := givenPaths[i]
			givenOpenAPISchema, err := loader.LoadFromFile(givenPath)
			if err != nil {
				t.Fatalf("failed to load given file %q: %v", givenPath, err)
			}

			if givenOpenAPISchema.Components == nil ||
				givenOpenAPISchema.Components.Schemas == nil {
				t.Fatalf("given file %q has no components.schemas", givenPath)
			}
			givenSchemas := givenOpenAPISchema.Components.Schemas
			givens = append(givens, givenSchemas)
		}

		expectedPath := filepath.Join("testdata", tc.Name(), "expected.yaml")
		expectedOpenAPISchema, err := loader.LoadFromFile(expectedPath)
		if err != nil {
			t.Fatalf("failed to load expected file %q: %v", expectedPath, err)
		}

		if expectedOpenAPISchema.Components == nil ||
			expectedOpenAPISchema.Components.Schemas == nil {
			t.Fatalf("expected file %q has no components.schemas", expectedPath)
		}

		expectedSchemas := expectedOpenAPISchema.Components.Schemas

		actualSchema := alignSchemas(map[string]any{}, givens)
		assert.Equal(t, expectedSchemas, actualSchema, tc.Name())
	}
}
