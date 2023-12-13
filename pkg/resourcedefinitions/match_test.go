package resourcedefinitions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
)

func TestMatchEnvironment(t *testing.T) {
	testCases := []struct {
		name              string
		matchRules        []*model.ResourceDefinitionMatchingRule
		projectName       string
		environmentName   string
		environmentType   string
		environmentLabels map[string]string
		expected          *model.ResourceDefinitionMatchingRule
	}{
		{
			name: "match-all",
			matchRules: []*model.ResourceDefinitionMatchingRule{
				{
					Selector: types.Selector{
						ProjectName: "p-x",
					},
				},
				{
					Selector: types.Selector{
						ProjectName:     "p",
						EnvironmentName: "e",
						EnvironmentType: "development",
						EnvironmentLabels: map[string]string{
							"l1": "v1",
						},
					},
				},
			},
			projectName:       "p",
			environmentName:   "e",
			environmentType:   "development",
			environmentLabels: map[string]string{"l1": "v1"},
			expected: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					ProjectName:     "p",
					EnvironmentName: "e",
					EnvironmentType: "development",
					EnvironmentLabels: map[string]string{
						"l1": "v1",
					},
				},
			},
		},
		{
			name: "match-partial",
			matchRules: []*model.ResourceDefinitionMatchingRule{
				{
					Selector: types.Selector{
						ProjectName: "p-x",
					},
				},
				{
					Selector: types.Selector{
						EnvironmentType: "development",
						EnvironmentLabels: map[string]string{
							"l1": "v1",
						},
					},
				},
			},
			projectName:       "p",
			environmentName:   "e",
			environmentType:   "development",
			environmentLabels: map[string]string{"l1": "v1"},
			expected: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					EnvironmentType: "development",
					EnvironmentLabels: map[string]string{
						"l1": "v1",
					},
				},
			},
		},
		{
			name: "match-none",
			matchRules: []*model.ResourceDefinitionMatchingRule{
				{
					Selector: types.Selector{
						ProjectName: "p-x",
						EnvironmentLabels: map[string]string{
							"l1": "v1",
						},
					},
				},
				{
					Selector: types.Selector{
						EnvironmentType: "development",
						EnvironmentLabels: map[string]string{
							"l2": "v2",
						},
					},
				},
			},
			projectName:       "p",
			environmentName:   "e",
			environmentType:   "development",
			environmentLabels: map[string]string{"l3": "v3"},
			expected:          nil,
		},
	}

	for _, tc := range testCases {
		actual := MatchEnvironment(
			tc.matchRules,
			tc.projectName,
			tc.environmentName,
			tc.environmentType,
			tc.environmentLabels,
		)
		assert.Equal(t, tc.expected, actual, tc.name)
	}
}
