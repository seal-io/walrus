package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
)

// Basic APIs.

type GetRequest struct {
	model.TemplateVersionQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid() {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.TemplateVersionOutput

// Batch APIs.

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.TemplateVersion, templateversion.OrderOption] `query:",inline"`

	TemplateNames []string `query:"templateNames"`
}

func (r *CollectionGetRequest) Validate() error {
	if len(r.TemplateNames) == 0 {
		return errors.New("invalid request: missing template name")
	}

	for _, name := range r.TemplateNames {
		if name == "" {
			return errors.New("invalid template id: blank")
		}
	}

	return nil
}

type CollectionGetResponse = []*model.TemplateVersionOutput

// Extensional APIs.
