package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceDefinition is the schema for the resource definitions API.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:crd-gen:resource:scope="Namespaced",subResources=["status"]
type ResourceDefinition struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceDefinitionSpec   `json:"spec,omitempty"`
	Status ResourceDefinitionStatus `json:"status,omitempty"`
}

var _ runtime.Object = (*ResourceDefinition)(nil)

// ResourceDefinitionSpec defines the desired state of ResourceDefinition.
type ResourceDefinitionSpec struct {
	// TODO: your spec here
}

// ResourceDefinitionStatus defines the observed state of ResourceDefinition.
type ResourceDefinitionStatus struct {
	// TODO: your status here
}

// ResourceDefinitionList holds the list of ResourceDefinition.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceDefinitionList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []ResourceDefinition `json:"items"`
}

var _ runtime.Object = (*ResourceDefinitionList)(nil)
