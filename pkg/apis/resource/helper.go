package resource

import (
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
)

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
