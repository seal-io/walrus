package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceRun is the schema for the resource runs API.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:crd-gen:resource:scope="Namespaced",subResources=["status"]
type ResourceRun struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceDefinitionSpec   `json:"spec,omitempty"`
	Status ResourceDefinitionStatus `json:"status,omitempty"`
}

var _ runtime.Object = (*ResourceRun)(nil)

// ResourceRunSpec defines the desired state of ResourceRun.
type ResourceRunSpec struct {
	// TODO: your spec here
}

// ResourceRunStatus defines the observed state of ResourceRun.
type ResourceRunStatus struct {
	// TODO: your status here
}

// ResourceRunList holds the list of ResourceRun.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceRunList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []ResourceRun `json:"items"`
}

var _ runtime.Object = (*ResourceRunList)(nil)
