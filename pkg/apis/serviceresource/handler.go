package serviceresource

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/apis/serviceresource/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/operator"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	pkgresource "github.com/seal-io/seal/pkg/serviceresources"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/topic"
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
	return "ServiceResource"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

// Batch APIs.

var (
	queryFields = []string{
		serviceresource.FieldName,
	}
	getFields = serviceresource.WithoutFields(
		serviceresource.FieldUpdateTime)
	sortFields = []string{
		serviceresource.FieldMode,
		serviceresource.FieldType,
		serviceresource.FieldName,
		serviceresource.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.ServiceResources().Query().
		Where(
			serviceresource.ServiceID(req.ServiceID),
			serviceresource.Mode(types.ServiceResourceModeManaged),
			serviceresource.Shape(types.ServiceResourceShapeInstance),
		)

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}

	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	if orders, ok := req.Sorting(sortFields, model.Desc(serviceresource.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := getCollection(ctx, query, req.WithoutKeys)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeServiceResources(entities), cnt, nil
}

func (h Handler) CollectionStream(
	ctx runtime.RequestUnidiStream,
	req view.CollectionStreamRequest,
) error {
	t, err := topic.Subscribe(datamessage.ServiceResource)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()

	query := h.modelClient.ServiceResources().Query().
		Where(
			serviceresource.Mode(types.ServiceResourceModeManaged),
			serviceresource.Shape(types.ServiceResourceShapeInstance),
		)

	if req.ServiceID != "" {
		query.Where(serviceresource.ServiceID(req.ServiceID))
	}

	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	for {
		var event topic.Event

		event, err = t.Receive(ctx)
		if err != nil {
			return err
		}

		dm, ok := event.Data.(datamessage.Message[oid.ID])
		if !ok {
			continue
		}

		var resp view.CollectionStreamResponse

		switch dm.Type {
		case datamessage.EventCreate, datamessage.EventUpdate:
			var entities model.ServiceResources

			entities, err = getCollection(
				ctx,
				query.Clone().Where(serviceresource.IDIn(dm.Data...)),
				req.WithoutKeys)
			if err != nil {
				return err
			}

			resp = view.CollectionStreamResponse{
				Type:       dm.Type,
				Collection: model.ExposeServiceResources(entities),
			}
		case datamessage.EventDelete:
			resp = view.CollectionStreamResponse{
				Type: dm.Type,
				IDs:  dm.Data,
			}
		}

		if len(resp.IDs) == 0 && len(resp.Collection) == 0 {
			continue
		}

		err = ctx.SendJSON(resp)
		if err != nil {
			return err
		}
	}
}

// Extensional APIs.

func (h Handler) GetKeys(ctx *gin.Context, req view.GetKeysRequest) (view.GetKeysResponse, error) {
	res := req.Entity

	op, err := operator.Get(ctx, optypes.CreateOptions{Connector: *res.Edges.Connector})
	if err != nil {
		return nil, err
	}

	if err = op.IsConnected(ctx); err != nil {
		return nil, fmt.Errorf("unreachable connector: %w", err)
	}

	return op.GetKeys(ctx, res)
}

func (h Handler) StreamLog(ctx runtime.RequestUnidiStream, req view.StreamLogRequest) error {
	res := req.Entity

	op, err := operator.Get(ctx, optypes.CreateOptions{Connector: *res.Edges.Connector})
	if err != nil {
		return err
	}

	if err = op.IsConnected(ctx); err != nil {
		return fmt.Errorf("unreachable connector: %w", err)
	}

	opts := optypes.LogOptions{
		Out:          ctx,
		Previous:     req.Previous,
		Tail:         req.Tail,
		SinceSeconds: req.SinceSeconds,
		Timestamps:   req.Timestamps,
	}

	return op.Log(ctx, req.Key, opts)
}

func (h Handler) StreamExec(ctx runtime.RequestBidiStream, req view.StreamExecRequest) error {
	res := req.Entity

	op, err := operator.Get(ctx, optypes.CreateOptions{Connector: *res.Edges.Connector})
	if err != nil {
		return err
	}

	if err = op.IsConnected(ctx); err != nil {
		return fmt.Errorf("unreachable connector: %w", err)
	}

	ts := asTermStream(ctx, req.Width, req.Height)
	defer func() { _ = ts.Close() }()
	opts := optypes.ExecOptions{
		Out:     ts,
		In:      ts,
		Shell:   req.Shell,
		Resizer: ts,
	}

	err = op.Exec(ts, req.Key, opts)
	if err != nil {
		if strings.Contains(err.Error(), "OCI runtime exec failed: exec failed:") {
			return &websocket.CloseError{
				Code: websocket.CloseUnsupportedData,
				Text: "unresolved exec shell: " + req.Shell,
			}
		}

		return err
	}

	return nil
}

func (h Handler) CollectionGetGraph(
	ctx *gin.Context,
	req view.CollectionGetGraphRequest,
) (*view.CollectionGetGraphResponse, error) {
	fields := []string{
		serviceresource.FieldServiceID,
		serviceresource.FieldDeployerType,
		serviceresource.FieldType,
		serviceresource.FieldID,
		serviceresource.FieldName,
		serviceresource.FieldMode,
		serviceresource.FieldShape,
		serviceresource.FieldClassID,
		serviceresource.FieldCreateTime,
		serviceresource.FieldUpdateTime,
		serviceresource.FieldStatus,
	}

	// Fetch entities.
	query := h.modelClient.ServiceResources().Query().
		Select(fields...).
		Order(model.Desc(serviceresource.FieldCreateTime)).
		Where(
			serviceresource.ServiceID(req.ServiceID),
			serviceresource.Mode(types.ServiceResourceModeManaged),
			serviceresource.Shape(types.ServiceResourceShapeClass),
		)

	entities, err := dao.ServiceResourceShapeClassQuery(query, fields...).All(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate capacity for allocation.
	var verticesCap, edgesCap int
	{
		// Count the number of ServiceResource.
		verticesCap = len(entities)
		for i := 0; i < len(entities); i++ {
			// Count the vertex size of sub ServiceResource,
			// and the edge size from ServiceResource to sub ServiceResource.
			verticesCap += len(entities[i].Edges.Components)
			edgesCap += len(entities[i].Edges.Components)
		}
	}

	// Construct response.
	var (
		vertices  = make([]types.GraphVertex, 0, verticesCap)
		edges     = make([]types.GraphEdge, 0, edgesCap)
		operators = make(map[oid.ID]optypes.Operator)
	)

	pkgresource.SetKeys(ctx, entities, operators)
	vertices, edges = pkgresource.GetVerticesAndEdges(entities, vertices, edges)

	return &view.CollectionGetGraphResponse{
		Vertices: vertices,
		Edges:    edges,
	}, nil
}

func getCollection(
	ctx context.Context,
	query *model.ServiceResourceQuery,
	withoutKeys bool,
) (model.ServiceResources, error) {
	wcOpts := func(cq *model.ConnectorQuery) {
		cq.Select(
			connector.FieldName,
			connector.FieldType,
			connector.FieldCategory,
			connector.FieldConfigVersion,
			connector.FieldConfigData)
	}
	// Query service resource with its components.
	entities, err := query.Unique(false).
		// Must extract connector.
		Select(serviceresource.FieldConnectorID).
		WithConnector(wcOpts).
		// Must extract components.
		WithComponents(func(rq *model.ServiceResourceQuery) {
			rq.Select(getFields...).
				Order(model.Desc(serviceresource.FieldCreateTime)).
				Where(serviceresource.Mode(types.ServiceResourceModeDiscovered)).
				WithConnector(wcOpts)
		}).All(ctx)
	if err != nil {
		return nil, err
	}

	// Return directly if no need next operations, e.g. Log, Exec and so on.
	if withoutKeys {
		return entities, nil
	}

	operators := make(map[oid.ID]optypes.Operator)
	entities = pkgresource.SetKeys(ctx, entities, operators)

	return entities, nil
}
