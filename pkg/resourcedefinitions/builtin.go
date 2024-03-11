package resourcedefinitions

import (
	"context"
	"sort"
	"strings"

	"github.com/seal-io/walrus/pkg/bus/builtin"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinition"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	pkgtemplates "github.com/seal-io/walrus/pkg/templates"
	"github.com/seal-io/walrus/utils/log"
)

func SyncBuiltinResourceDefinitions(ctx context.Context, m builtin.BusMessage) error {
	logger := log.WithName("builtin").WithName("resource-definitions")

	mc := m.TransactionalModelClient

	ts, err := mc.Templates().Query().
		Where(template.CatalogID(m.Refer.ID)).
		All(ctx)
	if err != nil {
		return err
	}

	resourceTypeToConnectorTypes := make(map[string][]*model.Template)

	for _, t := range ts {
		labels := t.Labels

		if len(labels) == 0 {
			continue
		}

		rt, ok := labels[types.LabelWalrusResourceType]
		if !ok {
			logger.Warnf("builtin template %s missing label %s", t.Name, types.LabelWalrusResourceType)
			continue
		}

		_, ok = labels[types.LabelWalrusConnectorType]
		if !ok {
			logger.Warnf("builtin template %s missing label %s", t.Name, types.LabelWalrusConnectorType)
			continue
		}

		resourceTypeToConnectorTypes[rt] = append(resourceTypeToConnectorTypes[rt], t)
	}

	resourceDefinitions := make([]*model.ResourceDefinition, 0, len(resourceTypeToConnectorTypes))

	for res, ts := range resourceTypeToConnectorTypes {
		var definition *model.ResourceDefinition

		definition, err = newResourceDefinition(ctx, mc, res, ts)
		if err != nil {
			logger.Errorf("failed to create builtin %s resource definition: %v", res, err)
			continue
		}

		resourceDefinitions = append(resourceDefinitions, definition)

		logger.Debugf("created builtin %s resource definition", res)
	}

	err = mc.ResourceDefinitions().CreateBulk().
		Set(resourceDefinitions...).
		OnConflictColumns(resourcedefinition.FieldName).
		UpdateNewValues().
		ExecE(ctx, dao.ResourceDefinitionMatchingRulesEdgeSave)
	if err != nil {
		return err
	}

	return nil
}

func newResourceDefinition(
	ctx context.Context,
	mc model.ClientSet,
	resourceType string,
	templates []*model.Template,
) (*model.ResourceDefinition, error) {
	logger := log.WithName("builtin").WithName("resource-definitions")

	matchingRules := make([]*model.ResourceDefinitionMatchingRule, 0, len(templates))

	for _, t := range templates {
		ct := strings.ToLower(t.Labels[types.LabelWalrusConnectorType])

		m, err := newMatchingRule(ctx, mc, t.ID, ct)
		if err != nil {
			logger.Errorf("failed to create matching rule for builtin %s-%s: %v", ct, resourceType, err)
			continue
		}

		matchingRules = append(matchingRules, m)
	}

	// Sort the matchingRules to ensure the order is deterministic.
	sort.SliceStable(matchingRules, func(i, j int) bool {
		return matchingRules[i].Name < matchingRules[j].Name
	})

	bn := "builtin-" + resourceType
	rd := &model.ResourceDefinition{
		Name:        bn,
		Type:        resourceType,
		Description: "Walrus Builtin Resource Definition",
		Builtin:     true,
		Edges: model.ResourceDefinitionEdges{
			MatchingRules: matchingRules,
		},
	}

	err := pkgtemplates.SetResourceDefinitionSchemaDefault(ctx, mc, rd)
	if err != nil {
		return nil, err
	}

	err = GenSchema(ctx, mc, rd)
	if err != nil {
		return nil, err
	}

	return rd, nil
}

func newMatchingRule(
	ctx context.Context,
	mc model.ClientSet,
	templateID object.ID,
	connectorType string,
) (*model.ResourceDefinitionMatchingRule, error) {
	version, err := mc.TemplateVersions().Query().
		Where(templateversion.TemplateID(templateID)).
		Order(dao.OrderSemverVersionFunc).
		First(ctx)
	if err != nil {
		return nil, err
	}

	m := &model.ResourceDefinitionMatchingRule{
		Name:       connectorType,
		TemplateID: version.ID,
		Selector: types.Selector{EnvironmentLabels: map[string]string{
			dao.ProviderLabelPrefix + connectorType: dao.LabelValueTrue,
		}},
		Edges: model.ResourceDefinitionMatchingRuleEdges{
			Template: version,
		},
	}

	return m, nil
}
