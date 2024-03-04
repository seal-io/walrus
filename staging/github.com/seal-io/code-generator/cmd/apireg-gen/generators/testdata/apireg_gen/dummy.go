package apireg_gen

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Dummy is the schema for the projects API.
//
// +k8s:apireg-gen:resource:categories=["all","walrus"],scope="Cluster",shortName=["proj"],plural="projects",subResources=["status","scale","others"]
type Dummy struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   DummySpec   `json:"spec,omitempty"`
	Status DummyStatus `json:"status,omitempty"`
}

// DummySpec defines the desired state of Dummy.
type DummySpec struct {
}

// DummyStatus defines the observed state of Dummy.
type DummyStatus struct {
}

var _ runtime.Object = (*DummyList)(nil)

// DummyList holds the list of Dummy.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DummyList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []Dummy `json:"items"`
}
