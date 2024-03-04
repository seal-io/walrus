package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Catalog is the schema for the catalogs API.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:crd-gen:resource:scope="Namespaced",subResources=["status"]
type Catalog struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   CatalogSpec   `json:"spec,omitempty"`
	Status CatalogStatus `json:"status,omitempty"`
}

var _ runtime.Object = (*Catalog)(nil)

// CatalogSpec defines the desired state of Catalog.
type CatalogSpec struct {
	// TODO: your spec here
}

// CatalogStatus defines the observed state of Catalog.
type CatalogStatus struct {
	// TODO: your status here
}

// CatalogList holds the list of Catalog.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CatalogList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []Catalog `json:"items"`
}

var _ runtime.Object = (*CatalogList)(nil)
