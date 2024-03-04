package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Connector is the schema for the connectors API.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:crd-gen:resource:scope="Namespaced",subResources=["status"]
type Connector struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConnectorSpec   `json:"spec,omitempty"`
	Status ConnectorStatus `json:"status,omitempty"`
}

var _ runtime.Object = (*Connector)(nil)

// ConnectorSpec defines the desired state of Connector.
type ConnectorSpec struct {
	// TODO: your spec here
}

// ConnectorStatus defines the observed state of Connector.
type ConnectorStatus struct {
	// TODO: your status here
}

// ConnectorList holds the list of Connector.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ConnectorList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []Connector `json:"items"`
}

var _ runtime.Object = (*ConnectorList)(nil)
