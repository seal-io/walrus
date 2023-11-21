package resource

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/pkg/operator/k8s/intercept"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	pkgcomponent "github.com/seal-io/walrus/pkg/resourcecomponents"
	tfparser "github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/topic"
	"github.com/seal-io/walrus/utils/validation"
)

func (h Handler) RouteUpgrade(req RouteUpgradeRequest) error {
	entity := req.Model()

	// Update service, mark status from deploying.
	status.ResourceStatusDeployed.Reset(entity, "Upgrading")
	entity.Status.SetSummary(status.WalkResource(&entity.Status))

	err := h.modelClient.WithTx(req.Context, func(tx *model.Tx) (err error) {
		entity, err = tx.Resources().UpdateOne(entity).
			Set(entity).
			SaveE(req.Context, dao.ResourceDependenciesEdgeSave)

		return err
	})
	if err != nil {
		return errorx.Wrap(err, "error updating service")
	}

	dp, err := h.getDeployer(req.Context)
	if err != nil {
		return err
	}

	// Apply resource.
	applyOpts := pkgresource.Options{
		Deployer: dp,
	}

	ready, err := pkgresource.CheckDependencyStatus(req.Context, h.modelClient, entity)
	if err != nil {
		return errorx.Wrap(err, "error checking dependency status")
	}

	if ready {
		return pkgresource.Apply(
			req.Context,
			h.modelClient,
			entity,
			applyOpts)
	}

	return nil
}

func (h Handler) RouteRollback(req RouteRollbackRequest) error {
	rev, err := h.modelClient.ResourceRevisions().Query().
		Where(
			resourcerevision.ID(req.RevisionID),
			resourcerevision.ResourceID(req.ID)).
		WithResource().
		Only(req.Context)
	if err != nil {
		return err
	}

	entity := rev.Edges.Resource
	entity.Attributes = rev.Attributes

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

	entity, err = h.modelClient.Resources().UpdateOne(entity).
		Set(entity).
		SaveE(req.Context, dao.ResourceDependenciesEdgeSave)
	if err != nil {
		return errorx.Wrap(err, "error updating service")
	}

	dp, err := h.getDeployer(req.Context)
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

	status.ResourceStatusUnDeployed.Remove(entity)
	status.ResourceStatusStopped.Remove(entity)
	status.ResourceStatusDeployed.Reset(entity, "Deploying")
	entity.Status.SetSummary(status.WalkResource(&entity.Status))

	err := h.modelClient.WithTx(req.Context, func(tx *model.Tx) (err error) {
		entity, err = tx.Resources().UpdateOne(entity).
			Set(entity).
			SaveE(req.Context, dao.ResourceDependenciesEdgeSave)

		return err
	})
	if err != nil {
		return errorx.Wrap(err, "error updating resource")
	}

	dp, err := h.getDeployer(req.Context)
	if err != nil {
		return err
	}

	applyOpts := pkgresource.Options{
		Deployer: dp,
	}

	ready, err := pkgresource.CheckDependencyStatus(req.Context, h.modelClient, entity)
	if err != nil {
		return errorx.Wrap(err, "error checking dependency status")
	}

	if ready {
		return pkgresource.Apply(
			req.Context,
			h.modelClient,
			entity,
			applyOpts)
	}

	return nil
}

func (h Handler) RouteStop(req RouteStopRequest) error {
	dp, err := h.getDeployer(req.Context)
	if err != nil {
		return err
	}

	opts := pkgresource.Options{
		Deployer: dp,
	}

	return pkgresource.Stop(req.Context, req.Client, req.Model(), opts)
}

func (h Handler) RouteGetAccessEndpoints(req RouteGetAccessEndpointsRequest) (RouteGetAccessEndpointsResponse, error) {
	if stream := req.Stream; stream != nil {
		t, err := topic.Subscribe(modelchange.ResourceRevision)
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

			for _, id := range dm.IDs {
				ar, err := h.modelClient.ResourceRevisions().Query().
					Where(resourcerevision.ID(id)).
					Only(stream)
				if err != nil {
					return nil, err
				}

				if ar.ResourceID != req.ID {
					continue
				}

				eps, err := h.getAccessEndpoints(stream, req.ID)
				if err != nil {
					return nil, err
				}

				if len(eps) == 0 {
					continue
				}

				var resp *runtime.ResponseCollection

				switch dm.Type {
				case modelchange.EventTypeCreate:
					// While create new resource revision,
					// the previous endpoints from outputs and resources need to be deleted.
					resp = runtime.TypedResponse(modelchange.EventTypeDelete.String(), eps)
				case modelchange.EventTypeUpdate:
					// While the resource revision status is succeeded,
					// the endpoints is updated to the current revision.
					if status.ResourceRevisionStatusReady.IsTrue(ar) {
						continue
					}

					resp = runtime.TypedResponse(modelchange.EventTypeUpdate.String(), eps)
				}

				if err = stream.SendJSON(resp); err != nil {
					return nil, err
				}
			}
		}
	}

	return h.getAccessEndpoints(req.Context, req.ID)
}

func (h Handler) getAccessEndpoints(ctx context.Context, id object.ID) ([]AccessEndpoint, error) {
	// Endpoints from output.
	eps, err := h.getEndpointsFromOutput(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(eps) != 0 {
		return eps, nil
	}

	// Endpoints from resources.
	return h.getEndpointsFromResources(ctx, id)
}

func (h Handler) getEndpointsFromOutput(ctx context.Context, id object.ID) ([]AccessEndpoint, error) {
	outputs, err := h.getResourceOutputs(ctx, id, true)
	if err != nil {
		return nil, err
	}

	var (
		invalidTypeErr = errorx.NewHttpError(http.StatusBadRequest,
			"element type of output endpoints should be string")
		endpoints = make([]AccessEndpoint, 0, len(outputs))
	)

	for _, v := range outputs {
		if !strings.HasPrefix(v.Name, "endpoint") {
			continue
		}

		switch {
		case v.Schema.Type == openapi3.TypeString:
			ep, _, err := property.GetString(v.Value)
			if err != nil {
				return nil, err
			}

			if err = validation.IsValidEndpoint(ep); err != nil {
				return nil, errorx.NewHttpError(http.StatusBadRequest, err.Error())
			}

			endpoints = append(endpoints, AccessEndpoint{
				Endpoints: []string{ep},
				Name:      v.Name,
			})
		case v.Schema.Type == openapi3.TypeArray:
			if v.Schema.Items != nil && v.Schema.Items.Value != nil {
				if v.Schema.Items.Value.Type != openapi3.TypeObject &&
					v.Schema.Items.Value.Type != openapi3.TypeString {
					return nil, invalidTypeErr
				}
			}

			eps, _, err := property.GetSlice[string](v.Value)
			if err != nil {
				return nil, err
			}

			if err := validation.IsValidEndpoints(eps); err != nil {
				return nil, err
			}

			if len(eps) == 0 {
				continue
			}

			endpoints = append(endpoints, AccessEndpoint{
				Endpoints: eps,
				Name:      v.Name,
			})
		}
	}

	return endpoints, nil
}

func (h Handler) getEndpointsFromResources(ctx context.Context, id object.ID) ([]AccessEndpoint, error) {
	res, err := h.modelClient.ResourceComponents().Query().
		Where(
			resourcecomponent.ResourceID(id),
			resourcecomponent.Mode(types.ResourceComponentModeManaged),
			resourcecomponent.TypeIn(intercept.TFEndpointsTypes...)).
		Select(
			resourcecomponent.FieldConnectorID,
			resourcecomponent.FieldType,
			resourcecomponent.FieldName,
			resourcecomponent.FieldStatus,
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	var endpoints []AccessEndpoint

	for _, v := range res {
		for _, eps := range v.Status.ResourceEndpoints {
			endpoints = append(endpoints, AccessEndpoint{
				Name:      filepath.Join(eps.EndpointType, filepath.Base(v.Name)),
				Endpoints: eps.Endpoints,
			})
		}
	}

	return endpoints, nil
}

func (h Handler) RouteGetOutputs(req RouteGetOutputsRequest) (RouteGetOutputsResponse, error) {
	if stream := req.Stream; stream != nil {
		t, err := topic.Subscribe(modelchange.ResourceRevision)
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

			for _, id := range dm.IDs {
				ar, err := h.modelClient.ResourceRevisions().Query().
					Where(resourcerevision.ID(id)).
					Only(stream)
				if err != nil {
					return nil, err
				}

				if ar.ResourceID != req.ID {
					continue
				}

				outs, err := h.getResourceOutputs(stream, ar.ResourceID, false)
				if err != nil {
					return nil, err
				}

				if len(outs) == 0 {
					continue
				}

				var resp *runtime.ResponseCollection

				switch dm.Type {
				case modelchange.EventTypeCreate:
					// While create new resource revision,
					// the outputs of new revision is the previous outputs.
					resp = runtime.TypedResponse(modelchange.EventTypeDelete.String(), outs)
				case modelchange.EventTypeUpdate:
					// While the resource revision status is succeeded,
					// the outputs is updated to the current revision.
					if status.ResourceRevisionStatusReady.IsTrue(ar) {
						continue
					}

					resp = runtime.TypedResponse(modelchange.EventTypeUpdate.String(), outs)
				}

				if err = stream.SendJSON(resp); err != nil {
					return nil, err
				}
			}
		}
	}

	return h.getResourceOutputs(req.Context, req.ID, true)
}

func (h Handler) getResourceOutputs(ctx context.Context, id object.ID, onlySuccess bool) ([]types.OutputValue, error) {
	sr, err := h.modelClient.ResourceRevisions().Query().
		Where(resourcerevision.ResourceID(id)).
		Select(
			resourcerevision.FieldOutput,
			resourcerevision.FieldTemplateName,
			resourcerevision.FieldTemplateVersion,
			resourcerevision.FieldAttributes,
			resourcerevision.FieldStatus,
		).
		WithResource(func(sq *model.ResourceQuery) {
			sq.Select(resource.FieldName)
		}).
		Order(model.Desc(resourcerevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, fmt.Errorf("error getting the latest service revision")
	}

	if sr == nil {
		return nil, nil
	}

	if onlySuccess && !status.ResourceRevisionStatusReady.IsTrue(sr) {
		return nil, nil
	}

	o, err := tfparser.ParseStateOutput(sr)
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
