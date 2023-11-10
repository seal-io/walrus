package templateversion

import "github.com/seal-io/walrus/pkg/dao"

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.TemplateVersions().
		Get(req.Context, req.ID)
	if err != nil {
		return nil, err
	}

	return dao.ExposeTemplateVersion(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	_, err := h.modelClient.TemplateVersions().UpdateOne(entity).
		SetUiSchema(entity.UiSchema).
		Save(req.Context)

	return err
}
