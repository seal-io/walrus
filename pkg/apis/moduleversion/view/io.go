package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// Basic APIs

type GetRequest struct {
	*model.ModuleVersionQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetResponse = *model.ModuleVersionOutput

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.ModuleVersion] `query:",inline"`

	ModuleID []string `query:"moduleID"`
}

func (r *CollectionGetRequest) Validate() error {
	if len(r.ModuleID) == 0 {
		return errors.New("invalid request: missing module id")
	}

	for _, id := range r.ModuleID {
		if id == "" {
			return errors.New("invalid module id: blank")
		}
	}
	return nil
}

type CollectionGetResponse = []*model.ModuleVersionOutput

// Extensional APIs
