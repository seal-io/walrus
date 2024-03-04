package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceComponents holds the list of ResourceComponent.
//
// ResourceComponents is the subresource of Resource.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"],shortName=["rescomp"]
type ResourceComponents struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	// +patchStrategy=merge
	// +patchMergeKey=name
	// +listType=map
	// +listMapKey=name
	Items []ResourceComponent `json:"items,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

var _ runtime.Object = (*ResourceComponents)(nil)

// ResourceComponent is the schema for the resource component API.
type ResourceComponent struct {
	/* NB(thxCode): All attributes must be comparable. */

	// Name is the name of the component.
	Name string `json:"name"`
}
