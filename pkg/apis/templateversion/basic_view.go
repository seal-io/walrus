package templateversion

import (
	"github.com/seal-io/walrus/pkg/dao/model"
)

type (
	GetRequest = model.TemplateVersionQueryInput

	GetResponse = *model.TemplateVersionOutput
)

type UpdateRequest struct {
	model.TemplateVersionUpdateInput `path:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.TemplateVersionUpdateInput.Validate(); err != nil {
		return err
	}

	if r.TemplateVersionUpdateInput.Schema.IsEmpty() {
		tv, err := r.Client.TemplateVersion.Get(r.Context, r.ID)
		if err != nil {
			return err
		}

		r.TemplateVersionUpdateInput.Schema = tv.Schema
	}

	return nil
}
