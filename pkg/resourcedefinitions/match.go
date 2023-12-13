package resourcedefinitions

import "github.com/seal-io/walrus/pkg/dao/model"

func Match(
	matchingRules []*model.ResourceDefinitionMatchingRule,
	projectName, environmentName, environmentType string,
	environmentLabels, resourceLabels map[string]string,
) *model.ResourceDefinitionMatchingRule {
	for _, rule := range matchingRules {
		switch {
		case rule.Selector.ProjectName != "" && rule.Selector.ProjectName != projectName:
			continue
		case rule.Selector.EnvironmentName != "" && rule.Selector.EnvironmentName != environmentName:
			continue
		case rule.Selector.EnvironmentType != "" && rule.Selector.EnvironmentType != environmentType:
			continue
		case !matchLabels(rule.Selector.EnvironmentLabels, environmentLabels):
			continue
		case !matchLabels(rule.Selector.ResourceLabels, resourceLabels):
			continue
		default:
			return rule
		}
	}

	return nil
}

// MatchEnvironment returns the matching rule that pairs with the environment regardless of resource labels.
func MatchEnvironment(
	matchingRules []*model.ResourceDefinitionMatchingRule,
	projectName, environmentName, environmentType string,
	environmentLabels map[string]string,
) *model.ResourceDefinitionMatchingRule {
	for _, rule := range matchingRules {
		switch {
		case rule.Selector.ProjectName != "" && rule.Selector.ProjectName != projectName:
			continue
		case rule.Selector.EnvironmentName != "" && rule.Selector.EnvironmentName != environmentName:
			continue
		case rule.Selector.EnvironmentType != "" && rule.Selector.EnvironmentType != environmentType:
			continue
		case !matchLabels(rule.Selector.EnvironmentLabels, environmentLabels):
			continue
		default:
			return rule
		}
	}

	return nil
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
