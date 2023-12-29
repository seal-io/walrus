package resourcedefinitions

import (
	"golang.org/x/exp/slices"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
)

const (
	// Matching scores are in order from least to most specific.
	// Environment Type < Project Label < Project Name < Environment Label < Environment Name < Resource Label
	// Higher score represents higher priority.
	scoreEnvironmentType = 1 << iota
	scoreProjectLabel
	scoreProjectName
	scoreEnvironmentLabel
	scoreEnvironmentName
	scoreResourceLabel
)

// MatchResourceDefinition takes a list of resource definitions and the resource metadata,
// and returns the matched resource definition and the matched matching rule.
// If no resource definition matches, it returns nil.
func MatchResourceDefinition(
	resourceDefinitions []*model.ResourceDefinition,
	metadata types.MatchResourceMetadata,
) (*model.ResourceDefinition, *model.ResourceDefinitionMatchingRule) {
	var (
		matchedScore      = -1
		matchedDefinition *model.ResourceDefinition
		matchedRule       *model.ResourceDefinitionMatchingRule
	)

	for _, def := range resourceDefinitions {
		for _, rule := range def.Edges.MatchingRules {
			score := computeMatchingScore(rule, metadata)
			if score < 0 {
				continue
			} else if score > matchedScore {
				matchedScore = score
				matchedRule = rule
				matchedDefinition = def
			}
		}
	}

	if matchedScore == -1 {
		return nil, nil
	}

	return matchedDefinition, matchedRule
}

func matchLabels(selectors, labels map[string]string) bool {
	if len(selectors) == 0 {
		return true
	}

	for key, value := range selectors {
		if labelValue, exists := labels[key]; !exists || labelValue != value {
			return false
		}
	}

	return true
}

// computeMatchingScore computes the matching score given a matching rule and metadata of a resource.
// If the rule does not match, it returns -1.
func computeMatchingScore(
	rule *model.ResourceDefinitionMatchingRule,
	md types.MatchResourceMetadata,
) int {
	if (len(rule.Selector.EnvironmentTypes) > 0 &&
		!slices.Contains(rule.Selector.EnvironmentTypes, md.EnvironmentType)) ||
		(len(rule.Selector.ProjectLabels) > 0 &&
			!matchLabels(rule.Selector.ProjectLabels, md.ProjectLabels)) ||
		(len(rule.Selector.ProjectNames) > 0 &&
			!slices.Contains(rule.Selector.ProjectNames, md.ProjectName)) ||
		(len(rule.Selector.EnvironmentLabels) > 0 &&
			!matchLabels(rule.Selector.EnvironmentLabels, md.EnvironmentLabels)) ||
		(len(rule.Selector.EnvironmentNames) > 0 &&
			!slices.Contains(rule.Selector.EnvironmentNames, md.EnvironmentName)) ||
		(len(rule.Selector.ResourceLabels) > 0 &&
			!matchLabels(rule.Selector.ResourceLabels, md.ResourceLabels)) {
		// Conditions are logical AND. If any condition is not met, return 0.
		return -1
	}

	score := 0

	if len(rule.Selector.EnvironmentTypes) > 0 &&
		slices.Contains(rule.Selector.EnvironmentTypes, md.EnvironmentType) {
		score += scoreEnvironmentType
	}

	if len(rule.Selector.ProjectLabels) > 0 &&
		matchLabels(rule.Selector.ProjectLabels, md.ProjectLabels) {
		score += scoreProjectLabel
	}

	if len(rule.Selector.ProjectNames) > 0 &&
		slices.Contains(rule.Selector.ProjectNames, md.ProjectName) {
		score += scoreProjectName
	}

	if len(rule.Selector.EnvironmentLabels) > 0 &&
		matchLabels(rule.Selector.EnvironmentLabels, md.EnvironmentLabels) {
		score += scoreEnvironmentLabel
	}

	if len(rule.Selector.EnvironmentNames) > 0 &&
		slices.Contains(rule.Selector.EnvironmentNames, md.EnvironmentName) {
		score += scoreEnvironmentName
	}

	if len(rule.Selector.ResourceLabels) > 0 &&
		matchLabels(rule.Selector.ResourceLabels, md.ResourceLabels) {
		score += scoreResourceLabel
	}

	return score
}
