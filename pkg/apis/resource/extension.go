package resource

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	pkgcomponent "github.com/seal-io/walrus/pkg/resourcecomponents"
	pkgresource "github.com/seal-io/walrus/pkg/resources"
	tfparser "github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) RouteUpgrade(req RouteUpgradeRequest) (*RouteUpgradeResponse, error) {
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

	return &RouteUpgradeResponse{
		ResourceOutput: model.ExposeResource(entity),
		Run:            model.ExposeResourceRun(run),
	}, nil
}

func (h Handler) RouteRollback(req RouteRollbackRequest) error {
	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return err
	}

	return pkgresource.Rollback(
		req.Context,
		h.modelClient,
		req.ID,
		req.RunID,
		pkgresource.Options{
			Deployer:      dp,
			ChangeComment: req.ChangeComment,
		})
}

func (h Handler) RouteStart(req RouteStartRequest) (*RouteStartResponse, error) {
	entity := req.resource

	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return nil, err
	}

	run, err := pkgresource.Start(req.Context, h.modelClient, entity, pkgresource.Options{
		Deployer:      dp,
		ChangeComment: req.ChangeComment,
	})
	if err != nil {
		return nil, err
	}

	return &RouteStartResponse{
		ResourceOutput: model.ExposeResource(entity),
		Run:            model.ExposeResourceRun(run),
	}, nil
}

func (h Handler) RouteStop(req RouteStopRequest) error {
	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return err
	}

	opts := pkgresource.Options{
		Deployer:      dp,
		ChangeComment: req.ChangeComment,
	}

	entity, err := h.modelClient.Resources().Query().
		Where(resource.ID(req.ID)).
		Only(req.Context)
	if err != nil {
		return err
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

	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return err
	}

	return pkgresource.CollectionStart(req.Context, h.modelClient, resources, pkgresource.Options{
		Deployer:      dp,
		ChangeComment: req.ChangeComment,
	})
}

func (h Handler) CollectionRouteStop(req CollectionRouteStopRequest) error {
	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return err
	}

	opts := pkgresource.Options{
		Deployer:      dp,
		ChangeComment: req.ChangeComment,
	}

	return pkgresource.CollectionStop(req.Context, h.modelClient, req.Resources, opts)
}

func (h Handler) CollectionRouteUpgrade(req CollectionRouteUpgradeRequest) error {
	resources := req.Model()

	dp, err := getDeployer(req.Context, h.kubeConfig)
	if err != nil {
		return err
	}

	return pkgresource.CollectionUpgrade(req.Context, h.modelClient, resources, pkgresource.Options{
		Deployer:      dp,
		ChangeComment: req.ChangeComment,
		Draft:         req.Draft,
	})
}
