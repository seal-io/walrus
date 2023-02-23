package connector

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/connector/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
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
	return "Connector"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.ConnectorCreateRequest) (*model.Connector, error) {
	var creates, err = dao.ConnectorCreates(h.modelClient, req.Connector)
	if err != nil {
		return nil, err
	}

	o, err := creates[0].Save(ctx)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (h Handler) Delete(ctx *gin.Context, req view.IDRequest) error {
	return h.modelClient.Connectors().DeleteOneID(req.ID).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.ConnectorUpdateRequest) (*model.Connector, error) {
	var update, err = dao.ConnectorUpdate(h.modelClient, req.Connector)
	if err != nil {
		return nil, err
	}
	o, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (h Handler) Get(ctx *gin.Context, req view.IDRequest) (*model.Connector, error) {
	return h.modelClient.Connectors().Get(ctx, req.ID)
}

// Batch APIs

var (
	getFields  = connector.Columns
	sortFields = []string{connector.FieldID, connector.FieldCreateTime, connector.FieldUpdateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.Connectors().Query()

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

	return entities, cnt, nil
}

// Extensional APIs
