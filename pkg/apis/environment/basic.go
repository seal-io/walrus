package environment

import (
	"net/http"

	"github.com/seal-io/walrus/pkg/auths/session"
	envbus "github.com/seal-io/walrus/pkg/bus/environment"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/log"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	return createEnvironment(req.Context, h.modelClient, req.Model())
}

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	query := h.modelClient.Environments().Query().
		Where(environment.ID(req.ID))

	if req.IncludeSummary {
		query.WithResources()
	}

	entity, err := query.
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

	return exposeEnvironment(entity), nil
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

	if req.IncludeSummary {
		query.WithResources()
	}

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

	return exposeEnvironments(entities), cnt, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	// Validate whether the subject has permission to delete environments.
	sj := session.MustGetSubject(req.Context)
	if !sj.IsAdmin() {
		for i := range ids {
			ress := []session.ActionResource{
				{Name: "projects", Refer: req.Project.ID.String()},
				{Name: "environments", Refer: ids[i].String()},
			}

			if sj.Enforce(http.MethodDelete, ress, "") {
				continue
			}

			return errorx.HttpErrorf(http.StatusForbidden,
				"cannot delete environment %s that type not in: %v",
				ids[i], sj.ApplicableEnvironmentTypes)
		}
	}

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
