package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Resource is the schema for the resources API.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:crd-gen:resource:scope="Namespaced",subResources=["status"]
type Resource struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceSpec   `json:"spec,omitempty"`
	Status ResourceStatus `json:"status,omitempty"`
}

var _ runtime.Object = (*Resource)(nil)

// ResourceSpec defines the desired state of Resource.
type ResourceSpec struct {
	// TODO: your spec here
}

// ResourceStatus defines the observed state of Resource.
type ResourceStatus struct {
	// TODO: your status here
}

// ResourceList holds the list of Resource.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []Resource `json:"items"`
}

var _ runtime.Object = (*ResourceList)(nil)
