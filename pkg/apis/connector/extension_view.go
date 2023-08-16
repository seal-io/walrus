package connector

import (
	"context"
	"errors"

	"github.com/drone/go-scm/scm"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type RouteApplyCostToolsRequest struct {
	_ struct{} `route:"POST=/apply-cost-tools"`

	model.ConnectorQueryInput `path:",inline"`
}

func (r *RouteApplyCostToolsRequest) Validate() error {
	if err := r.ConnectorQueryInput.Validate(); err != nil {
		return err
	}

	return validateConnectorType(r.Context, r.Client, r.ID)
}

type RouteSyncCostDataRequest struct {
	_ struct{} `route:"POST=/sync-cost-data"`

	model.ConnectorQueryInput `path:",inline"`
}

func (r *RouteSyncCostDataRequest) Validate() error {
	if err := r.ConnectorQueryInput.Validate(); err != nil {
		return err
	}

	return validateConnectorType(r.Context, r.Client, r.ID)
}

type (
	RouteGetRepositoriesRequest struct {
		_ struct{} `route:"GET=/repositories"`

		model.ConnectorQueryInput `path:",inline"`

		Query   *string `query:"query,omitempty"`
		Page    int     `query:"page,default=1"`
		PerPage int     `query:"perPage,default=10"`
	}

	RouteGetRepositoriesResponse = []*scm.Repository
)

type (
	RouteGetRepositoryBranchesRequest struct {
		_ struct{} `route:"GET=/repository-branches"`

		model.ConnectorQueryInput `path:",inline"`

		Repository string  `query:"repository"`
		Query      *string `query:"query,omitempty"`
		Page       int     `query:"page,default=1"`
		PerPage    int     `query:"perPage,default=10"`
	}

	RouteGetRepositoryBranchesResponse = []*scm.Reference
)

func validateConnectorType(ctx context.Context, modelClient model.ClientSet, id object.ID) error {
	conn, err := modelClient.Connectors().Query().
		Select(connector.FieldType).
		Where(connector.ID(id)).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get connector")
	}

	if conn.Type != types.ConnectorTypeK8s {
		return errors.New("invalid type: not support")
	}

	return nil
}
