package templates

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/template"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

// CreateTemplate creates or updates a template.
func CreateTemplate(ctx context.Context, mc model.ClientSet, entity *model.Template) (*model.Template, error) {
	if entity == nil {
		return nil, errors.New("template is nil")
	}

	q := mc.Templates().Query().
		Where(template.Name(entity.Name))

	if entity.CatalogID.Valid() {
		q.Where(template.CatalogID(entity.CatalogID))
	}

	find, err := q.Only(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if find != nil {
		return find, nil
	}

	status.TemplateStatusInitialized.Unknown(entity, "Initializing template")
	entity.Status.SetSummary(status.WalkTemplate(&entity.Status))

	entity, err = mc.Templates().Create().
		Set(entity).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return entity, nil
}
