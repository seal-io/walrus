package serviceresource

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/apis/serviceresource/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/operator"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/log"
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
		Where(serviceresource.ModeNEQ(types.ServiceResourceModeDiscovered))

	switch {
	case req.ServiceID != "":
		query.Where(serviceresource.ServiceID(req.ServiceID))
	case req.ServiceName != "":
		query.QueryService().Where(service.Name(req.ServiceName))
	}

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

	resp, err := getCollection(ctx, query, req.WithoutKeys)
	if err != nil {
		return nil, 0, err
	}

	return resp, cnt, nil
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

	query := h.modelClient.ServiceResources().Query()
	if req.ServiceID != "" {
		query.Where(serviceresource.ServiceID(req.ServiceID))
	}

	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	for {
		var (
			event topic.Event
			resp  view.CollectionGetResponse
		)

		event, err = t.Receive(ctx)
		if err != nil {
			return err
		}

		dm, ok := event.Data.(datamessage.Message[oid.ID])
		if !ok {
			continue
		}

		var streamData view.StreamResponse

		switch dm.Type {
		case datamessage.EventCreate, datamessage.EventUpdate:
			resp, err = getCollection(ctx,
				query.Clone().Where(serviceresource.IDIn(dm.Data...)),
				req.WithoutKeys)
			if err != nil {
				return err
			}
			streamData = view.StreamResponse{
				Type:       dm.Type,
				Collection: resp,
			}
		case datamessage.EventDelete:
			streamData = view.StreamResponse{
				Type: dm.Type,
				IDs:  dm.Data,
			}
		}

		if len(streamData.IDs) == 0 && len(streamData.Collection) == 0 {
			continue
		}

		err = ctx.SendJSON(streamData)
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

func getCollection(
	ctx context.Context,
	query *model.ServiceResourceQuery,
	withoutKeys bool,
) (view.CollectionGetResponse, error) {
	logger := log.WithName("api").WithName("service-resource")

	// Allow returning without sorting keys.
	entities, err := query.Unique(false).
		// Must extract connector.
		Select(serviceresource.FieldConnectorID).
		WithConnector(func(cq *model.ConnectorQuery) {
			cq.Select(
				connector.FieldName,
				connector.FieldType,
				connector.FieldCategory,
				connector.FieldConfigVersion,
				connector.FieldConfigData)
		}).
		// Must extract components.
		WithComponents(func(rq *model.ServiceResourceQuery) {
			rq.Select(getFields...).
				Order(model.Desc(serviceresource.FieldCreateTime)).
				Where(serviceresource.Mode(types.ServiceResourceModeDiscovered))
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Expose resources.
	resp := make(view.CollectionGetResponse, len(entities))
	for i := 0; i < len(entities); i++ {
		resp[i].ServiceResourceOutput = model.ExposeServiceResource(entities[i])
	}

	// Fetch keys for each resource without error returning.
	if !withoutKeys {
		// NB(thxCode): we can safety index the connector with its pointer here,
		// as the ent can keep the connector pointer is the same between those resources related by the same connector.
		m := make(map[*model.Connector][]int)
		for i := 0; i < len(entities); i++ {
			m[entities[i].Edges.Connector] = append(m[entities[i].Edges.Connector], i)
		}

		for c, idxs := range m {
			// Get operator by connector.
			op, err := operator.Get(ctx, optypes.CreateOptions{Connector: *c})
			if err != nil {
				logger.Warnf("cannot get operator of connector: %v", err)
				continue
			}

			if err = op.IsConnected(ctx); err != nil {
				logger.Warnf("unreachable connector: %v", err)
				continue
			}
			// Fetch keys for the resources that related to same connector.
			for _, i := range idxs {
				resp[i].Keys, err = op.GetKeys(ctx, entities[i])
				if err != nil {
					logger.Errorf("error getting keys: %v", err)
				}
			}
		}
	}

	return resp, nil
}

func (h Handler) CollectionGetGraph(
	ctx *gin.Context,
	req view.CollectionGetGraphRequest,
) (*view.CollectionGetGraphResponse, error) {
	// Fetch entities.
	entities, err := h.modelClient.ServiceResources().Query().
		Order(model.Asc(serviceresource.FieldCreateTime)).
		Where(
			serviceresource.ServiceID(req.ServiceID),
			serviceresource.ModeNEQ(types.ServiceResourceModeDiscovered)).
		Select(
			serviceresource.FieldServiceID,
			serviceresource.FieldType,
			serviceresource.FieldID,
			serviceresource.FieldName,
			serviceresource.FieldCreateTime,
			serviceresource.FieldUpdateTime,
			serviceresource.FieldStatus).
		WithComponents(func(rq *model.ServiceResourceQuery) {
			rq.Order(model.Asc(serviceresource.FieldCreateTime)).
				Where(serviceresource.Mode(types.ServiceResourceModeDiscovered)).
				Select(
					serviceresource.FieldServiceID,
					serviceresource.FieldType,
					serviceresource.FieldID,
					serviceresource.FieldName,
					serviceresource.FieldCreateTime,
					serviceresource.FieldUpdateTime,
					serviceresource.FieldStatus)
		}).
		All(ctx)
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
		vertices = make([]view.GraphVertex, 0, verticesCap)
		edges    = make([]view.GraphEdge, 0, edgesCap)
	)

	for i := 0; i < len(entities); i++ {
		// Append ServiceResource to vertices.
		vertices = append(vertices, view.GraphVertex{
			GraphVertexID: view.GraphVertexID{
				Kind: "ServiceResource",
				ID:   entities[i].ID,
			},
			Name:       entities[i].Name,
			CreateTime: entities[i].CreateTime,
			UpdateTime: entities[i].UpdateTime,
			Status:     entities[i].Status.Summary,
			Extensions: map[string]any{
				"type": entities[i].Type,
			},
		})

		for j := 0; j < len(entities[i].Edges.Components); j++ {
			// Append sub ServiceResource to vertices.
			vertices = append(vertices, view.GraphVertex{
				GraphVertexID: view.GraphVertexID{
					Kind: "ServiceResource",
					ID:   entities[i].Edges.Components[j].ID,
				},
				Name:       entities[i].Edges.Components[j].Name,
				CreateTime: entities[i].Edges.Components[j].CreateTime,
				UpdateTime: entities[i].Edges.Components[j].UpdateTime,
				Status:     entities[i].Edges.Components[j].Status.Summary,
				Extensions: map[string]any{
					"type": entities[i].Edges.Components[j].Type,
				},
			})

			// Append the edge of ServiceResource to sub ServiceResource.
			edges = append(edges, view.GraphEdge{
				Type: "Composition",
				Start: view.GraphVertexID{
					Kind: "ServiceResource",
					ID:   entities[i].ID,
				},
				End: view.GraphVertexID{
					Kind: "ServiceResource",
					ID:   entities[i].Edges.Components[j].ID,
				},
			})
		}
	}

	// TODO(thxCode): return operation keys.

	return &view.CollectionGetGraphResponse{
		Vertices: vertices,
		Edges:    edges,
	}, nil
}
