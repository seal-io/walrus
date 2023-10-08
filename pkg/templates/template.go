package templates

import (
	"context"
	"errors"

	"entgo.io/ent/dialect/sql"

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

	var conflictOptions []sql.ConflictOption
	if entity.ProjectID == "" {
		conflictOptions = append(
			conflictOptions,
			sql.ConflictWhere(sql.P().
				IsNull(template.FieldProjectID)),
			sql.ConflictColumns(template.FieldName),
		)
	} else {
		conflictOptions = append(
			conflictOptions,
			sql.ConflictWhere(sql.P().
				NotNull(template.FieldProjectID)),
			sql.ConflictColumns(
				template.FieldName,
				template.FieldProjectID,
			),
		)
	}

	id, err := mc.Templates().Create().
		Set(entity).
		OnConflict(conflictOptions...).
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
