package connector

import (
	"net/http"

	"github.com/drone/go-scm/scm"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/connector/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/vcs"
)

// Extensional APIs for SCM connectors.

func (h Handler) RouteRepositories(
	ctx *gin.Context,
	req view.GetRepositoriesRequest,
) (view.GetRepositoriesResponse, error) {
	conn, err := h.modelClient.Connectors().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if conn.Category != types.ConnectorCategoryVersionControl {
		return nil, runtime.Errorf(http.StatusBadRequest,
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

	repositories, _, err := client.Repositories.List(ctx, listOptions)
	if err != nil {
		return nil, err
	}

	return repositories, nil
}

func (h Handler) RouteRepositoryBranches(
	ctx *gin.Context,
	req view.GetBranchesRequest,
) (view.GetBranchesResponse, error) {
	conn, err := h.modelClient.Connectors().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if conn.Category != types.ConnectorCategoryVersionControl {
		return nil, runtime.Errorf(http.StatusBadRequest,
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

	branches, _, err := client.Git.ListBranches(ctx, req.Repository, listOptions)
	if err != nil {
		return nil, err
	}

	return branches, nil
}
