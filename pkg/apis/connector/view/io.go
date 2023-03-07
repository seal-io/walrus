package view

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
)

// Basic APIs

type CreateRequest struct {
	*model.ConnectorCreateInput `json:",inline"`
}

func (r *CreateRequest) ValidateWith(ctx context.Context, input any) error {
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	var entity = r.Model()
	op, err := platform.GetOperator(ctx, operator.CreateOptions{
		Connector: *entity,
	})
	if err != nil {
		return fmt.Errorf("invalid connector config: %w", err)
	}
	connected, err := op.IsConnected(ctx)
	if err != nil {
		return fmt.Errorf("invalid connector: %w", err)
	}
	if !connected {
		return errors.New("invalid connector: unreachable")
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

type CollectionDeleteRequest []*model.ConnectorQueryInput

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
