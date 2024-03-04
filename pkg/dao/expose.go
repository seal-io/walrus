// Implements custom expose functions.

package dao

import "github.com/seal-io/walrus/pkg/dao/model"

// ExposeResourceDefinition converts the ResourceDefinition to ResourceDefinitionOutput.
func ExposeResourceDefinition(_rd *model.ResourceDefinition) *model.ResourceDefinitionOutput {
	if _rd == nil {
		return nil
	}

	// Will generate UI schema while it isn't exist, the UI schema without variables consider as exist.
	if _rd.UISchema != nil && _rd.UISchema.OpenAPISchema == nil {
		uiSchema := _rd.Schema.Expose()
		_rd.UISchema = &uiSchema
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
