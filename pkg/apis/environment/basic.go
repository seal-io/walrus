package environment

import (
	"context"
	"net/http"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/auths/session"
	envbus "github.com/seal-io/walrus/pkg/bus/environment"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/pkg/deployer"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	pkgresource "github.com/seal-io/walrus/pkg/resources"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	dp, err := h.getDeployer(req.Context)
	if err != nil {
		return nil, err
	}

	return createEnvironment(req.Context, h.modelClient, req.Model(), pkgresource.Options{
		Deployer:            dp,
		Draft:               false,
		RunApprovalRequired: req.ApprovalRequired,
	})
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

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(environment.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		t, err := topic.Subscribe(modelchange.Resource)
		if err != nil {
			return nil, 0, err
		}

		defer func() { t.Unsubscribe() }()

		for {
			event, err := t.Receive(stream)
			if err != nil {
				return nil, 0, err
			}

			dm, ok := event.Data.(modelchange.Event)
			if !ok {
				continue
			}

			entities, err := query.Clone().
				Where(environment.IDIn(dm.EnvironmentIDs()...)).
				// Must append project ID.
				Select(environment.FieldProjectID).
				// Must extract connectors.
				Select(environment.FieldID).
				WithResources(func(rq *model.ResourceQuery) {
					rq.Select(
						resource.FieldID,
						resource.FieldStatus,
						resource.FieldEnvironmentID,
					)
				}).
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
				All(stream)
			if err != nil {
				return nil, 0, err
			}

			resp := runtime.TypedResponse(modelchange.EventTypeUpdate.String(), exposeEnvironments(entities))
			if err := stream.SendJSON(resp); err != nil {
				return nil, 0, err
			}
		}
	}

	// Handle normal request.

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
		WithResources(func(rq *model.ResourceQuery) {
			rq.Select(
				resource.FieldID,
				resource.FieldStatus,
				resource.FieldEnvironmentID,
			)
		}).
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

func (h Handler) getDeployer(ctx context.Context) (deptypes.Deployer, error) {
	dep, err := deployer.Get(ctx, deptypes.CreateOptions{
		Type:       types.DeployerTypeTF,
		KubeConfig: h.kubeConfig,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to get deployer")
	}

	return dep, nil
}
