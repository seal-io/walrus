package applicationrevision

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/applicationrevision/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	revisionbus "github.com/seal-io/seal/pkg/bus/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platformtf"
	tftopic "github.com/seal-io/seal/pkg/topic/platformtf"
	"github.com/seal-io/seal/utils/log"
)

func Handle(mc model.ClientSet, kc *rest.Config) Handler {
	return Handler{
		modelClient: mc,
		kubeConfig:  kc,
	}
}

type Handler struct {
	modelClient model.ClientSet
	kubeConfig  *rest.Config
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
	getFields = applicationrevision.WithoutFields(
		applicationrevision.FieldStatusMessage,
		applicationrevision.FieldInputPlan,
		applicationrevision.FieldOutput,
		applicationrevision.FieldInputVariables,
		applicationrevision.FieldModules,
	)
	sortFields = []string{
		applicationrevision.FieldCreateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.ApplicationRevisions().Query()
	if req.InstanceID != "" {
		query.Where(applicationrevision.InstanceID(req.InstanceID))
	}

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
		WithInstance(func(aiq *model.ApplicationInstanceQuery) {
			aiq.Select(
				applicationinstance.FieldID,
				applicationinstance.FieldName,
				applicationinstance.FieldApplicationID).
				WithApplication(func(aq *model.ApplicationQuery) {
					aq.Select(application.FieldID, application.FieldName)
				})
		}).
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(environment.FieldID, environment.FieldName)
		}).
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

			updateRevision, updateErr := h.modelClient.ApplicationRevisions().UpdateOne(req.Model()).
				SetStatus(status.ApplicationRevisionStatusFailed).
				SetStatusMessage(statusMessage).
				Save(updateCtx)

			if updateErr != nil {
				log.Errorf("update application revision status failed: %v", updateErr)
				return
			}

			if err = revisionbus.Notify(ctx, h.modelClient, updateRevision); err != nil {
				log.Errorf("add application revision update notify err: %w", err)
			}
		}
	}()

	if err = revisionbus.Notify(ctx, h.modelClient, entity); err != nil {
		return err
	}

	return tftopic.Notify(ctx, tftopic.Name, tftopic.TopicMessage{
		ModelClient:         h.modelClient,
		ApplicationRevision: entity,
	})
}

func (h Handler) StreamLog(ctx runtime.RequestStream, req view.StreamLogRequest) error {
	var cli, err = coreclient.NewForConfig(h.kubeConfig)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}

	return platformtf.StreamJobLogs(ctx, cli, req.ID, ctx)
}
