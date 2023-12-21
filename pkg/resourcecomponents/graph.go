package resourcecomponents

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// GetVerticesAndEdges constructs a graph of the given model.ResourceComponent entities with DFS algorithm,
// and appends the vertices and edges to the given slices.
func GetVerticesAndEdges(
	entities model.ResourceComponents,
	vertices []types.GraphVertex,
	edges []types.GraphEdge,
) ([]types.GraphVertex, []types.GraphEdge) {
	var (
		visited = sets.New[object.ID]()
		dfs     func(entity *model.ResourceComponent)
	)

	dfs = func(entity *model.ResourceComponent) {
		if visited.Has(entity.ID) {
			return
		}

		visited.Insert(entity.ID)
		kind := GetGraphVertexType(entity)

		// Append ResourceComponent to vertices.
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
				"type":           entity.Type,
				"keys":           entity.Keys,
				"projectID":      entity.ProjectID,
				"environmentID":  entity.EnvironmentID,
				"resourceID":     entity.ResourceID,
				"connectorID":    entity.ConnectorID,
				"driftDetection": entity.DriftDetection,
			},
		})

		for i := 0; i < len(entity.Edges.Components); i++ {
			// Append Composition to edges.
			edges = append(edges, types.GraphEdge{
				Type: types.EdgeTypeComposition,
				Start: types.GraphVertexID{
					Kind: types.VertexKindResourceComponent,
					ID:   entity.ID,
				},
				End: types.GraphVertexID{
					Kind: types.VertexKindResourceComponent,
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
					Kind: types.VertexKindResourceComponentGroup,
					ID:   entity.ID,
				},
				End: types.GraphVertexID{
					Kind: types.VertexKindResourceComponent,
					ID:   entity.Edges.Instances[j].ID,
				},
			})

			dfs(entity.Edges.Instances[j])
		}

		// Hide resource component's dependencies.
		if entity.Shape == types.ResourceComponentShapeInstance {
			return
		}

		for k := 0; k < len(entity.Edges.Dependencies); k++ {
			// Append the edge of class resource to dependencies.
			edges = append(edges, types.GraphEdge{
				Type: entity.Edges.Dependencies[k].Type,
				Start: types.GraphVertexID{
					Kind: types.VertexKindResourceComponentGroup,
					ID:   entity.Edges.Dependencies[k].ResourceComponentID,
				},
				End: types.GraphVertexID{
					Kind: types.VertexKindResourceComponentGroup,
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

func GetGraphVertexType(m *model.ResourceComponent) string {
	if m.Shape == types.ResourceComponentShapeClass {
		return types.VertexKindResourceComponentGroup
	}

	return types.VertexKindResourceComponent
}
