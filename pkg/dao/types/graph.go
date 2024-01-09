package types

import (
	"time"

	"github.com/seal-io/walrus/pkg/dao/types/status"
)

const (
	// VertexKindResource indicates the vertex kind of resource, it contains resource groups.
	VertexKindResource = "Resource"
	// VertexKindResourceComponentGroup indicates the group resource that generates same type resources.
	VertexKindResourceComponentGroup = "ResourceComponentGroup"
	// VertexKindResourceComponent indicates the vertex kind of resource component.
	VertexKindResourceComponent = "ResourceComponent"

	// EdgeTypeComposition indicates vertex is composed of another vertex.
	EdgeTypeComposition = "Composition"
	// EdgeTypeRealization indicates vertex are realized by another vertex.
	EdgeTypeRealization = "Realization"
	// EdgeTypeDependency indicates vertex has dependency on another vertex.
	EdgeTypeDependency = "Dependency"
)

// GraphVertexID defines the identifier of the vertex,
// which uniquely represents an API resource.
type GraphVertexID struct {
	// Kind indicates the kind of the resource,
	// which should be the same as the API handler's Kind returning.
	Kind string `json:"kind"`
	// ID indicates the identifier of the resource.
	ID any `json:"id,omitempty"`
}

// GraphVertex defines the vertex of graph.
type GraphVertex struct {
	GraphVertexID `json:",inline"`

	// Name indicates a human-readable string of the vertex.
	Name string `json:"name,omitempty"`
	// Description indicates the detail of the vertex.
	Description string `json:"description,omitempty"`
	// Icon indicates the icon of the vertex.
	Icon string `json:"icon,omitempty"`
	// Labels indicates the labels of the vertex.
	Labels map[string]string `json:"labels,omitempty"`
	// CreateTime indicates the time to create the vertex.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// UpdateTime indicates the time to update the vertex.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Status observes the status of the vertex.
	Status status.Summary `json:"status,omitempty"`
	// Extensions records the other information of the vertex,
	// e.g. the physical type, logical category or the operational keys of the resource.
	Extensions map[string]any `json:"extensions,omitempty"`
}

// GraphEdge defines the edge of graph.
type GraphEdge struct {
	// Type indicates the type of the edge, like: Dependency or Composition.
	Type string `json:"type"`
	// Start indicates the beginning of the edge.
	Start GraphVertexID `json:"start"`
	// End indicates the ending of the edge.
	End GraphVertexID `json:"end"`
}
