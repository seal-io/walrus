package resource

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	pkgenv "github.com/seal-io/walrus/pkg/environment"
	"github.com/seal-io/walrus/utils/json"
)

// GenComputedAttributes generate computed attributes for resource.
// Required:
// entity.EnvironmentID, entity.Attributes
// entity.TemplateID/entity.ResourceDefinitionMatchingRuleID.
func GenComputedAttributes(
	ctx context.Context,
	modelClient model.ClientSet,
	entity *model.Resource,
) (property.Values, error) {
	// Check.
	if (entity.TemplateID == nil || !entity.TemplateID.Valid()) &&
		(entity.ResourceDefinitionMatchingRuleID == nil || !entity.ResourceDefinitionMatchingRuleID.Valid()) {
		return nil, fmt.Errorf("failed to generate computed attributes, " +
			"both template and resource definition matching rule id are empty")
	}

	if !entity.EnvironmentID.Valid() {
		return nil, fmt.Errorf("failed to generate computed attributes, environment id is empty")
	}

	// Compute.
	var computedAttrs property.Values

	env, err := modelClient.Environments().Query().
		Where(environment.ID(entity.EnvironmentID)).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(
				project.FieldID,
				project.FieldName)
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Walrus Context.
	wctx := *types.NewContext().
		SetProject(env.Edges.Project.ID, env.Edges.Project.Name).
		SetEnvironment(env.ID, env.Name, pkgenv.GetManagedNamespaceName(env)).
		SetResource(entity.ID, entity.Name)

	switch {
	case entity.TemplateID != nil:
		// Get template version.
		tv, err := modelClient.TemplateVersions().Query().
			Where(templateversion.ID(*entity.TemplateID)).
			Select(
				templateversion.FieldID,
				templateversion.FieldName,
				templateversion.FieldSchemaDefaultValue).
			Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get template version: %w", err)
		}

		computedAttrs, err = computedAttributeWithTemplate(wctx, entity.Attributes, tv)
		if err != nil {
			return nil, err
		}

	case entity.ResourceDefinitionMatchingRuleID != nil:
		rule, err := modelClient.ResourceDefinitionMatchingRules().Query().
			Where(resourcedefinitionmatchingrule.ID(*entity.ResourceDefinitionMatchingRuleID)).
			Select(
				resourcedefinitionmatchingrule.FieldTemplateID,
				resourcedefinitionmatchingrule.FieldResourceDefinitionID,
				resourcedefinitionmatchingrule.FieldAttributes,
				resourcedefinitionmatchingrule.FieldSchemaDefaultValue,
			).
			WithTemplate(func(tq *model.TemplateVersionQuery) {
				tq.Select(
					templateversion.FieldID,
					templateversion.FieldSchemaDefaultValue,
				)
			}).
			Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get resource definition matching rule: %w", err)
		}

		computedAttrs, err = computedAttributeWithResourceDefinition(wctx, entity.Attributes, rule)
		if err != nil {
			return nil, err
		}
	}

	return computedAttrs, nil
}

// computedAttributeWithTemplate generate computed attribute from template.
func computedAttributeWithTemplate(
	wctx types.Context,
	attrs property.Values,
	template *model.TemplateVersion,
) (property.Values, error) {
	var (
		err       error
		wctxByte  []byte
		attrsByte []byte
	)

	wctxByte, err = json.Marshal(map[string]any{"context": wctx})
	if err != nil {
		return nil, err
	}

	if len(attrs) != 0 {
		attrsByte, err = json.Marshal(attrs)
		if err != nil {
			return nil, err
		}
	}

	merged, err := json.ApplyPatches(wctxByte, template.SchemaDefaultValue, attrsByte)
	if err != nil {
		return nil, err
	}

	var ca property.Values

	err = json.Unmarshal(merged, &ca)
	if err != nil {
		return nil, err
	}

	return ca, nil
}

// computedAttributeWithResourceDefinition computed attribute with resource definition.
// required: rule.Edges.Template
func computedAttributeWithResourceDefinition(
	wctx types.Context,
	attrs property.Values,
	rule *model.ResourceDefinitionMatchingRule,
) (property.Values, error) {
	// Check.
	if rule.Edges.Template == nil {
		return nil, fmt.Errorf("edge template is empty")
	}

	var (
		err             error
		wctxByte        []byte
		attrsByte       []byte
		rdSchemaDefault = rule.SchemaDefaultValue
		tvSchemaDefault = rule.Edges.Template.SchemaDefaultValue
	)

	wctxByte, err = json.Marshal(map[string]any{"context": wctx})
	if err != nil {
		return nil, err
	}

	if len(attrs) != 0 {
		attrsByte, err = json.Marshal(attrs)
		if err != nil {
			return nil, err
		}
	}

	merged, err := json.ApplyPatches(
		wctxByte,
		tvSchemaDefault,
		rdSchemaDefault,
		attrsByte)
	if err != nil {
		return nil, err
	}

	var ca property.Values

	err = json.Unmarshal(merged, &ca)
	if err != nil {
		return nil, err
	}

	return ca, nil
}
