package serviceresources

import (
	"context"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/operator"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/utils/log"
)

// GetVerticesAndEdges constructs a graph of the given model.ServiceResource entities with DFS algorithm,
// and appends the vertices and edges to the given slices.
func GetVerticesAndEdges(
	entities model.ServiceResources,
	vertices []types.GraphVertex,
	edges []types.GraphEdge,
) ([]types.GraphVertex, []types.GraphEdge) {
	var (
		visited = sets.New[object.ID]()
		dfs     func(entity *model.ServiceResource)
	)

	dfs = func(entity *model.ServiceResource) {
		if visited.Has(entity.ID) {
			return
		}

		visited.Insert(entity.ID)
		kind := GetGraphVertexType(entity)

		// Append ServiceResource to vertices.
		vertices = append(vertices, types.GraphVertex{
			GraphVertexID: types.GraphVertexID{
				Kind: kind,
				ID:   entity.ID,
			},
			Name:       entity.Name,
			CreateTime: entity.CreateTime,
			UpdateTime: entity.UpdateTime,
			Status:     entity.Status.Summary,
			Extensions: map[string]any{
				"type":          entity.Type,
				"keys":          entity.Keys,
				"projectID":     entity.ProjectID,
				"environmentID": entity.EnvironmentID,
				"serviceID":     entity.ServiceID,
				"connectorID":   entity.ConnectorID,
			},
		})

		for i := 0; i < len(entity.Edges.Components); i++ {
			// Append Composition to edges.
			edges = append(edges, types.GraphEdge{
				Type: types.EdgeTypeComposition,
				Start: types.GraphVertexID{
					Kind: types.VertexKindServiceResource,
					ID:   entity.ID,
				},
				End: types.GraphVertexID{
					Kind: types.VertexKindServiceResource,
					ID:   entity.Edges.Components[i].ID,
				},
			})

			dfs(entity.Edges.Components[i])
		}

		for j := 0; j < len(entity.Edges.Instances); j++ {
			// Append Realization to edges.
			edges = append(edges, types.GraphEdge{
				Type: types.EdgeTypeRealization,
				Start: types.GraphVertexID{
					Kind: types.VertexKindServiceResourceGroup,
					ID:   entity.ID,
				},
				End: types.GraphVertexID{
					Kind: types.VertexKindServiceResource,
					ID:   entity.Edges.Instances[j].ID,
				},
			})

			dfs(entity.Edges.Instances[j])
		}

		// Hide instance resources's dependencies.
		if entity.Shape == types.ServiceResourceShapeInstance {
			return
		}

		for k := 0; k < len(entity.Edges.Dependencies); k++ {
			// Append the edge of class resource to dependencies.
			edges = append(edges, types.GraphEdge{
				Type: entity.Edges.Dependencies[k].Type,
				Start: types.GraphVertexID{
					Kind: types.VertexKindServiceResourceGroup,
					ID:   entity.Edges.Dependencies[k].ServiceResourceID,
				},
				End: types.GraphVertexID{
					Kind: types.VertexKindServiceResourceGroup,
					ID:   entity.Edges.Dependencies[k].DependencyID,
				},
			})
		}
	}

	for i := 0; i < len(entities); i++ {
		dfs(entities[i])
	}

	return vertices, edges
}

// SetKeys sets the keys of the resources for operations like log and exec.
func SetKeys(
	ctx context.Context,
	entities model.ServiceResources,
	operators map[object.ID]optypes.Operator,
) model.ServiceResources {
	logger := log.WithName("service-resource")
	cache := make(map[object.ID]*model.ServiceResource)

	if operators == nil {
		operators = make(map[object.ID]optypes.Operator)
	}

	// DFS function to get resource keys.
	var fn func(entity *model.ServiceResource)
	fn = func(entity *model.ServiceResource) {
		if _, ok := cache[entity.ID]; ok {
			return
		}

		cache[entity.ID] = entity

		if IsOperable(entity) && entity.Edges.Connector != nil {
			var err error

			op, ok := operators[entity.Edges.Connector.ID]
			if !ok {
				op, err = operator.Get(ctx, optypes.CreateOptions{Connector: *entity.Edges.Connector})
				if err != nil {
					logger.Warnf("cannot get operator of connector: %v", err)
					return
				}
				operators[entity.Edges.Connector.ID] = op
			}

			entity.Keys, err = op.GetKeys(ctx, entity)
			if err != nil {
				logger.Errorf("error getting keys for %q: %v", entity.ID, err)
				return
			}
		}

		for i := 0; i < len(entity.Edges.Components); i++ {
			fn(entity.Edges.Components[i])
		}

		for j := 0; j < len(entity.Edges.Instances); j++ {
			fn(entity.Edges.Instances[j])
		}
	}

	if operators == nil {
		operators = make(map[object.ID]optypes.Operator)
	}

	for i := 0; i < len(entities); i++ {
		fn(entities[i])
	}

	return entities
}

func GetGraphVertexType(m *model.ServiceResource) string {
	if m.Shape == types.ServiceResourceShapeClass {
		return types.VertexKindServiceResourceGroup
	}

	return types.VertexKindServiceResource
}

func IsOperable(m *model.ServiceResource) bool {
	return m.Shape == types.ServiceResourceShapeInstance &&
		(m.Mode == types.ServiceResourceModeManaged || m.Mode == types.ServiceResourceModeDiscovered)
}
