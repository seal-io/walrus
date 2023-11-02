package templateversion

import (
	modbus "github.com/seal-io/walrus/pkg/bus/template"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
)

func (h Handler) RouteReset(req RouteResetRequest) error {
	entity, err := h.modelClient.TemplateVersions().Query().
		Select(
			templateversion.FieldID,
			templateversion.FieldTemplateID).
		Where(templateversion.ID(req.ID)).
		WithTemplate().
		Only(req.Context)
	if err != nil {
		return err
	}

	_, err = h.modelClient.TemplateVersions().UpdateOne(entity).
		SetUiSchema(types.UISchema{}).
		Save(req.Context)
	if err != nil {
		return err
	}

	return modbus.Notify(req.Context, entity.Edges.Template)
}
