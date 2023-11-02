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

	return nil
}
