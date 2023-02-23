package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs

type ConnectorCreateRequest struct {
	*model.Connector `json:",inline"`
}

type ConnectorUpdateRequest struct {
	UriID types.ID `uri:"id"`

	*model.Connector `json:",inline"`
}

func (r *ConnectorUpdateRequest) Validate() error {
	if !r.UriID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	r.ID = r.UriID
	return nil
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

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`
}

type CollectionGetResponse = []*model.Connector

// Extensional APIs
