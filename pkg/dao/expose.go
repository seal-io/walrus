// Implements custom expose functions.

package dao

import "github.com/seal-io/walrus/pkg/dao/model"

// ExposeResourceDefinition converts the ResourceDefinition to ResourceDefinitionOutput.
func ExposeResourceDefinition(_rd *model.ResourceDefinition) *model.ResourceDefinitionOutput {
	if _rd == nil {
		return nil
	}

	if _rd.UiSchema != nil && _rd.UiSchema.IsEmpty() {
		uiSchema := _rd.Schema.Expose()
		_rd.UiSchema = &uiSchema
	}

	return model.ExposeResourceDefinition(_rd)
}

// ExposeResourceDefinitions converts the ResourceDefinition slice to ResourceDefinitionOutput pointer slice.
func ExposeResourceDefinitions(_rds []*model.ResourceDefinition) []*model.ResourceDefinitionOutput {
	if len(_rds) == 0 {
		return nil
	}

	rdos := make([]*model.ResourceDefinitionOutput, len(_rds))
	for i := range _rds {
		rdos[i] = ExposeResourceDefinition(_rds[i])
	}

	return rdos
}

// ExposeTemplateVersion converts the TemplateVersion to TemplateVersionOutput.
func ExposeTemplateVersion(_tv *model.TemplateVersion) *model.TemplateVersionOutput {
	if _tv == nil {
		return nil
	}

	if _tv.UiSchema.IsEmpty() {
		_tv.UiSchema = _tv.Schema.Expose()
	}

	return model.ExposeTemplateVersion(_tv)
}

// ExposeTemplateVersions converts the TemplateVersion slice to TemplateVersionOutput pointer slice.
func ExposeTemplateVersions(_tvs []*model.TemplateVersion) []*model.TemplateVersionOutput {
	if len(_tvs) == 0 {
		return nil
	}

	tvos := make([]*model.TemplateVersionOutput, len(_tvs))
	for i := range _tvs {
		tvos[i] = ExposeTemplateVersion(_tvs[i])
	}

	return tvos
}
