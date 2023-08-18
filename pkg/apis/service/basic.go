package service

import (
	"context"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/service"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/pkg/deployer"
	deployertf "github.com/seal-io/walrus/pkg/deployer/terraform"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	pkgservice "github.com/seal-io/walrus/pkg/service"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	dp, err := h.getDeployer(req.Context)
	if err != nil {
		return nil, err
	}

	if err = pkgservice.SetSubjectID(req.Context, entity); err != nil {
		return nil, err
	}

	createOpts := pkgservice.Options{
		TlsCertified: h.tlsCertified,
		Tags:         req.RemarkTags,
	}

	return pkgservice.Create(
		req.Context,
		h.modelClient,
		dp,
		entity,
		createOpts,
	)
}

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.Services().Query().
		Where(service.ID(req.ID)).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldID,
				templateversion.FieldName,
				templateversion.FieldVersion)
		}).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName)
		}).
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(environment.FieldName)
		}).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	return model.ExposeService(entity), nil
}

func (h Handler) Delete(req DeleteRequest) (err error) {
	if req.WithoutCleanup {
		// Do not clean deployed native resources.
		return h.modelClient.Services().DeleteOneID(req.ID).
			Exec(req.Context)
	}

	dp, err := h.getDeployer(req.Context)
	if err != nil {
		return err
	}

	destroyOpts := pkgservice.Options{
		TlsCertified: h.tlsCertified,
	}

	return pkgservice.Destroy(
		req.Context,
		h.modelClient,
		dp,
		req.Model(),
		destroyOpts)
}

func (h Handler) CollectionCreate(req CollectionCreateRequest) (CollectionCreateResponse, error) {
	entities := req.Model()

	err := h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		if err := pkgservice.SetSubjectID(req.Context, entities...); err != nil {
			return err
		}

		_, err := pkgservice.CreateScheduledServices(req.Context, tx, entities)

		return err
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to create services")
	}

	return model.ExposeServices(entities), nil
}

var (
	queryFields = []string{
		service.FieldName,
	}
	getFields = service.WithoutFields(
		service.FieldUpdateTime)
	sortFields = []string{
		service.FieldName,
		service.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.Services().Query().
		Where(service.EnvironmentID(req.Environment.ID))

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(service.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		t, err := topic.Subscribe(modelchange.Service)
		if err != nil {
			return nil, 0, err
		}

		defer func() { t.Unsubscribe() }()

		for {
			var event topic.Event

			event, err = t.Receive(stream)
			if err != nil {
				return nil, 0, err
			}

			dm, ok := event.Data.(modelchange.Event)
			if !ok {
				continue
			}

			var items []*model.ServiceOutput

			switch dm.Type {
			case modelchange.EventTypeCreate, modelchange.EventTypeUpdate:
				entities, err := query.Clone().
					Where(service.IDIn(dm.IDs...)).
					// Must append environment ID.
					Select(service.FieldEnvironmentID).
					// Must extract template.
					Select(service.FieldTemplateID).
					WithTemplate(func(tvq *model.TemplateVersionQuery) {
						tvq.Select(
							templateversion.FieldID,
							templateversion.FieldName,
							templateversion.FieldVersion)
						if req.WithSchema {
							tvq.Select(templateversion.FieldSchema)
						}
					}).
					Unique(false).
					All(stream)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeServices(entities)
			case modelchange.EventTypeDelete:
				items = make([]*model.ServiceOutput, len(dm.IDs))
				for i := range dm.IDs {
					items[i] = &model.ServiceOutput{
						ID: dm.IDs[i],
					}
				}
			}

			if len(items) == 0 {
				continue
			}

			resp := runtime.TypedResponse(dm.Type.String(), items)
			if err = stream.SendJSON(resp); err != nil {
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

	if orders, ok := req.Sorting(sortFields, model.Desc(service.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Must append environment ID.
		Select(service.FieldEnvironmentID).
		// Must extract template.
		Select(service.FieldTemplateID).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldID,
				templateversion.FieldName,
				templateversion.FieldVersion)
			if req.WithSchema {
				tvq.Select(templateversion.FieldSchema)
			}
		}).
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeServices(entities), cnt, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		if req.WithoutCleanup {
			// Do not clean deployed native resources.
			_, err := tx.Services().Delete().
				Where(service.IDIn(ids...)).
				Exec(req.Context)

			return err
		}

		services, err := tx.Services().Query().
			Where(service.IDIn(ids...)).
			All(req.Context)
		if err != nil {
			return err
		}

		deployerOpts := deptypes.CreateOptions{
			Type:        deployertf.DeployerType,
			ModelClient: tx,
			KubeConfig:  h.kubeConfig,
		}

		dp, err := deployer.Get(req.Context, deployerOpts)
		if err != nil {
			return err
		}

		destroyOpts := pkgservice.Options{
			TlsCertified: h.tlsCertified,
		}

		for _, s := range services {
			err = pkgservice.Destroy(req.Context, tx, dp, s, destroyOpts)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (h Handler) getDeployer(ctx context.Context) (deptypes.Deployer, error) {
	return deployer.Get(ctx, deptypes.CreateOptions{
		Type:        deployertf.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	})
}
