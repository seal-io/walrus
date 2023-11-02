package templateversion

import (
	"github.com/seal-io/walrus/pkg/dao/model"
)

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.TemplateVersions().
		Get(req.Context, req.ID)
	if err != nil {
		return nil, err
	}

	return exposeTemplateVersion(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	_, err := h.modelClient.TemplateVersions().UpdateOne(entity).
		SetUiSchema(entity.UiSchema).
		Save(req.Context)

	return err
}

func exposeTemplateVersion(
	entity *model.TemplateVersion,
) *model.TemplateVersionOutput {
	// Set expose schema.
	if !entity.UiSchema.IsEmpty() || entity.Schema.IsEmpty() {
		return model.ExposeTemplateVersion(entity)
	}

	entity.UiSchema = entity.Schema.Expose()

	return model.ExposeTemplateVersion(entity)
}
