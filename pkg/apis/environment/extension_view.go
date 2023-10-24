package environment

import (
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
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

type (
	RouteCloneEnvironmentRequest struct {
		_ struct{} `route:"POST=/clone"`

		CloneEnvironmentId object.ID `path:"environment"`

		model.EnvironmentCreateInput `path:",inline" json:",inline"`
	}

	RouteCloneEnvironmentResponse = model.EnvironmentOutput
)

func (r *RouteCloneEnvironmentRequest) Validate() error {
	return validateEnvironmentCreateInput(r.EnvironmentCreateInput)
}

type (
	RouteGetResourceDefinitionsRequest struct {
		_ struct{} `route:"GET=/resource-definitions"`

		model.EnvironmentQueryInput `path:",inline"`
	}

	RouteGetResourceDefinitionsResponse = []*model.ResourceDefinitionOutput
)
