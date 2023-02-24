package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/json"
)

// Basic APIs

type SafeApplicationRevision model.ApplicationRevision

func (r SafeApplicationRevision) MarshalJSON() ([]byte, error) {
	// blind sensitive fields
	r.Output = ""
	r.InputPlan = ""
	r.InputVariables = nil

	return json.Marshal(model.ApplicationRevision(r))
}

type SafeApplicationRevisions []*model.ApplicationRevision

func (r SafeApplicationRevisions) MarshalJSON() ([]byte, error) {
	// blind sensitive fields
	r2 := make([]*SafeApplicationRevision, len(r))

	for i := 0; i < len(r); i++ {
		r2[i] = (*SafeApplicationRevision)(r[i])
	}

	return json.Marshal([]*model.ApplicationRevision(r))
}

type IDRequest struct {
	ID types.ID `uri:"id"`
}

func (r *IDRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *SafeApplicationRevision

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`
}

type CollectionGetResponse = model.ApplicationRevisions

// Extensional APIs

type GetOutputsResponse = json.RawMessage

type UpdateOutputRequest struct {
	json.RawMessage

	ID types.ID `uri:"id"`
}

func (r *UpdateOutputRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}
