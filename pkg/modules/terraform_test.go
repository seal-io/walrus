package modules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
)

func TestGetVersionedSource(t *testing.T) {
	testCases := []struct {
		name    string
		source  string
		version string
		result  string
	}{
		{
			name:    "github with subdirectory",
			source:  "github.com/foo/bar//module1",
			version: "0.0.1",
			result:  "github.com/foo/bar//module1/0.0.1",
		},
		{
			name:    "github root",
			source:  "github.com/foo/bar",
			version: "0.0.1",
			result:  "github.com/foo/bar//0.0.1",
		},
		{
			name:    "github with ref",
			source:  "github.com/foo/bar//module1?ref=dev",
			version: "0.0.1",
			result:  "github.com/foo/bar//module1/0.0.1?ref=dev",
		},
		{
			name:    "generic git",
			source:  "git::https://github.com/foo/bar.git",
			version: "0.0.1",
			result:  "git::https://github.com/foo/bar.git//0.0.1",
		},
		{
			name:    "generic git with subdirectory",
			source:  "git::https://github.com/foo/bar.git//module1",
			version: "0.0.1",
			result:  "git::https://github.com/foo/bar.git//module1/0.0.1",
		},
		{
			name:    "generic git with ref",
			source:  "git::https://github.com/foo/bar.git//module1?ref=dev",
			version: "0.0.1",
			result:  "git::https://github.com/foo/bar.git//module1/0.0.1?ref=dev",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := getVersionedSource(tc.source, tc.version)
			assert.Equal(t, tc.result, actualOutput)
		})
	}
}

func TestLoadTerraformSchema(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput *types.ModuleSchema
		expectedError  bool
	}{
		{
			name:           "Invalid",
			input:          "testdata/invalid",
			expectedOutput: &types.ModuleSchema{},
			expectedError:  true,
		},
		{
			name:  "With README.md",
			input: "testdata/with_readme",
			expectedOutput: &types.ModuleSchema{
				Readme: "# test readme",
			},
			expectedError: false,
		},
		{
			name:  "With output",
			input: "testdata/with_output",
			expectedOutput: &types.ModuleSchema{
				Outputs: []types.ModuleOutput{
					{
						Name:        "first",
						Description: "The first output.",
					},
					{
						Name:        "second",
						Description: "The second output.",
						Sensitive:   true,
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "With variable",
			input: "testdata/with_variable",
			expectedOutput: &types.ModuleSchema{
				Variables: []types.ModuleVariable{
					{
						Name:    "foo",
						Type:    "string",
						Default: "bar",
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "Full schema",
			input: "testdata/full_schema",
			expectedOutput: &types.ModuleSchema{
				Readme:                 "# test readme",
				RequiredConnectorTypes: []string{"mycloud", "null"},
				Outputs: []types.ModuleOutput{
					{
						Name:        "first",
						Description: "The first output.",
					},
					{
						Name:        "second",
						Description: "The second output.",
						Sensitive:   true,
					},
				},
				Variables: []types.ModuleVariable{
					{
						Name:    "foo",
						Type:    "string",
						Default: "foo",
						Label:   "Foo Label",
						Options: []string{"F1", "F2", "F3"},
						Group:   "Test Group",
					},
					{
						Name:    "bar",
						Type:    "string",
						Default: "bar",
						Label:   "Bar Label",
						Options: []string{"B1", "B2", "B3"},
						Group:   "Test Group",
						ShowIf:  "foo=F1",
					},
					{
						Name:    "thee",
						Type:    "string",
						Default: "thee",
					},
					{
						Name:    "subgroup1_1",
						Type:    "string",
						Default: "subgroup1_1",
						Label:   "Subgroup1_1 Label",
						Group:   "Test Subgroup/Subgroup 1",
					},
					{
						Name:    "subgroup1_2",
						Type:    "string",
						Default: "subgroup1_2",
						Label:   "Subgroup1_2 Label",
						Group:   "Test Subgroup/Subgroup 1",
					},
					{
						Name:    "subgroup2_1",
						Type:    "string",
						Default: "subgroup2_1",
						Label:   "Subgroup2_1 Label",
						Group:   "Test Subgroup/Subgroup 2",
					},
				},
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput, actualError := loadTerraformModuleSchema(tc.input)

			if tc.expectedError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}

			if actualOutput == nil {
				return
			}

			assert.Equal(t, tc.expectedOutput, actualOutput)
		})
	}
}

func TestLoadTerraformModuleVersions(t *testing.T) {
	testCases := []struct {
		name           string
		input          *model.Module
		expectedOutput []*model.ModuleVersion
		expectedError  bool
	}{
		{
			name: "versioned-modules",
			input: &model.Module{
				ID:     "mock-id",
				Source: "testdata/versioned_module",
			},
			expectedOutput: []*model.ModuleVersion{
				{
					ModuleID: "mock-id",
					Version:  "0.0.1",
					Schema: &types.ModuleSchema{
						Readme: "# Version 0.0.1",
					},
				},
				{
					ModuleID: "mock-id",
					Version:  "0.0.2",
					Schema: &types.ModuleSchema{
						Readme: "# Version 0.0.2",
					},
				},
				{
					ModuleID: "mock-id",
					Version:  "100.0.0",
					Schema: &types.ModuleSchema{
						Readme: "# Version 100.0.0",
					},
				},
			},
		},
	}

	// Absolute path is required for getter.Get(dest, local_source)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get pwd: %v", err)
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.Source = filepath.Join(pwd, tc.input.Source)
			actualOutput, actualError := loadTerraformModuleVersions(tc.input)
			if tc.expectedError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}

			assert.Equal(t, len(tc.expectedOutput), len(actualOutput))

			for i, v := range tc.expectedOutput {
				actualV := actualOutput[i]

				assert.Equal(t, actualV.ModuleID, v.ModuleID)
				assert.Equal(t, actualV.Version, v.Version)

				assert.Equal(t, actualV.Schema, v.Schema)
			}
		})
	}
}
