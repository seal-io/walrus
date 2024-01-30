package resourcedefinitions

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"

	"github.com/seal-io/walrus/utils/json"
)

func Test_alignSchemas(t *testing.T) {
	dir := filepath.Join("testdata", "align_schemas")

	testCases, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}

	loader := openapi3.NewLoader()

	for _, tc := range testCases {
		if !tc.IsDir() {
			continue
		}

		var givenSchemasSlice []openapi3.Schemas
		{
			givenPaths, err := filepath.Glob(
				filepath.Join(dir, tc.Name(), "given.*.yaml"))
			if err != nil {
				t.Fatalf("failed to glob given files: %v", err)
			}

			givenSchemasSlice = make([]openapi3.Schemas, 0, len(givenPaths))

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
				givenSchemasSlice = append(givenSchemasSlice, givenSchemas)
			}
		}

		var expectedSchemas openapi3.Schemas
		{
			expectedPath := filepath.Join(dir, tc.Name(), "expected.yaml")
			expectedOpenAPISchema, err := loader.LoadFromFile(expectedPath)
			if err != nil {
				t.Fatalf("failed to load expected file %q: %v", expectedPath, err)
			}

			if expectedOpenAPISchema.Components == nil ||
				expectedOpenAPISchema.Components.Schemas == nil {
				t.Fatalf("expected file %q has no components.schemas", expectedPath)
			}

			expectedSchemas = expectedOpenAPISchema.Components.Schemas
		}

		nb := map[string]any{
			variablesSchemaKey: map[string]any{},
			outputsSchemaKey:   map[string]any{},
		}
		actualSchemas := alignSchemas(nb, givenSchemasSlice)
		assert.Equal(t, expectedSchemas, actualSchemas, tc.Name())
	}
}

func Test_refillVariableSchemaRef(t *testing.T) {
	dir := filepath.Join("testdata", "refill_variable_schema_ref")

	testCases, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}

	loader := openapi3.NewLoader()

	for _, tc := range testCases {
		if !tc.IsDir() {
			continue
		}

		var givenSchemas openapi3.Schemas
		{
			givenPath := filepath.Join(dir, tc.Name(), "given.yaml")
			givenOpenAPISchema, err := loader.LoadFromFile(givenPath)
			if err != nil {
				t.Fatalf("failed to load given file %q: %v", givenPath, err)
			}

			if givenOpenAPISchema.Components == nil ||
				givenOpenAPISchema.Components.Schemas == nil {
				t.Fatalf("given file %q has no components.schemas", givenPath)
			}
			givenSchemas = givenOpenAPISchema.Components.Schemas
		}

		var expectedSchemas openapi3.Schemas
		{
			expectedPath := filepath.Join(dir, tc.Name(), "expected.yaml")
			expectedOpenAPISchema, err := loader.LoadFromFile(expectedPath)
			if err != nil {
				t.Fatalf("failed to load expected file %q: %v", expectedPath, err)
			}

			if expectedOpenAPISchema.Components == nil ||
				expectedOpenAPISchema.Components.Schemas == nil {
				t.Fatalf("expected file %q has no components.schemas", expectedPath)
			}

			expectedSchemas = expectedOpenAPISchema.Components.Schemas
		}

		var nb map[string]any
		{
			nbJSONPath := filepath.Join(dir, tc.Name(), "nb.json")

			nbJSONBytes, err := os.ReadFile(nbJSONPath)
			if err != nil {
				t.Fatalf("failed to read nb.json: %v", err)
			}

			err = json.Unmarshal(nbJSONBytes, &nb)
			if err != nil {
				t.Fatalf("failed to unmarshal nb.json: %v", err)
			}
		}

		var defs [][]byte
		{
			defPaths, err := filepath.Glob(
				filepath.Join(dir, tc.Name(), "default.*.json"))
			if err != nil {
				t.Fatalf("failed to glob default files: %v", err)
			}

			defs = make([][]byte, 0, len(defPaths))

			for i := range defPaths {
				defPath := defPaths[i]
				defJSONBytes, err := os.ReadFile(defPath)
				if err != nil {
					t.Fatalf("failed to read %s: %v", defPath, err)
				}

				defs = append(defs, defJSONBytes)
			}
		}

		if sr, ok := givenSchemas[variablesSchemaKey]; ok {
			refillVariableSchemaRef(nb[variablesSchemaKey].(map[string]any), "", sr, defs, "", sr.Value)
		}

		actualSchemas := givenSchemas
		actualJSONBytes := json.ShouldMarshalIndent(actualSchemas, "", "  ")
		expectedJSONBytes := json.ShouldMarshalIndent(expectedSchemas, "", "  ")

		assert.Equal(t, string(expectedJSONBytes), string(actualJSONBytes), tc.Name())
	}
}
