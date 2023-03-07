package applicationrevision

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"

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

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
	}
}

type Handler struct {
	modelClient model.ClientSet
}

func (h Handler) Kind() string {
	return "ApplicationRevision"
}

func (h Handler) Validating() any {
	return h.modelClient
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

// Batch APIs

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.ApplicationRevisions().DeleteOne(req[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return
	})
}

var (
	getFields = []string{
		applicationrevision.FieldID,
		applicationrevision.FieldStatus,
		applicationrevision.FieldStatusMessage,
		applicationrevision.FieldCreateTime,
	}
	sortFields = []string{
		applicationrevision.FieldCreateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.InstanceID(req.InstanceID))

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
	if orders, ok := req.Sorting(sortFields, model.Desc(applicationrevision.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		Unique(false). // allow returning without sorting keys.
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeApplicationRevisions(entities), cnt, nil
}

// Extensional APIs

// GetTerraformStates get the terraform states of the application revision deployment.
func (h Handler) GetTerraformStates(ctx *gin.Context, req view.GetTerraformStatesRequest) (view.GetTerraformStatesResponse, error) {
	var entity, err = h.modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.ID(req.ID)).
		Select(applicationrevision.FieldOutput).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if entity.Output == "" {
		return nil, nil
	}
	return view.GetTerraformStatesResponse(entity.Output), nil
}

// UpdateTerraformStates update the terraform states of the application revision deployment.
func (h Handler) UpdateTerraformStates(ctx *gin.Context, req view.UpdateTerraformStatesRequest) error {
	var entity, err = h.modelClient.ApplicationRevisions().UpdateOne(req.Model()).
		SetOutput(string(req.RawMessage)).
		Save(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			statusMessage := entity.StatusMessage + err.Error()

			// timeout context
			updateCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, updateErr := h.modelClient.ApplicationRevisions().UpdateOne(req.Model()).
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
	applicationResources, err := parser.ParseAppRevision(entity)
	if err != nil {
		return err
	}

	existResourceIDs := make([]types.ID, 0)
	newResources := make(model.ApplicationResources, 0)

	// fetch the old resources of the application
	oldResources, err := h.modelClient.ApplicationResources().
		Query().
		Where(applicationresource.InstanceID(entity.InstanceID)).
		All(ctx)
	if err != nil {
		return err
	}
	oldResourceSet := sets.NewString()
	for _, r := range oldResources {
		oldResourceSet.Insert(getFingerprint(r))
	}

	for _, ar := range applicationResources {
		// check if the resource is exists
		exists := oldResourceSet.Has(getFingerprint(ar))
		if exists {
			existResourceIDs = append(existResourceIDs, ar.ID)
		} else {
			newResources = append(newResources, ar)
		}
	}

	// diff application resource of this revision and the latest revision
	// if the resource is not in the latest revision, delete it
	_, err = h.modelClient.ApplicationResources().Delete().
		Where(
			applicationresource.InstanceID(entity.InstanceID),
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
		_, err = h.modelClient.ApplicationResources().CreateBulk(resourcesToCreate...).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO(thxCode): generate by entc.
func getFingerprint(r *model.ApplicationResource) string {
	// align to schema definition.
	return strs.Join("-", string(r.ConnectorID), r.Module, r.Mode, r.Type, r.Name)
}
