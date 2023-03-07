package view

import (
	"context"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/json"
)

// Basic APIs

type GetRequest struct {
	*model.ApplicationRevisionQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetResponse = *model.ApplicationRevisionOutput

// Batch APIs

type CollectionDeleteRequest = []*model.ApplicationRevisionQueryInput

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`

	InstanceID types.ID `query:"instanceID"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.InstanceID.Valid(0) {
		return errors.New("invalid instance id: blank")
	}
	_, err := modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ID(r.InstanceID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Error(http.StatusNotFound, "invalid instance id: not found")
	}
	return nil
}

type CollectionGetResponse = []*model.ApplicationRevisionOutput

// Extensional APIs

type GetTerraformStatesRequest = GetRequest

type GetTerraformStatesResponse = json.RawMessage

type UpdateTerraformStatesRequest struct {
	GetRequest      `uri:",inline"`
	json.RawMessage `json:",inline"`
}
