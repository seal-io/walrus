package resource

import (
	"context"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	pkgrun "github.com/seal-io/walrus/pkg/resourceruns"
	pkgresource "github.com/seal-io/walrus/pkg/resources"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return nil, err
	}

	entity, _, err = pkgresource.Create(
		req.Context,
		h.modelClient,
		entity,
		pkgresource.Options{
			Deployer: dp,
			Draft:    req.Draft,
			Preview:  req.Preview,
		},
	)

	return model.ExposeResource(entity), err
}

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.Resources().Query().
		Where(resource.ID(req.ID)).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldID,
				templateversion.FieldTemplateID,
				templateversion.FieldName,
				templateversion.FieldVersion,
				templateversion.FieldProjectID)
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

	return model.ExposeResource(entity), nil
}

func (h Handler) Delete(req DeleteRequest) (err error) {
	err = h.cleanResourcePlanFiles(req.Context, h.modelClient, req.ID)
	if err != nil {
		return err
	}

	if req.WithoutCleanup {
		// Do not clean deployed native resources.
		return h.modelClient.Resources().DeleteOneID(req.ID).
			Exec(req.Context)
	}

	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return err
	}

	deleteOptions := pkgresource.DeleteOptions{
		Options: pkgresource.Options{
			Deployer: dp,
			Preview:  req.Preview,
		},
		WithoutCleanup: req.WithoutCleanup,
	}

	return pkgresource.Delete(
		req.Context,
		h.modelClient,
		req.Model(),
		deleteOptions)
}

func (h Handler) Patch(req PatchRequest) (*PatchResponse, error) {
	entity := req.Model()

	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return nil, err
	}

	run, err := pkgresource.Upgrade(req.Context, h.modelClient, entity, pkgresource.Options{
		Deployer:      dp,
		ChangeComment: req.ChangeComment,
		Draft:         req.Draft,
		Preview:       req.Preview,
	})
	if err != nil {
		return nil, err
	}

	return &PatchResponse{
		ResourceOutput: model.ExposeResource(entity),
		Run:            model.ExposeResourceRun(run),
	}, nil
}

func (h Handler) CollectionCreate(req CollectionCreateRequest) (CollectionCreateResponse, error) {
	entities := req.Model()

	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return nil, err
	}

	entities, err = pkgresource.CollectionCreate(req.Context, h.modelClient, entities, pkgresource.Options{
		Deployer: dp,
		Draft:    req.Draft,
		Preview:  req.Preview,
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeResources(entities), nil
}

var (
	queryFields = []string{
		resource.FieldName,
	}
	getFields = resource.WithoutFields(
		resource.FieldUpdateTime)
	sortFields = []string{
		resource.FieldName,
		resource.FieldCreateTime,
		resource.FieldType,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.Resources().Query().
		Where(resource.EnvironmentID(req.Environment.ID))

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(resource.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		t, err := topic.Subscribe(modelchange.Resource)
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

			var items []*model.ResourceOutput

			ids := dm.IDs()

			switch dm.Type {
			case modelchange.EventTypeCreate, modelchange.EventTypeUpdate:
				entities, err := query.Clone().
					Where(resource.IDIn(ids...)).
					// Must append environment ID.
					Select(resource.FieldEnvironmentID).
					// Must extract template.
					Select(resource.FieldTemplateID).
					WithTemplate(func(tvq *model.TemplateVersionQuery) {
						tvq.Select(
							templateversion.FieldID,
							templateversion.FieldTemplateID,
							templateversion.FieldName,
							templateversion.FieldVersion,
							templateversion.FieldProjectID)
						if req.WithSchema {
							tvq.Select(templateversion.FieldSchema)
							tvq.Select(templateversion.FieldUISchema)
						}
					}).
					Unique(false).
					All(stream)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeResources(entities)
			case modelchange.EventTypeDelete:
				items = make([]*model.ResourceOutput, len(ids))
				for i := range ids {
					items[i] = &model.ResourceOutput{
						ID:   ids[i],
						Name: dm.Data[i].Name,
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

	if orders, ok := req.Sorting(sortFields, model.Desc(resource.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Must append environment ID.
		Select(resource.FieldEnvironmentID).
		// Must extract template.
		Select(resource.FieldTemplateID).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldID,
				templateversion.FieldTemplateID,
				templateversion.FieldName,
				templateversion.FieldVersion,
				templateversion.FieldProjectID)
			if req.WithSchema {
				tvq.Select(
					templateversion.FieldSchema,
					templateversion.FieldUISchema,
				)
			}
		}).
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeResources(entities), cnt, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	resources, err := h.modelClient.Resources().Query().
		Where(resource.IDIn(req.IDs()...)).
		All(req.Context)
	if err != nil {
		return err
	}

	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return err
	}

	err = h.cleanResourcePlanFiles(req.Context, h.modelClient, req.IDs()...)
	if err != nil {
		return err
	}

	return pkgresource.CollectionDelete(req.Context, h.modelClient, resources, pkgresource.DeleteOptions{
		Options: pkgresource.Options{
			Deployer: dp,
			Preview:  req.Preview,
		},
		WithoutCleanup: req.WithoutCleanup,
	})
}

func (h Handler) cleanResourcePlanFiles(ctx context.Context, mc model.ClientSet, resourceIDs ...object.ID) error {
	runIDs, err := h.modelClient.ResourceRuns().Query().
		Where(resourcerun.ResourceIDIn(resourceIDs...)).
		IDs(ctx)
	if err != nil {
		return err
	}

	return pkgrun.CleanPlanFiles(ctx, mc, h.storageManager, runIDs...)
}
