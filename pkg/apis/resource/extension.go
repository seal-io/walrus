package resource

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	pkgcomponent "github.com/seal-io/walrus/pkg/resourcecomponents"
	tfparser "github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) RouteUpgrade(req RouteUpgradeRequest) error {
	entity := req.Model()

	return upgrade(req.Context, h.kubeConfig, h.modelClient, entity, req.Draft)
}

func (h Handler) RouteRollback(req RouteRollbackRequest) error {
	rev, err := h.modelClient.ResourceRuns().Query().
		Where(
			resourcerun.ID(req.RunID),
			resourcerun.ResourceID(req.ID)).
		WithResource().
		Only(req.Context)
	if err != nil {
		return err
	}

	entity := rev.Edges.Resource
	entity.Attributes = rev.Attributes
	entity.ComputedAttributes = rev.ComputedAttributes
	entity.ChangeComment = req.ChangeComment

	if entity.TemplateID != nil {
		// Find previous template version when the resource is using template not definition.
		tv, err := h.modelClient.TemplateVersions().Query().
			Where(
				templateversion.Version(rev.TemplateVersion),
				templateversion.TemplateID(rev.TemplateID)).
			Only(req.Context)
		if err != nil {
			return err
		}

		entity.TemplateID = &tv.ID
	}

	status.ResourceStatusDeployed.Reset(entity, "Rolling back")
	entity.Status.SetSummary(status.WalkResource(&entity.Status))

	if err := pkgresource.SetSubjectID(req.Context, entity); err != nil {
		return err
	}

	entity, err = h.modelClient.Resources().UpdateOne(entity).
		Set(entity).
		SaveE(req.Context, dao.ResourceDependenciesEdgeSave)
	if err != nil {
		return errorx.Wrap(err, "error updating resource")
	}

	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return err
	}

	applyOpts := pkgresource.Options{
		Deployer: dp,
	}

	return pkgresource.Apply(
		req.Context,
		h.modelClient,
		entity,
		applyOpts)
}

func (h Handler) RouteStart(req RouteStartRequest) error {
	entity := req.resource
	entity.ChangeComment = req.ChangeComment

	return h.start(req.Context, entity)
}

func (h Handler) start(ctx context.Context, entity *model.Resource) error {
	status.ResourceStatusUnDeployed.Remove(entity)
	status.ResourceStatusStopped.Remove(entity)
	status.ResourceStatusDeployed.Reset(entity, "Deploying")
	entity.Status.SetSummary(status.WalkResource(&entity.Status))

	if err := pkgresource.SetSubjectID(ctx, entity); err != nil {
		return err
	}

	err := h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		entity, err = tx.Resources().UpdateOne(entity).
			Set(entity).
			SaveE(ctx, dao.ResourceDependenciesEdgeSave)

		return err
	})
	if err != nil {
		return errorx.Wrap(err, "error updating resource")
	}

	dp, err := getDeployer(ctx, h.kubeConfig)
	if err != nil {
		return err
	}

	applyOpts := pkgresource.Options{
		Deployer: dp,
	}

	ready, err := pkgresource.CheckDependencyStatus(ctx, h.modelClient, dp, entity)
	if err != nil {
		return errorx.Wrap(err, "error checking dependency status")
	}

	if ready {
		return pkgresource.Apply(
			ctx,
			h.modelClient,
			entity,
			applyOpts)
	}

	return nil
}

func (h Handler) RouteStop(req RouteStopRequest) error {
	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return err
	}

	opts := pkgresource.Options{
		Deployer: dp,
	}

	entity, err := h.modelClient.Resources().Query().
		Where(resource.ID(req.ID)).
		Only(req.Context)
	if err != nil {
		return err
	}
	entity.ChangeComment = req.ChangeComment

	if err := pkgresource.SetSubjectID(req.Context, entity); err != nil {
		return err
	}

	entity, err = h.modelClient.Resources().UpdateOne(entity).
		Set(entity).
		Save(req.Context)
	if err != nil {
		return errorx.Wrap(err, "error updating resource")
	}

	return pkgresource.Stop(req.Context, h.modelClient, entity, opts)
}

func (h Handler) RouteGetEndpoints(req RouteGetEndpointsRequest) (RouteGetEndpointsResponse, error) {
	if stream := req.Stream; stream != nil {
		t, err := topic.Subscribe(modelchange.Resource)
		if err != nil {
			return nil, err
		}

		defer func() { t.Unsubscribe() }()

		for {
			var event topic.Event

			event, err = t.Receive(stream)
			if err != nil {
				return nil, err
			}

			dm, ok := event.Data.(modelchange.Event)
			if !ok {
				continue
			}

			if dm.Type == modelchange.EventTypeDelete {
				continue
			}

			for _, id := range dm.IDs() {
				if id != req.ID {
					continue
				}

				entity, err := h.modelClient.Resources().Query().
					Where(resource.ID(id)).
					Select(resource.FieldEndpoints).
					Only(stream)
				if err != nil {
					return nil, err
				}

				resp := runtime.TypedResponse(modelchange.EventTypeUpdate.String(), entity.Endpoints)
				if err = stream.SendJSON(resp); err != nil {
					return nil, err
				}

				break // NB(thxCode): reduce duplicated sending in the same event.
			}
		}
	}

	entity, err := h.modelClient.Resources().Query().
		Where(resource.ID(req.ID)).
		Select(resource.FieldEndpoints).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	return entity.Endpoints, nil
}

func (h Handler) RouteGetOutputs(req RouteGetOutputsRequest) (RouteGetOutputsResponse, error) {
	query := h.modelClient.Resources().Query().
		Where(resource.ID(req.ID))

	if stream := req.Stream; stream != nil {
		t, err := topic.Subscribe(modelchange.Resource)
		if err != nil {
			return nil, err
		}

		defer func() { t.Unsubscribe() }()

		for {
			var event topic.Event

			event, err = t.Receive(stream)
			if err != nil {
				return nil, err
			}

			dm, ok := event.Data.(modelchange.Event)
			if !ok {
				continue
			}

			if dm.Type == modelchange.EventTypeDelete {
				continue
			}

			for _, id := range dm.IDs() {
				if id != req.ID {
					continue
				}

				res, err := query.Clone().
					WithState().
					Only(stream)
				if err != nil {
					return nil, err
				}

				outs, err := h.getResourceOutputs(res)
				if err != nil {
					return nil, err
				}

				if len(outs) == 0 {
					continue
				}

				resp := runtime.TypedResponse(dm.Type.String(), outs)

				if err = stream.SendJSON(resp); err != nil {
					return nil, err
				}
			}
		}
	}

	res, err := query.Clone().
		WithState().
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	return h.getResourceOutputs(res)
}

func (h Handler) getResourceOutputs(resource *model.Resource) ([]types.OutputValue, error) {
	var p tfparser.StateParser

	o, err := p.GetOriginalOutputs(resource.Edges.State.Data, resource.Name)
	if err != nil {
		return nil, fmt.Errorf("error get outputs: %w", err)
	}

	return o, nil
}

func (h Handler) RouteGetGraph(req RouteGetGraphRequest) (*RouteGetGraphResponse, error) {
	entities, err := dao.ResourceComponentShapeClassQuery(h.modelClient.ResourceComponents().Query()).
		Where(resourcecomponent.ResourceID(req.ID)).
		All(req.Context)
	if err != nil {
		return nil, err
	}

	// Calculate capacity for allocation.
	var verticesCap, edgesCap int
	{
		// Count the number of ResourceComponent.
		verticesCap = len(entities)
		for i := 0; i < len(entities); i++ {
			// Count the vertex size of sub ResourceComponent,
			// and the edge size from ResourceComponent to sub ResourceComponent.
			verticesCap += len(entities[i].Edges.Components)
			edgesCap += len(entities[i].Edges.Components)
		}
	}

	// Construct response.
	var (
		vertices = make([]GraphVertex, 0, verticesCap)
		edges    = make([]GraphEdge, 0, edgesCap)
	)

	// Set keys for next operations, e.g. Log, Exec and so on.
	if !req.WithoutKeys {
		pkgcomponent.SetKeys(
			req.Context,
			log.WithName("api").WithName("resource"),
			h.modelClient,
			entities,
			nil)
	}

	vertices, edges = pkgcomponent.GetVerticesAndEdges(
		entities, vertices, edges)

	return &RouteGetGraphResponse{
		Vertices: vertices,
		Edges:    edges,
	}, nil
}

func (h Handler) CollectionRouteStart(req CollectionRouteStartRequest) error {
	// Start resources in topological order.
	resources, err := pkgresource.TopologicalSortResources(req.Resources)
	if err != nil {
		return err
	}

	for i := range resources {
		entity := resources[i]
		entity.ChangeComment = req.ChangeComment

		if err := h.start(req.Context, entity); err != nil {
			return err
		}
	}

	return nil
}

func (h Handler) CollectionRouteStop(req CollectionRouteStopRequest) error {
	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return err
	}

	opts := pkgresource.Options{
		Deployer: dp,
	}

	resources, err := pkgresource.ReverseTopologicalSortResources(req.Resources)
	if err != nil {
		return err
	}

	for i := range resources {
		res := resources[i]
		res.ChangeComment = req.ChangeComment

		if err := pkgresource.SetSubjectID(req.Context, res); err != nil {
			return err
		}

		res, err := h.modelClient.Resources().UpdateOne(res).
			Set(res).
			Save(req.Context)
		if err != nil {
			return errorx.Wrap(err, "error updating resource")
		}

		if err := pkgresource.Stop(req.Context, h.modelClient, res, opts); err != nil {
			return err
		}
	}

	return nil
}

func (h Handler) CollectionRouteUpgrade(req CollectionRouteUpgradeRequest) error {
	return UpgradeResources(req, h.modelClient, h.kubeConfig)
}
