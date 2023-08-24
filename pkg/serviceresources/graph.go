package serviceresources

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
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

		// Hide service resource's dependencies.
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

func GetGraphVertexType(m *model.ServiceResource) string {
	if m.Shape == types.ServiceResourceShapeClass {
		return types.VertexKindServiceResourceGroup
	}

	return types.VertexKindServiceResource
}
