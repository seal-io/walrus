package resource

import (
	"context"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/pkg/deployer"
	deployertf "github.com/seal-io/walrus/pkg/deployer/terraform"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	if req.Draft {
		_, err := pkgresource.CreateDraftResources(req.Context, req.Client, entity)
		return model.ExposeResource(entity), err
	}

	dp, err := h.getDeployer(req.Context)
	if err != nil {
		return nil, err
	}

	if err = pkgresource.SetSubjectID(req.Context, entity); err != nil {
		return nil, err
	}

	createOpts := pkgresource.Options{
		Deployer: dp,
	}

	return pkgresource.Create(
		req.Context,
		h.modelClient,
		entity,
		createOpts,
	)
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
	if req.WithoutCleanup {
		// Do not clean deployed native resources.
		return h.modelClient.Resources().DeleteOneID(req.ID).
			Exec(req.Context)
	}

	dp, err := h.getDeployer(req.Context)
	if err != nil {
		return err
	}

	destroyOpts := pkgresource.Options{
		Deployer: dp,
	}

	return pkgresource.Destroy(
		req.Context,
		h.modelClient,
		req.Model(),
		destroyOpts)
}

func (h Handler) CollectionCreate(req CollectionCreateRequest) (CollectionCreateResponse, error) {
	entities := req.Model()

	dp, err := h.getDeployer(req.Context)
	if err != nil {
		return nil, err
	}

	err = h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		if err := pkgresource.SetSubjectID(req.Context, entities...); err != nil {
			return err
		}

		if req.Draft {
			_, err := pkgresource.CreateDraftResources(req.Context, tx, entities...)
			return err
		}

		_, err := pkgresource.CreateScheduledResources(req.Context, tx, dp, entities)

		return err
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to create resources")
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
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.Resources().Query().
		Where(resource.EnvironmentID(req.Environment.ID))

	// At the moment, a resource is considered a service if it directly uses templates.
	// In the long run, it may be extended to both template/definition resources if the category label is present.
	if req.IsService != nil && *req.IsService {
		query.Where(resource.TemplateIDNotNil())
	} else if req.IsService != nil && !*req.IsService {
		query.Where(resource.ResourceDefinitionIDNotNil())
	}

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
							tvq.Select(templateversion.FieldUiSchema)
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
						ID: ids[i],
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
					templateversion.FieldUiSchema,
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
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		if req.WithoutCleanup {
			// Do not clean deployed native resources.
			_, err := tx.Resources().Delete().
				Where(resource.IDIn(ids...)).
				Exec(req.Context)

			return err
		}

		resources, err := tx.Resources().Query().
			Where(resource.IDIn(ids...)).
			All(req.Context)
		if err != nil {
			return err
		}

		deployerOpts := deptypes.CreateOptions{
			Type:       deployertf.DeployerType,
			KubeConfig: h.kubeConfig,
		}

		dp, err := deployer.Get(req.Context, deployerOpts)
		if err != nil {
			return err
		}

		destroyOpts := pkgresource.Options{
			Deployer: dp,
		}

		for _, s := range resources {
			err = pkgresource.Destroy(req.Context, tx, s, destroyOpts)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (h Handler) getDeployer(ctx context.Context) (deptypes.Deployer, error) {
	dep, err := deployer.Get(ctx, deptypes.CreateOptions{
		Type:       deployertf.DeployerType,
		KubeConfig: h.kubeConfig,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to get deployer")
	}

	return dep, nil
}
