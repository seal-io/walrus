package applicationresource

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/applicationresource/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
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
	return "ApplicationResource"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

// Batch APIs

var (
	getFields = applicationresource.WithoutFields(
		applicationresource.FieldUpdateTime)
	sortFields = []string{
		applicationresource.FieldModule,
		applicationresource.FieldMode,
		applicationresource.FieldType,
		applicationresource.FieldName,
		applicationresource.FieldCreateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.ApplicationResources().Query().
		Where(applicationresource.InstanceID(req.InstanceID))

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
	if orders, ok := req.Sorting(sortFields, model.Desc(applicationresource.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		// allow returning without sorting keys.
		Unique(false).
		// must extract connector.
		Select(applicationresource.FieldConnectorID).
		WithConnector(func(cq *model.ConnectorQuery) {
			cq.Select(
				connector.FieldName,
				connector.FieldType,
				connector.FieldConfigVersion,
				connector.FieldConfigData)
		}).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var resp = make(view.CollectionGetResponse, len(entities))
	for i := 0; i < len(entities); i++ {
		resp[i].ApplicationResourceOutput = model.ExposeApplicationResource(entities[i])

		if !req.WithoutKeys {
			// fetch operable keys.
			var op operator.Operator
			op, err = platform.GetOperator(ctx,
				operator.CreateOptions{Connector: *entities[i].Edges.Connector})
			if err != nil {
				return nil, 0, err
			}
			resp[i].Keys, err = op.GetKeys(ctx, *entities[i])
			if err != nil {
				return nil, 0, err
			}
		}
	}
	return resp, cnt, nil
}

// Extensional APIs

func (h Handler) GetKeys(ctx *gin.Context, req view.GetKeysRequest) (view.GetKeysResponse, error) {
	var res = req.Entity

	var op, err = platform.GetOperator(ctx, operator.CreateOptions{Connector: *res.Edges.Connector})
	if err != nil {
		return nil, err
	}
	ok, err := op.IsConnected(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("cannot connect %s", res.Edges.Connector.Name)
	}

	return op.GetKeys(ctx, *res)
}

func (h Handler) StreamLog(ctx runtime.RequestStream, req view.StreamLogRequest) error {
	var res = req.Entity

	var op, err = platform.GetOperator(ctx, operator.CreateOptions{Connector: *res.Edges.Connector})
	if err != nil {
		return err
	}
	ok, err := op.IsConnected(ctx)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("cannot connect %s", res.Edges.Connector.Name)
	}

	var opts = operator.LogOptions{
		Out:          ctx,
		Previous:     req.Previous,
		Tail:         req.Tail,
		SinceSeconds: req.SinceSeconds,
		Timestamps:   req.Timestamps,
	}
	return op.Log(ctx, req.Key, opts)
}

func (h Handler) StreamExec(ctx runtime.RequestStream, req view.StreamExecRequest) error {
	var res = req.Entity

	op, err := platform.GetOperator(ctx, operator.CreateOptions{Connector: *res.Edges.Connector})
	if err != nil {
		return err
	}
	ok, err := op.IsConnected(ctx)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("cannot connect %s", res.Edges.Connector.Name)
	}

	var ts = asTermStream(ctx, req.Width, req.Height)
	defer func() { _ = ts.Close() }()
	var opts = operator.ExecOptions{
		Out:     ts,
		In:      ts,
		Shell:   req.Shell,
		Resizer: ts,
	}
	return op.Exec(ctx, req.Key, opts)
}
