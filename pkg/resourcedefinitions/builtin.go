package resourcedefinitions

import (
	"context"
	"fmt"
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
	"github.com/seal-io/walrus/pkg/templates"
	"github.com/seal-io/walrus/utils/log"
)

type templateWithConnector struct {
	templateID    object.ID
	connectorType string
}

func SyncBuiltinResourceDefinitions(ctx context.Context, m builtin.BusMessage) error {
	logger := log.WithName("builtin").WithName("resource-definitions")

	mc := m.TransactionalModelClient

	ts, err := mc.Templates().Query().
		Where(template.CatalogID(m.Refer.ID)).
		All(ctx)
	if err != nil {
		return err
	}

	resourceDefinitionToConnectorTypes := make(map[string][]templateWithConnector)

	for _, t := range ts {
		labels := t.Labels

		if len(labels) == 0 {
			continue
		}

		rdn, ok := labels[types.LabelWalrusResourceDefinitionName]
		if !ok {
			logger.Warnf("builtin template %s missing label %s", t.Name, types.LabelWalrusResourceDefinitionName)
			continue
		}

		rt, ok := labels[types.LabelWalrusResourceType]
		if !ok {
			logger.Warnf("builtin template %s missing label %s", t.Name, types.LabelWalrusResourceType)
			continue
		}

		ct, ok := labels[types.LabelWalrusConnectorType]
		if !ok {
			logger.Warnf("builtin template %s missing label %s", t.Name, types.LabelWalrusConnectorType)
			continue
		}

		key := fmt.Sprintf("%s/%s", rdn, rt)
		resourceDefinitionToConnectorTypes[key] = append(resourceDefinitionToConnectorTypes[key], templateWithConnector{
			templateID:    t.ID,
			connectorType: ct,
		})
	}

	resourceDefinitions := make([]*model.ResourceDefinition, 0, len(resourceDefinitionToConnectorTypes))

	for key, conns := range resourceDefinitionToConnectorTypes {
		keys := strings.SplitN(key, "/", 2)
		rdn := keys[0]
		rt := keys[1]

		// Sort the connector types to ensure the order of matching rules is deterministic.
		sort.Slice(conns, func(i, j int) bool {
			return conns[i].connectorType < conns[j].connectorType
		})

		var definition *model.ResourceDefinition

		definition, err = newResourceDefinition(ctx, mc, rdn, rt, conns)
		if err != nil {
			logger.Errorf("failed to create builtin %s resource definition: %v", rdn, err)
			continue
		}

		resourceDefinitions = append(resourceDefinitions, definition)

		logger.Debugf("created builtin %s resource definition", rdn)
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
	resourceDefinitionName string,
	resourceType string,
	templateWithConnector []templateWithConnector,
) (*model.ResourceDefinition, error) {
	logger := log.WithName("builtin").WithName("resource-definitions")

	matchingRules := make([]*model.ResourceDefinitionMatchingRule, 0, len(templateWithConnector))

	for _, tmpl := range templateWithConnector {
		ct := strings.ToLower(tmpl.connectorType)

		m, err := newMatchingRule(ctx, mc, tmpl.templateID, ct)
		if err != nil {
			logger.Errorf("failed to create matching rule for builtin %s-%s: %v", ct, resourceType, err)
			continue
		}

		matchingRules = append(matchingRules, m)
	}

	bn := "builtin-" + resourceDefinitionName
	rd := &model.ResourceDefinition{
		Name:        bn,
		Type:        resourceType,
		Description: "Walrus Builtin Resource Definition",
		Builtin:     true,
		Edges: model.ResourceDefinitionEdges{
			MatchingRules: matchingRules,
		},
	}

	err := templates.SetResourceDefinitionSchemaDefault(ctx, mc, rd)
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
