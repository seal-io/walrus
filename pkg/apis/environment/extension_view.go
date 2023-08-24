package environment

import (
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
)

type (
	// GraphVertexID defines the identifier of the vertex,
	// which uniquely represents an API resource.
	GraphVertexID = types.GraphVertexID
	// GraphVertex defines the vertex of graph.
	GraphVertex = types.GraphVertex
	// GraphEdge defines the edge of graph.
	GraphEdge = types.GraphEdge

	RouteGetGraphRequest struct {
		_ struct{} `route:"GET=/graph"`

		model.EnvironmentQueryInput `path:",inline"`

		WithoutKeys bool `query:"withoutKeys,omitempty"`
	}

	RouteGetGraphResponse struct {
		Vertices []GraphVertex `json:"vertices"`
		Edges    []GraphEdge   `json:"edges"`
	}
)
