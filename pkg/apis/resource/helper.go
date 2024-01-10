package resource

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinition"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/deployer/terraform"
	pkgenv "github.com/seal-io/walrus/pkg/environment"
	"github.com/seal-io/walrus/pkg/resourcedefinitions"
	"github.com/seal-io/walrus/utils/json"
)

// genComputedAttributes generate computed attributes for resource.
func genComputedAttributes(
	ctx context.Context,
	entity *model.Resource,
	modelClient *model.Client,
) (property.Values, error) {
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
	wctx := *terraform.NewContext().
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

	case entity.ResourceDefinitionID != nil:
		rd, err := modelClient.ResourceDefinitions().Query().
			Where(resourcedefinition.ID(*entity.ResourceDefinitionID)).
			WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
				rq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
					Select(
						resourcedefinitionmatchingrule.FieldTemplateID,
						resourcedefinitionmatchingrule.FieldResourceDefinitionID,
						resourcedefinitionmatchingrule.FieldAttributes,
						resourcedefinitionmatchingrule.FieldSchemaDefaultValue).
					Unique(false).
					WithTemplate(func(tq *model.TemplateVersionQuery) {
						tq.Select(
							templateversion.FieldID,
							templateversion.FieldSchemaDefaultValue,
						)
					})
			}).
			Select(
				resourcedefinition.FieldID,
				resourcedefinition.FieldName).
			Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get resource definition: %w", err)
		}

		rule := resourcedefinitions.Match(
			rd.Edges.MatchingRules,
			env.Edges.Project.Name,
			env.Name,
			env.Type,
			env.Labels,
			entity.Labels,
		)
		if rule == nil {
			return nil, fmt.Errorf("resource definition %s does not match environment %s", rd.Name, env.Name)
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
	wctx terraform.Context,
	attrs property.Values,
	template *model.TemplateVersion,
) (property.Values, error) {
	wctxByte, err := json.Marshal(map[string]any{
		"context": wctx,
	})
	if err != nil {
		return nil, err
	}

	attrsByte, err := json.Marshal(attrs)
	if err != nil {
		return nil, err
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
	wctx terraform.Context,
	attrs property.Values,
	rule *model.ResourceDefinitionMatchingRule,
) (property.Values, error) {
	tvSchemaDefault := rule.Edges.Template.SchemaDefaultValue
	rdSchemaDefault := rule.SchemaDefaultValue

	wctxByte, err := json.Marshal(map[string]any{
		"context": wctx,
	})
	if err != nil {
		return nil, err
	}

	attrsBytes, err := json.Marshal(attrs)
	if err != nil {
		return nil, err
	}

	merged, err := json.ApplyPatches(
		wctxByte,
		tvSchemaDefault,
		rdSchemaDefault,
		attrsBytes)
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

func createInputsItemToResource(
	input *model.ResourceCreateInputsItem,
	p *model.ProjectQueryInput,
	e *model.EnvironmentQueryInput,
) *model.Resource {
	return toResource(
		nil, input.Name, input.Type, input.Attributes, input.Labels,
		input.Template, input.ResourceDefinition,
		p, e)
}

func updateInputsItemToResource(
	typ string,
	input *model.ResourceUpdateInputsItem,
	p *model.ProjectQueryInput,
	e *model.EnvironmentQueryInput,
) *model.Resource {
	return toResource(
		nil, input.Name, typ, input.Attributes, input.Labels,
		input.Template, input.ResourceDefinition,
		p, e)
}

func toResource(
	id *object.ID, name, typ string, attr property.Values, labels map[string]string,
	tmpl *model.TemplateVersionQueryInput, rd *model.ResourceDefinitionQueryInput,
	p *model.ProjectQueryInput, e *model.EnvironmentQueryInput,
) *model.Resource {
	r := &model.Resource{
		Name:       name,
		Labels:     labels,
		Type:       typ,
		Attributes: attr,
	}

	if id != nil {
		r.ID = *id
	}

	if tmpl != nil {
		r.TemplateID = &tmpl.ID
	}

	if rd != nil {
		r.ResourceDefinitionID = &rd.ID
	}

	if p != nil {
		r.ProjectID = p.ID
	}

	if e != nil {
		r.EnvironmentID = e.ID
	}

	return r
}
