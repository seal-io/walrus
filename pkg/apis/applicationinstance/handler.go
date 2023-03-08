package applicationinstance

import (
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/applicationinstance/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/environment"
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
	return "ApplicationInstance"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.ApplicationInstances().DeleteOne(req.Model()).Exec(ctx)
}

// Batch APIs

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.ApplicationInstances().DeleteOne(req[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return
	})
}

var (
	getFields = applicationinstance.WithoutFields(
		applicationinstance.FieldApplicationID,
		applicationinstance.FieldUpdateTime)
	sortFields = []string{
		applicationinstance.FieldName,
		applicationinstance.FieldCreateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ApplicationID(req.ApplicationID))

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
	if orders, ok := req.Sorting(sortFields, model.Desc(applicationinstance.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		// allow returning without sorting keys.
		Unique(false).
		// must extract environment.
		Select(applicationinstance.FieldEnvironmentID).
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(environment.FieldName)
		}).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeApplicationInstances(entities), cnt, nil
}

// Extensional APIs
