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
	*model.TemplateVersionQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.TemplateVersionOutput

// Batch APIs.

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.TemplateVersion, templateversion.OrderOption] `query:",inline"`

	TemplateIDs []string `query:"templateID"`
}

func (r *CollectionGetRequest) Validate() error {
	if len(r.TemplateIDs) == 0 {
		return errors.New("invalid request: missing template id")
	}

	for _, id := range r.TemplateIDs {
		if id == "" {
			return errors.New("invalid template id: blank")
		}
	}

	return nil
}

type CollectionGetResponse = []*model.TemplateVersionOutput

// Extensional APIs.
