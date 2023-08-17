package connector

import (
	"net/http"

	"github.com/drone/go-scm/scm"

	pkgconn "github.com/seal-io/walrus/pkg/connectors"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/utils/errorx"
)

func (h Handler) RouteApplyCostTools(
	req RouteApplyCostToolsRequest,
) error {
	entity, err := h.modelClient.Connectors().Get(req.Context, req.ID)
	if err != nil {
		return err
	}

	status.ConnectorStatusCostToolsDeployed.Unknown(entity, "Redeploying cost tools actively")

	if err = pkgconn.UpdateStatus(req.Context, h.modelClient, entity); err != nil {
		return err
	}

	err = applyFinOps(h.modelClient, entity, true)
	if err != nil {
		return err
	}

	return nil
}

func (h Handler) RouteSyncCostData(
	req RouteSyncCostDataRequest,
) error {
	entity, err := h.modelClient.Connectors().Get(req.Context, req.ID)
	if err != nil {
		return err
	}

	status.ConnectorStatusCostSynced.Unknown(entity, "Syncing cost data actively")

	if err = pkgconn.UpdateStatus(req.Context, h.modelClient, entity); err != nil {
		return err
	}

	syncer := pkgconn.NewStatusSyncer(req.Client)

	return syncer.SyncFinOpsStatus(req.Context, entity)
}

func (h Handler) RouteGetRepositories(
	req RouteGetRepositoriesRequest,
) (RouteGetRepositoriesResponse, error) {
	conn, err := h.modelClient.Connectors().Get(req.Context, req.ID)
	if err != nil {
		return nil, err
	}

	if conn.Category != types.ConnectorCategoryVersionControl {
		return nil, errorx.HttpErrorf(http.StatusBadRequest,
			"%q is not a supported version control driver", conn.Type)
	}

	client, err := vcs.NewClient(conn)
	if err != nil {
		return nil, err
	}

	listOptions := scm.ListOptions{
		IncludePrivate: true,
		Page:           req.Page,
		Size:           req.PerPage,
	}
	if req.Query != nil {
		listOptions.Search = *req.Query
	}

	repositories, _, err := client.Repositories.List(req.Context, listOptions)
	if err != nil {
		return nil, err
	}

	return repositories, nil
}

func (h Handler) RouteGetRepositoryBranches(
	req RouteGetRepositoryBranchesRequest,
) (RouteGetRepositoryBranchesResponse, error) {
	conn, err := h.modelClient.Connectors().Get(req.Context, req.ID)
	if err != nil {
		return nil, err
	}

	if conn.Category != types.ConnectorCategoryVersionControl {
		return nil, errorx.HttpErrorf(http.StatusBadRequest,
			"%q is not a supported SCM driver", conn.Type)
	}

	client, err := vcs.NewClient(conn)
	if err != nil {
		return nil, err
	}

	listOptions := scm.ListOptions{
		IncludePrivate: true,
		Page:           req.Page,
		Size:           req.PerPage,
	}
	if req.Query != nil {
		listOptions.Search = *req.Query
	}

	branches, _, err := client.Git.ListBranches(req.Context, req.Repository, listOptions)
	if err != nil {
		return nil, err
	}

	return branches, nil
}
