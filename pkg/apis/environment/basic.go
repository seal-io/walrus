package environment

import (
	envbus "github.com/seal-io/walrus/pkg/bus/environment"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	pkgservice "github.com/seal-io/walrus/pkg/service"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/log"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	err := h.modelClient.WithTx(req.Context, func(tx *model.Tx) (err error) {
		entity, err = tx.Environments().Create().
			Set(entity).
			SaveE(req.Context, dao.EnvironmentConnectorsEdgeSave)
		if err != nil {
			return err
		}

		// TODO(thxCode): move the following codes into DAO.

		serviceInputs := make(model.Services, 0, len(req.Services))

		for _, s := range req.Services {
			svc := s.Model()
			svc.ProjectID = entity.ProjectID
			svc.EnvironmentID = entity.ID
			serviceInputs = append(serviceInputs, svc)
		}

		if err = pkgservice.SetSubjectID(req.Context, serviceInputs...); err != nil {
			return err
		}

		services, err := pkgservice.CreateScheduledServices(req.Context, tx, serviceInputs)
		if err != nil {
			return err
		}

		entity.Edges.Services = services

		return envbus.NotifyIDs(req.Context, tx, envbus.EventCreate, entity.ID)
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to create environment")
	}

	return model.ExposeEnvironment(entity), nil
}

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.Environments().Query().
		Where(environment.ID(req.ID)).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName)
		}).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			rq.Order(model.Desc(environmentconnectorrelationship.FieldCreateTime)).
				Select(environmentconnectorrelationship.FieldEnvironmentID).
				Unique(false).
				Select(environmentconnectorrelationship.FieldConnectorID).
				WithConnector(
					func(cq *model.ConnectorQuery) {
						cq.Select(
							connector.FieldID,
							connector.FieldType,
							connector.FieldName,
							connector.FieldProjectID)
					})
		}).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	return model.ExposeEnvironment(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		err := tx.Environments().UpdateOne(entity).
			Set(entity).
			ExecE(req.Context, dao.EnvironmentConnectorsEdgeSave)
		if err != nil {
			return err
		}

		return envbus.NotifyIDs(req.Context, tx, envbus.EventUpdate, req.ID)
	})
}

func (h Handler) Delete(req DeleteRequest) error {
	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		entity, err := dao.GetEnvironmentByID(req.Context, tx, req.ID)
		if err != nil {
			return err
		}

		err = tx.Environments().DeleteOneID(req.ID).
			Exec(req.Context)
		if err != nil {
			return err
		}

		if err = envbus.Notify(req.Context, tx, envbus.EventDelete, model.Environments{entity}); err != nil {
			// Proceed on clean up failure.
			log.Warnf("environment post deletion hook failed: %v", err)
		}

		return nil
	})
}

var (
	queryFields = []string{
		environment.FieldName,
	}
	getFields = environment.WithoutFields(
		environment.FieldUpdateTime)
	sortFields = []string{
		environment.FieldName,
		environment.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.Environments().Query().
		Where(environment.ProjectID(req.Project.ID))

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().Count(req.Context)
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

	if orders, ok := req.Sorting(sortFields, model.Desc(environment.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Must append project ID.
		Select(environment.FieldProjectID).
		// Must extract connectors.
		Select(environment.FieldID).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			// Includes connectors.
			rq.Order(model.Desc(environmentconnectorrelationship.FieldCreateTime)).
				WithConnector(func(cq *model.ConnectorQuery) {
					cq.Select(
						connector.FieldID,
						connector.FieldType,
						connector.FieldName,
						connector.FieldProjectID)
				})
		}).
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeEnvironments(entities), cnt, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		entities, err := dao.GetEnvironmentsByIDs(req.Context, tx, ids...)
		if err != nil {
			return err
		}

		_, err = tx.Environments().Delete().
			Where(environment.IDIn(ids...)).
			Exec(req.Context)
		if err != nil {
			return err
		}

		if err = envbus.Notify(req.Context, tx, envbus.EventDelete, entities); err != nil {
			// Proceed on clean up failure.
			log.Warnf("environment post deletion hook failed: %v", err)
		}

		return nil
	})
}
