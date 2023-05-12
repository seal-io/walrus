package application

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/application/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/topic"
)

func Handle(mc model.ClientSet, kc *rest.Config, tc bool) Handler {
	return Handler{
		modelClient:  mc,
		kubeConfig:   kc,
		tlsCertified: tc,
	}
}

type Handler struct {
	modelClient  model.ClientSet
	kubeConfig   *rest.Config
	tlsCertified bool
}

func (h Handler) Kind() string {
	return "Application"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	entity := req.Model()

	err := h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		creates, err := dao.ApplicationCreates(tx, entity)
		if err != nil {
			return err
		}
		entity, err = creates[0].Save(ctx)

		return err
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeApplication(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Applications().DeleteOne(req.Model()).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	entity := req.Model()

	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		updates, err := dao.ApplicationUpdates(tx, entity)
		if err != nil {
			return err
		}

		return updates[0].Exec(ctx)
	})
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	return h.getApplicationOutput(ctx, req.ID)
}

func (h Handler) getApplicationOutput(
	ctx context.Context,
	id types.ID,
) (*model.ApplicationOutput, error) {
	entity, err := h.modelClient.Applications().Query().
		Where(application.ID(id)).
		// Must extract modules.
		WithModules(func(rq *model.ApplicationModuleRelationshipQuery) {
			rq.Order(model.Asc(applicationmodulerelationship.FieldCreateTime)).
				Select(
					applicationmodulerelationship.FieldApplicationID,
					applicationmodulerelationship.FieldName,
					applicationmodulerelationship.FieldModuleID,
					applicationmodulerelationship.FieldVersion,
					applicationmodulerelationship.FieldAttributes).
				// Allow returning without sorting keys.
				Unique(false)
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return model.ExposeApplication(entity), nil
}

func (h Handler) Stream(ctx runtime.RequestUnidiStream, req view.StreamRequest) error {
	t, err := topic.Subscribe(datamessage.Application)
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

		dm, ok := event.Data.(datamessage.Message[oid.ID])
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
				entity, err := h.getApplicationOutput(ctx, id)
				if err != nil {
					return err
				}
				streamData = view.StreamResponse{
					Type:       dm.Type,
					Collection: []*model.ApplicationOutput{entity},
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

// Batch APIs.

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.Applications().DeleteOne(req[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		return
	})
}

var (
	queryFields = []string{
		application.FieldName,
	}
	getFields  = application.WithoutFields(application.FieldUpdateTime)
	sortFields = []string{
		application.FieldName,
		application.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.Applications().Query().
		Where(application.ProjectIDIn(req.ProjectIDs...))
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

	if orders, ok := req.Sorting(sortFields, model.Desc(application.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := h.getCollectionQuery(query).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeApplications(entities), cnt, nil
}

func (h Handler) getCollectionQuery(query *model.ApplicationQuery) *model.ApplicationQuery {
	// Get application with instances and environments.
	return query.
		// Allow returning without sorting keys.
		Unique(false).
		// Must extract application instances.
		WithInstances(func(iq *model.ApplicationInstanceQuery) {
			iq.Select(
				applicationinstance.FieldApplicationID,
				applicationinstance.FieldID,
				applicationinstance.FieldName,
				applicationinstance.FieldStatus).
				Where(func(s *sql.Selector) {
					// Sq generate instance with row number.
					sq := s.Clone().
						AppendSelectExprAs(
							sql.RowNumber().
								PartitionBy(applicationinstance.FieldApplicationID).
								OrderBy(sql.Desc(applicationinstance.FieldCreateTime)),
							"row_number",
						).
						Where(s.P()).
						From(s.Table()).
						As(applicationinstance.Table)

					// Query latest 5 instances.
					s.Where(sql.LTE(s.C("row_number"), 5)).
						From(sq)
				}).
				Select(
					applicationinstance.FieldEnvironmentID, // Must extract environment.
				).
				WithEnvironment(func(eq *model.EnvironmentQuery) {
					eq.Select(environment.FieldName)
				})
		})
}

func (h Handler) CollectionStream(
	ctx runtime.RequestUnidiStream,
	req view.CollectionStreamRequest,
) error {
	t, err := topic.Subscribe(datamessage.Application)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()

	query := h.modelClient.Applications().Query().
		WithProject(func(q *model.ProjectQuery) {
			q.Select(project.FieldID)
		})
	if len(req.ProjectIDs) != 0 {
		query.Where(application.ProjectIDIn(req.ProjectIDs...))
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

		var streamData view.StreamResponse

		switch dm.Type {
		case datamessage.EventCreate, datamessage.EventUpdate:
			entities, err := h.getCollectionQuery(query.Clone()).
				Where(application.IDIn(dm.Data...)).
				All(ctx)
			if err != nil {
				return err
			}
			streamData = view.StreamResponse{
				Type:       dm.Type,
				Collection: model.ExposeApplications(entities),
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
