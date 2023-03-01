package view

import (
	"context"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs

type CreateRequest struct {
	*model.ConnectorCreateInput `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	return nil
}

type CreateResponse = *model.ConnectorOutput

type DeleteRequest = GetRequest

type UpdateRequest struct {
	*model.ConnectorUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type UpdateResponse = *model.ConnectorOutput

type GetRequest struct {
	*model.ConnectorQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetResponse = *model.ConnectorOutput

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`
}

type CollectionGetResponse = []*model.ConnectorOutput

// Extensional APIs

type ApplyCostToolsRequest struct {
	_ struct{} `route:"POST=/apply-cost-tools"`

	ID types.ID `uri:"id"`
}

func (r *ApplyCostToolsRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return validateConnectorType(ctx, modelClient, r.ID)
}

type SyncCostDataRequest struct {
	_ struct{} `route:"POST=/sync-cost-data"`

	ID types.ID `uri:"id"`
}

func (r *SyncCostDataRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return validateConnectorType(ctx, modelClient, r.ID)
}

func validateConnectorType(ctx context.Context, modelClient model.ClientSet, id types.ID) error {
	conn, err := modelClient.Connectors().Get(ctx, id)
	if err != nil {
		return err
	}

	if conn.Type != types.ConnectorTypeK8s {
		return runtime.Errorf(http.StatusBadRequest, "invalid type: not support")
	}
	return nil
}
