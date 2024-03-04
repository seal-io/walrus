package templateversion

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
)

func (h Handler) RouteReset(req RouteResetRequest) error {
	entity, err := h.modelClient.TemplateVersions().Query().
		Where(templateversion.ID(req.ID)).
		Only(req.Context)
	if err != nil {
		return err
	}

	if entity.UiSchema.IsUserEdited() {
		entity.UiSchema.UnsetUserEdited()
	}

	_, err = h.modelClient.TemplateVersions().UpdateOne(entity).
		SetUiSchema(entity.OriginalUISchema).
		Save(req.Context)
	if err != nil {
		return fmt.Errorf("error unset ui schema tags: %w", err)
	}

	return nil
}
