package applicationresource

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/apis/applicationresource/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
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
	return "ApplicationResource"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Stream(ctx runtime.RequestStream, req view.StreamRequest) error {
	var t, err = topic.Subscribe(datamessage.ApplicationResource)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()
	for {
		var event topic.Event
		event, err = t.Receive(ctx)
		if err != nil {
			return err
		}
		dm, ok := event.Data.(datamessage.Message)
		if !ok {
			continue
		}

		var streamData view.StreamResponse
		for _, id := range dm.Data {
			if id != req.ID {
				continue
			}

			switch dm.Type {
			case datamessage.EventCreate, datamessage.EventUpdate:
				entity, err := h.modelClient.ApplicationResources().Get(ctx, id)
				if err != nil && !model.IsNotFound(err) {
					return err
				}
				entity.Module = parseModuleName(entity.Module)
				keys, err := getKeys(ctx, entity)
				if err != nil {
					return err
				}
				streamData = view.StreamResponse{
					Type: dm.Type,
					Collection: []view.ApplicationResource{
						{
							ApplicationResourceOutput: model.ExposeApplicationResource(entity),
							Keys:                      keys,
						},
					},
				}
			case datamessage.EventDelete:
				streamData = view.StreamResponse{
					Type: dm.Type,
					IDs:  dm.Data,
				}
			}
		}
		err = ctx.SendJSON(streamData)
		if err != nil {
			return err
		}
	}
}

// Batch APIs

var (
	queryFields = []string{
		applicationresource.FieldName,
	}
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
		Where(
			applicationresource.InstanceID(req.InstanceID),
			applicationresource.ModeNEQ(types.ApplicationResourceModeDiscovered))
	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
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
	if orders, ok := req.Sorting(sortFields, model.Desc(applicationresource.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	resp, err := getCollection(ctx, query, req.WithoutKeys)
	if err != nil {
		return nil, 0, err
	}
	return resp, cnt, nil
}

func (h Handler) CollectionStream(ctx runtime.RequestStream, req view.CollectionStreamRequest) error {
	var t, err = topic.Subscribe(datamessage.ApplicationResource)
	if err != nil {
		return err
	}

	var query = h.modelClient.ApplicationResources().Query()
	if req.InstanceID != "" {
		query.Where(applicationresource.InstanceID(req.InstanceID))
	}

	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	defer func() { t.Unsubscribe() }()
	for {
		var (
			event topic.Event
			resp  view.CollectionGetResponse
		)
		event, err = t.Receive(ctx)
		if err != nil {
			return err
		}
		dm, ok := event.Data.(datamessage.Message)
		if !ok {
			continue
		}

		var streamData view.StreamResponse
		switch dm.Type {
		case datamessage.EventCreate, datamessage.EventUpdate:
			// filter out component ids.
			entities, err := query.Clone().
				Where(applicationresource.IDIn(dm.Data...)).
				Select(
					applicationresource.FieldID,
					applicationresource.FieldCompositionID).
				All(ctx)
			if err != nil {
				return err
			}
			var ids = sets.Set[oid.ID]{}
			for i := range entities {
				if entities[i].CompositionID != "" {
					// consume composition id.
					ids.Insert(entities[i].CompositionID)
					continue
				}
				ids.Insert(entities[i].ID)
			}
			// get entities by filtered out ids.
			resp, err = getCollection(ctx,
				query.Clone().Where(applicationresource.IDIn(ids.UnsortedList()...)),
				req.WithoutKeys)
			if err != nil {
				return err
			}
			streamData = view.StreamResponse{
				Type:       dm.Type,
				IDs:        dm.Data,
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

// Extensional APIs

func (h Handler) GetKeys(ctx *gin.Context, req view.GetKeysRequest) (view.GetKeysResponse, error) {
	return getKeys(ctx, req.Entity)
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
	err = op.Exec(ctx, req.Key, opts)
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

func getCollection(ctx context.Context, query *model.ApplicationResourceQuery, withoutKeys bool) (view.CollectionGetResponse, error) {
	// allow returning without sorting keys.
	entities, err := query.Unique(false).
		// must extract connector.
		Select(applicationresource.FieldConnectorID).
		WithConnector(func(cq *model.ConnectorQuery) {
			cq.Select(
				connector.FieldName,
				connector.FieldType,
				connector.FieldCategory,
				connector.FieldConfigVersion,
				connector.FieldConfigData)
		}).
		// must extract components.
		WithComponents(func(rq *model.ApplicationResourceQuery) {
			rq.Select(getFields...).
				Order(model.Desc(applicationresource.FieldCreateTime)).
				Where(applicationresource.Mode(types.ApplicationResourceModeDiscovered))
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range entities {
		r.Module = parseModuleName(r.Module)
	}

	var resp = make(view.CollectionGetResponse, len(entities))
	for i := 0; i < len(entities); i++ {
		resp[i].ApplicationResourceOutput = model.ExposeApplicationResource(entities[i])
		if !withoutKeys {
			resp[i].Keys, err = getKeys(ctx, entities[i])
			if err != nil {
				return nil, err
			}
		}
	}
	return resp, nil
}

func getKeys(ctx context.Context, r *model.ApplicationResource) (*operator.Keys, error) {
	op, err := platform.GetOperator(ctx, operator.CreateOptions{Connector: *r.Edges.Connector})
	if err != nil {
		return nil, err
	}
	ok, err := op.IsConnected(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("cannot connect %s", r.Edges.Connector.Name)
	}
	return op.GetKeys(ctx, r)
}

// parseMoudleName parse moudle name only show application module name in UI.
func parseModuleName(moduleName string) string {
	return strings.Split(moduleName, "/")[0]
}
