package templateversion

import (
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/templates"
)

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.TemplateVersions().
		Get(req.Context, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeTemplateVersion(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	err := templates.SetTemplateSchemaDefault(req.Context, entity)
	if err != nil {
		return err
	}

	entity.UISchema.SetUserEdited()
	_, err = h.modelClient.TemplateVersions().UpdateOne(entity).
		SetUISchema(entity.UISchema).
		SetSchemaDefaultValue(entity.SchemaDefaultValue).
		Save(req.Context)

	return err
}
