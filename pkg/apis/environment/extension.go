package environment

import (
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/dao/model/variable"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/utils/errorx"

	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/service"
	"github.com/seal-io/walrus/pkg/dao/model/servicerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	pkgresource "github.com/seal-io/walrus/pkg/serviceresources"
	"github.com/seal-io/walrus/utils/log"
)

var getServiceFields = service.WithoutFields(
	service.FieldUpdateTime)

func (h Handler) RouteGetGraph(req RouteGetGraphRequest) (*RouteGetGraphResponse, error) {
	// Fetch service entities.
	entities, err := h.modelClient.Services().Query().
		Where(service.EnvironmentID(req.ID)).
		Order(model.Desc(service.FieldCreateTime)).
		Select(getServiceFields...).
		// Must extract dependency.
		WithDependencies(func(dq *model.ServiceRelationshipQuery) {
			dq.Select(servicerelationship.FieldDependencyID).
				Where(func(s *sql.Selector) {
					s.Where(sql.ColumnsNEQ(servicerelationship.FieldServiceID, servicerelationship.FieldDependencyID))
				})
		}).
		// Must extract resource.
		WithResources(func(rq *model.ServiceResourceQuery) {
			dao.ServiceResourceShapeClassQuery(rq)
		}).
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, err
	}

	verticesCap, edgesCap := getCaps(entities)

	// Construct response.
	var (
		vertices  = make([]GraphVertex, 0, verticesCap)
		edges     = make([]GraphEdge, 0, edgesCap)
		operators = make(map[object.ID]optypes.Operator)
	)

	for i := 0; i < len(entities); i++ {
		entity := entities[i]

		// Append Service to vertices.
		vertices = append(vertices, GraphVertex{
			GraphVertexID: GraphVertexID{
				Kind: types.VertexKindService,
				ID:   entity.ID,
			},
			Name:        entity.Name,
			Description: entity.Description,
			Labels:      entity.Labels,
			CreateTime:  entity.CreateTime,
			UpdateTime:  entity.UpdateTime,
			Status:      entity.Status.Summary,
			Extensions: map[string]any{
				"projectID":     entity.ProjectID,
				"environmentID": entity.EnvironmentID,
			},
		})

		// Append the link of related Services to edges.
		for j := 0; j < len(entity.Edges.Dependencies); j++ {
			edges = append(edges, GraphEdge{
				Type: types.EdgeTypeDependency,
				Start: GraphVertexID{
					Kind: types.VertexKindService,
					ID:   entity.ID,
				},
				End: GraphVertexID{
					Kind: types.VertexKindService,
					ID:   entity.Edges.Dependencies[j].DependencyID,
				},
			})
		}

		// Set keys for next operations, e.g. Log, Exec and so on.
		if !req.WithoutKeys {
			pkgresource.SetKeys(
				req.Context,
				log.WithName("api").WithName("environment"),
				h.modelClient,
				entity.Edges.Resources,
				operators)
		}

		// Append ServiceResource to vertices,
		// and append the link of related ServiceResources to edges.
		vertices, edges = pkgresource.GetVerticesAndEdges(
			entity.Edges.Resources, vertices, edges)

		for j := 0; j < len(entity.Edges.Resources); j++ {
			// Append the link from Service to ServiceResource into edges.
			edges = append(edges, GraphEdge{
				Type: types.EdgeTypeComposition,
				Start: GraphVertexID{
					Kind: types.VertexKindService,
					ID:   entity.ID,
				},
				End: GraphVertexID{
					Kind: types.VertexKindServiceResourceGroup,
					ID:   entity.Edges.Resources[j].ID,
				},
			})
		}
	}

	return &RouteGetGraphResponse{
		Vertices: vertices,
		Edges:    edges,
	}, nil
}

func getCaps(entities model.Services) (int, int) {
	// Calculate capacity for allocation.
	var verticesCap, edgesCap int

	// Count the number of Service.
	verticesCap = len(entities)
	for i := 0; i < len(entities); i++ {
		// Count the vertex size of ServiceResource,
		// and the edge size from Service to ServiceResource.
		verticesCap += len(entities[i].Edges.Resources)
		edgesCap += len(entities[i].Edges.Dependencies)

		for j := 0; j < len(entities[i].Edges.Resources); j++ {
			// Count the vertex size of instances,
			// and the edge size from ServiceResourceGroup to instance ServiceResource.
			verticesCap += len(entities[i].Edges.Resources[j].Edges.Instances)
			edgesCap += len(entities[i].Edges.Resources[j].Edges.Instances) +
				len(entities[i].Edges.Resources[j].Edges.Dependencies)

			for k := 0; k < len(entities[i].Edges.Resources[j].Edges.Instances); k++ {
				verticesCap += len(entities[i].Edges.Resources[j].Edges.Components)
				edgesCap += len(entities[i].Edges.Resources[j].Edges.Components)
			}
		}
	}

	return verticesCap, edgesCap
}

func (h Handler) RouteCloneEnvironment(req RouteCloneEnvironmentRequest) (*RouteCloneEnvironmentResponse, error) {
	entity := req.Model()

	var variableNames []string

	variables := entity.Edges.Variables
	for i := range variables {
		v := variables[i]
		if v == nil {
			return nil, errorx.New("invalid input: nil variable")
		}

		if v.Sensitive && v.Value == "" {
			variableNames = append(variableNames, v.Name)
		}
	}

	// Fetch and fill the value with sensitive variables from the cloned environment.
	if len(variableNames) > 0 {
		vs, err := h.modelClient.Variables().Query().
			Where(
				variable.EnvironmentID(req.CloneEnvironmentId),
				variable.NameIn(variableNames...),
			).
			All(req.Context)
		if err != nil {
			return nil, err
		}

		variableValueMap := make(map[string]crypto.String)
		for _, vv := range vs {
			variableValueMap[vv.Name] = vv.Value
		}

		for _, v := range variables {
			if v.Sensitive && v.Value == "" {
				v.Value = variableValueMap[v.Name]
			}
		}
	}

	return createEnvironment(req.Context, h.modelClient, entity)
}
