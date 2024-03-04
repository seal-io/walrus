package v1

import (
	"k8s.io/apimachinery/pkg/runtime"

	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
)

// Catalog is the schema for the catalogs API.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"],subResources=["status"]
type Catalog walruscore.Catalog

var _ runtime.Object = (*Catalog)(nil)

// CatalogList holds the list of Catalog.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CatalogList walruscore.CatalogList

var _ runtime.Object = (*CatalogList)(nil)
