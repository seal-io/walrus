package templates

import (
	"context"
	"errors"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/types/status"
)

// CreateTemplate creates or updates a template.
func CreateTemplate(ctx context.Context, mc model.ClientSet, entity *model.Template) (*model.Template, error) {
	if entity == nil {
		return nil, errors.New("template is nil")
	}

	status.TemplateStatusInitialized.Unknown(entity, "Initializing template")
	entity.Status.SetSummary(status.WalkTemplate(&entity.Status))

	id, err := mc.Templates().Create().
		Set(entity).
		OnConflictColumns(template.FieldName).
		Update(func(up *model.TemplateUpsert) {
			up.UpdateStatus().
				UpdateDescription().
				UpdateIcon()
		}).
		ID(ctx)
	if err != nil {
		return nil, err
	}

	entity.ID = id

	return entity, nil
}
