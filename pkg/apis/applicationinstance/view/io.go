package view

import (
	"context"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs

type DeleteRequest struct {
	*model.ApplicationInstanceQueryInput `uri:",inline"`
}

func (r *DeleteRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

// Batch APIs

type CollectionDeleteRequest []*model.ApplicationInstanceQueryInput

func (r CollectionDeleteRequest) Validate() error {
	if len(r) == 0 {
		return errors.New("invalid input: empty")
	}
	for _, i := range r {
		if !i.ID.Valid(0) {
			return errors.New("invalid id: blank")
		}
	}
	return nil
}

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`

	ApplicationID types.ID `query:"applicationID"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ApplicationID.Valid(0) {
		return errors.New("invalid application id: blank")
	}
	_, err := modelClient.Applications().Query().
		Where(application.ID(r.ApplicationID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Error(http.StatusNotFound, "invalid application id: not found")
	}
	return nil
}

type CollectionGetResponse = []*model.ApplicationInstanceOutput

// Extensional APIs
