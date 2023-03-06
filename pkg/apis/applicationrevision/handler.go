package applicationrevision

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/applicationrevision/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platformtf"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

func Handle(mc model.ClientSet, kubeconfig *rest.Config) Handler {
	return Handler{
		modelClient: mc,
		kubeconfig:  kubeconfig,
	}
}

type Handler struct {
	modelClient model.ClientSet
	kubeconfig  *rest.Config
}

func (h Handler) Kind() string {
	return "ApplicationRevision"
}

// Basic APIs

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	var entity, err = h.modelClient.ApplicationRevisions().
		Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeApplicationRevision(entity), nil
}

// Basic APIs

var (
	getFields  = applicationrevision.Columns
	sortFields = []string{applicationrevision.FieldCreateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.ApplicationRevisions().Query()

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}
	if orders, ok := req.Sorting(sortFields); ok {
		query.Order(orders...)
	}
	entities, err := query.
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeApplicationRevisions(entities), cnt, nil
}

// Extensional APIs

// GetTerraformStates get the terraform states of the application revision deployment.
func (h Handler) GetTerraformStates(ctx *gin.Context, req view.GetTerraformStatesRequest) (view.GetTerraformStatesResponse, error) {
	get, err := h.modelClient.ApplicationRevisions().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if get.Output == "" {
		return nil, nil
	}

	return view.GetTerraformStatesResponse(get.Output), nil
}

// UpdateTerraformStates update the terraform states of the application revision deployment.
func (h Handler) UpdateTerraformStates(ctx *gin.Context, req view.UpdateTerraformStatesRequest) error {
	revision, err := h.modelClient.
		ApplicationRevisions().
		UpdateOneID(req.ID).
		SetOutput(string(req.RawMessage)).
		Save(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			statusMessage := revision.StatusMessage + err.Error()

			// timeout context
			updateCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, updateErr := h.modelClient.
				ApplicationRevisions().
				UpdateOneID(req.ID).
				SetStatus(status.ApplicationRevisionStatusFailed).
				SetStatusMessage(statusMessage).
				Save(updateCtx)

			if updateErr != nil {
				log.Errorf("update application revision status failed: %v", updateErr)
			}
		}
	}()

	// TODO (alex) using Topic to parse the application resources
	var parser platformtf.Parser
	applicationResources, err := parser.ParseAppRevision(revision)
	if err != nil {
		return err
	}

	existResourceIDs := make([]types.ID, 0)
	newResources := make(model.ApplicationResources, 0)

	// fetch the old resources of the application
	oldResources, err := h.modelClient.ApplicationResources().
		Query().
		Where(applicationresource.ApplicationID(revision.ApplicationID)).
		All(ctx)
	if err != nil {
		return err
	}
	oldResourceSet := sets.NewString()
	for _, r := range oldResources {
		uniqueKey := strs.Join("-", string(r.ConnectorID), r.Name, r.Type, r.Module, r.Mode)
		oldResourceSet.Insert(uniqueKey)
	}

	for _, ar := range applicationResources {
		// check if the resource is exists
		key := strs.Join("-", string(ar.ConnectorID), ar.Name, ar.Type, ar.Module, ar.Mode)
		exists := oldResourceSet.Has(key)
		if exists {
			existResourceIDs = append(existResourceIDs, ar.ID)
		} else {
			newResources = append(newResources, ar)
		}
	}

	// diff application resource of this revision and the latest revision
	// if the resource is not in the latest revision, delete it
	_, err = h.modelClient.ApplicationResources().
		Delete().
		Where(
			applicationresource.ApplicationID(revision.ApplicationID),
			applicationresource.IDNotIn(existResourceIDs...),
		).
		Exec(ctx)
	if err != nil {
		return err
	}

	// create newResource
	if len(newResources) > 0 {
		resourcesToCreate, err := dao.ApplicationResourceCreates(h.modelClient, newResources...)
		if err != nil {
			return err
		}
		if _, err = h.modelClient.ApplicationResources().CreateBulk(resourcesToCreate...).Save(ctx); err != nil {
			return err
		}
	}

	return nil
}
