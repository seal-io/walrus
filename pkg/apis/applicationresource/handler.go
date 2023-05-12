package applicationresource

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

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
	return "ApplicationResource"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

// Batch APIs.

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
		applicationresource.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.ApplicationResources().Query().
		Where(
			applicationresource.InstanceID(req.InstanceID),
			applicationresource.ModeNEQ(types.ApplicationResourceModeDiscovered))
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

	if orders, ok := req.Sorting(sortFields, model.Desc(applicationresource.FieldCreateTime)); ok {
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
	t, err := topic.Subscribe(datamessage.ApplicationResource)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()

	query := h.modelClient.ApplicationResources().Query()
	if req.InstanceID != "" {
		query.Where(applicationresource.InstanceID(req.InstanceID))
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
				query.Clone().Where(applicationresource.IDIn(dm.Data...)),
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

	op, err := platform.GetOperator(ctx, operator.CreateOptions{Connector: *res.Edges.Connector})
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

	op, err := platform.GetOperator(ctx, operator.CreateOptions{Connector: *res.Edges.Connector})
	if err != nil {
		return err
	}

	if err = op.IsConnected(ctx); err != nil {
		return fmt.Errorf("unreachable connector: %w", err)
	}

	opts := operator.LogOptions{
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

	op, err := platform.GetOperator(ctx, operator.CreateOptions{Connector: *res.Edges.Connector})
	if err != nil {
		return err
	}

	if err = op.IsConnected(ctx); err != nil {
		return fmt.Errorf("unreachable connector: %w", err)
	}

	ts := asTermStream(ctx, req.Width, req.Height)
	defer func() { _ = ts.Close() }()
	opts := operator.ExecOptions{
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
	query *model.ApplicationResourceQuery,
	withoutKeys bool,
) (view.CollectionGetResponse, error) {
	logger := log.WithName("api").WithName("application-resource")

	// Allow returning without sorting keys.
	entities, err := query.Unique(false).
		// Must extract connector.
		Select(applicationresource.FieldConnectorID).
		WithConnector(func(cq *model.ConnectorQuery) {
			cq.Select(
				connector.FieldName,
				connector.FieldType,
				connector.FieldCategory,
				connector.FieldConfigVersion,
				connector.FieldConfigData)
		}).
		// Must extract components.
		WithComponents(func(rq *model.ApplicationResourceQuery) {
			rq.Select(getFields...).
				Order(model.Desc(applicationresource.FieldCreateTime)).
				Where(applicationresource.Mode(types.ApplicationResourceModeDiscovered))
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Expose resources.
	resp := make(view.CollectionGetResponse, len(entities))
	for i := 0; i < len(entities); i++ {
		resp[i].Resource = model.ExposeApplicationResource(entities[i])
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
			op, err := platform.GetOperator(ctx, operator.CreateOptions{Connector: *c})
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
