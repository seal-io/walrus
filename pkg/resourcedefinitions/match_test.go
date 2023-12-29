package resourcedefinitions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
)

func TestComputeMatchingScore(t *testing.T) {
	testCases := []struct {
		name     string
		rule     *model.ResourceDefinitionMatchingRule
		metadata types.MatchResourceMetadata
		expected int
	}{
		{
			name: "not-match",
			rule: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					ProjectNames:     []string{"x"},
					EnvironmentNames: []string{"y"},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:     "p",
				EnvironmentName: "e",
				EnvironmentType: "development",
			},
			expected: -1,
		},
		{
			name: "project-name",
			rule: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					ProjectNames: []string{"p"},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:     "p",
				EnvironmentName: "e",
				EnvironmentType: "development",
			},
			expected: scoreProjectName,
		},
		{
			name: "environment-name",
			rule: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					EnvironmentNames: []string{"e"},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:       "p",
				EnvironmentName:   "e",
				EnvironmentType:   "development",
				EnvironmentLabels: map[string]string{"e1": "v1"},
			},
			expected: scoreEnvironmentName,
		},
		{
			name: "project-label",
			rule: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					ProjectLabels: map[string]string{"p1": "v1"},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:       "p",
				EnvironmentName:   "e",
				EnvironmentType:   "development",
				ProjectLabels:     map[string]string{"p1": "v1"},
				EnvironmentLabels: map[string]string{"e1": "v1"},
			},
			expected: scoreProjectLabel,
		},
		{
			name: "environment-label",
			rule: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					EnvironmentLabels: map[string]string{"e1": "v1"},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:       "p",
				EnvironmentName:   "e",
				EnvironmentType:   "development",
				ProjectLabels:     map[string]string{"p1": "v1"},
				EnvironmentLabels: map[string]string{"e1": "v1"},
			},
			expected: scoreEnvironmentLabel,
		},
		{
			name: "resource-label",
			rule: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					ResourceLabels: map[string]string{"r1": "v1"},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:       "p",
				EnvironmentName:   "e",
				EnvironmentType:   "development",
				ProjectLabels:     map[string]string{"p1": "v1"},
				EnvironmentLabels: map[string]string{"e1": "v1"},
				ResourceLabels:    map[string]string{"r1": "v1"},
			},
			expected: scoreResourceLabel,
		},
		{
			name: "project-environment-label",
			rule: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					ProjectLabels:     map[string]string{"p1": "v1"},
					EnvironmentLabels: map[string]string{"e1": "v1"},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:       "p",
				EnvironmentName:   "e",
				EnvironmentType:   "development",
				ProjectLabels:     map[string]string{"p1": "v1"},
				EnvironmentLabels: map[string]string{"e1": "v1"},
				ResourceLabels:    map[string]string{"r1": "v1"},
			},
			expected: scoreProjectLabel + scoreEnvironmentLabel,
		},
		{
			name: "match-all",
			rule: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					ProjectNames:      []string{"p"},
					EnvironmentNames:  []string{"e"},
					EnvironmentTypes:  []string{"development"},
					ProjectLabels:     map[string]string{"p1": "v1"},
					EnvironmentLabels: map[string]string{"e1": "v1"},
					ResourceLabels:    map[string]string{"r1": "v1"},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:       "p",
				EnvironmentName:   "e",
				EnvironmentType:   "development",
				ProjectLabels:     map[string]string{"p1": "v1"},
				EnvironmentLabels: map[string]string{"e1": "v1"},
				ResourceLabels:    map[string]string{"r1": "v1"},
			},
			expected: scoreProjectName + scoreProjectLabel + scoreEnvironmentName +
				scoreEnvironmentLabel + scoreEnvironmentType + scoreResourceLabel,
		},
	}

	for _, tc := range testCases {
		actual := computeMatchingScore(tc.rule, tc.metadata)
		assert.Equal(t, tc.expected, actual, tc.name)
	}
}

func TestMatchResourceDefinition(t *testing.T) {
	testCases := []struct {
		name               string
		definitions        []*model.ResourceDefinition
		metadata           types.MatchResourceMetadata
		expectedDefinition *model.ResourceDefinition
		expectedRule       *model.ResourceDefinitionMatchingRule
	}{
		{
			name: "not-match",
			definitions: []*model.ResourceDefinition{
				{
					Name: "d",
					Edges: model.ResourceDefinitionEdges{
						MatchingRules: []*model.ResourceDefinitionMatchingRule{
							{
								Selector: types.Selector{
									ProjectNames: []string{"x"},
								},
							},
						},
					},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:       "p",
				EnvironmentName:   "e",
				EnvironmentType:   "development",
				ProjectLabels:     map[string]string{"p1": "v1"},
				EnvironmentLabels: map[string]string{"e1": "v1"},
				ResourceLabels:    map[string]string{"r1": "v1"},
			},
			expectedDefinition: nil,
			expectedRule:       nil,
		},
		{
			name: "match-empty-selector",
			definitions: []*model.ResourceDefinition{
				{
					Name: "d",
					Edges: model.ResourceDefinitionEdges{
						MatchingRules: []*model.ResourceDefinitionMatchingRule{
							{
								Selector: types.Selector{},
							},
						},
					},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:       "p",
				EnvironmentName:   "e",
				EnvironmentType:   "development",
				ProjectLabels:     map[string]string{"p1": "v1"},
				EnvironmentLabels: map[string]string{"e1": "v1"},
				ResourceLabels:    map[string]string{"r1": "v1"},
			},
			expectedDefinition: &model.ResourceDefinition{
				Name: "d",
				Edges: model.ResourceDefinitionEdges{
					MatchingRules: []*model.ResourceDefinitionMatchingRule{
						{
							Selector: types.Selector{},
						},
					},
				},
			},
			expectedRule: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{},
			},
		},
		{
			name: "match-definition-by-environment-name",
			definitions: []*model.ResourceDefinition{
				{
					Name: "by-project-name",
					Edges: model.ResourceDefinitionEdges{
						MatchingRules: []*model.ResourceDefinitionMatchingRule{
							{
								Selector: types.Selector{
									ProjectNames: []string{"p"},
								},
							},
						},
					},
				},
				{
					Name: "by-environment-name",
					Edges: model.ResourceDefinitionEdges{
						MatchingRules: []*model.ResourceDefinitionMatchingRule{
							{
								Selector: types.Selector{
									EnvironmentNames: []string{"e"},
								},
							},
						},
					},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:       "p",
				EnvironmentName:   "e",
				EnvironmentType:   "development",
				ProjectLabels:     map[string]string{"p1": "v1"},
				EnvironmentLabels: map[string]string{"e1": "v1"},
				ResourceLabels:    map[string]string{"r1": "v1"},
			},
			expectedDefinition: &model.ResourceDefinition{
				Name: "by-environment-name",
				Edges: model.ResourceDefinitionEdges{
					MatchingRules: []*model.ResourceDefinitionMatchingRule{
						{
							Selector: types.Selector{
								EnvironmentNames: []string{"e"},
							},
						},
					},
				},
			},
			expectedRule: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					EnvironmentNames: []string{"e"},
				},
			},
		},
		{
			name: "match-definition-by-project-and-environment-name",
			definitions: []*model.ResourceDefinition{
				{
					Name: "by-project-name",
					Edges: model.ResourceDefinitionEdges{
						MatchingRules: []*model.ResourceDefinitionMatchingRule{
							{
								Selector: types.Selector{
									ProjectNames: []string{"p"},
								},
							},
						},
					},
				},
				{
					Name: "by-environment-name",
					Edges: model.ResourceDefinitionEdges{
						MatchingRules: []*model.ResourceDefinitionMatchingRule{
							{
								Selector: types.Selector{
									EnvironmentNames: []string{"e"},
								},
							},
						},
					},
				},
				{
					Name: "by-project-and-environment-name",
					Edges: model.ResourceDefinitionEdges{
						MatchingRules: []*model.ResourceDefinitionMatchingRule{
							{
								Selector: types.Selector{
									ProjectNames:     []string{"p"},
									EnvironmentNames: []string{"e"},
								},
							},
						},
					},
				},
			},
			metadata: types.MatchResourceMetadata{
				ProjectName:       "p",
				EnvironmentName:   "e",
				EnvironmentType:   "development",
				ProjectLabels:     map[string]string{"p1": "v1"},
				EnvironmentLabels: map[string]string{"e1": "v1"},
				ResourceLabels:    map[string]string{"r1": "v1"},
			},
			expectedDefinition: &model.ResourceDefinition{
				Name: "by-project-and-environment-name",
				Edges: model.ResourceDefinitionEdges{
					MatchingRules: []*model.ResourceDefinitionMatchingRule{
						{
							Selector: types.Selector{
								ProjectNames:     []string{"p"},
								EnvironmentNames: []string{"e"},
							},
						},
					},
				},
			},
			expectedRule: &model.ResourceDefinitionMatchingRule{
				Selector: types.Selector{
					ProjectNames:     []string{"p"},
					EnvironmentNames: []string{"e"},
				},
			},
		},
	}

	for _, tc := range testCases {
		actualDefinition, actualRule := MatchResourceDefinition(
			tc.definitions,
			tc.metadata,
		)
		assert.Equal(t, tc.expectedDefinition, actualDefinition, tc.name)
		assert.Equal(t, tc.expectedRule, actualRule, tc.name)
	}
}
