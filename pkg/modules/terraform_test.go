package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/seal-io/seal/pkg/dao/types"
)

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
						Name:    "bar",
						Type:    "string",
						Default: "bar",
						Label:   "Bar Label",
						Options: []string{"B1", "B2", "B3"},
						Group:   "Test Group",
						ShowIf:  "foo=F1",
					},
					{
						Name:    "foo",
						Type:    "string",
						Default: "foo",
						Label:   "Foo Label",
						Options: []string{"F1", "F2", "F3"},
						Group:   "Test Group",
					},
					{
						Name:    "thee",
						Type:    "string",
						Default: "thee",
					},
				},
			},
			expectedError: false,
		},
	}

	// Absolute path is required for getter.Get(dest, local_source)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get pwd: %v", err)
	}
	for _, tc := range testCases {
		caseMessage := fmt.Sprintf("test case: %s", tc.name)
		actualOutput, actualError := loadTerraformModuleSchema(filepath.Join(pwd, tc.input))

		if tc.expectedError {
			assert.Error(t, actualError, caseMessage)
		} else {
			assert.NoError(t, actualError, caseMessage)
		}

		if actualOutput == nil {
			continue
		}

		// sort to avoid random-order results
		sort.Slice(actualOutput.Outputs, func(i, j int) bool {
			return actualOutput.Outputs[i].Name < actualOutput.Outputs[j].Name
		})
		sort.Slice(actualOutput.RequiredConnectorTypes, func(i, j int) bool {
			return actualOutput.RequiredConnectorTypes[i] < actualOutput.RequiredConnectorTypes[j]
		})
		sort.Slice(actualOutput.Variables, func(i, j int) bool {
			return actualOutput.Variables[i].Name < actualOutput.Variables[j].Name
		})

		assert.Equal(t, tc.expectedOutput, actualOutput, caseMessage)

	}
}
