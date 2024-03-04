package v1

import (
	"k8s.io/apimachinery/pkg/runtime"

	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
)

// Resource is the schema for the resources API.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"],shortName=["res"],subResources=["status"]
type Resource walruscore.Resource

var _ runtime.Object = (*Resource)(nil)

// ResourceList holds the list of Resource.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceList walruscore.ResourceList

var _ runtime.Object = (*ResourceList)(nil)
