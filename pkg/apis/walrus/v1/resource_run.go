package v1

import (
	"k8s.io/apimachinery/pkg/runtime"

	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
)

// ResourceRun is the schema for the resource runs API.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"],shortName=["resrun"],subResources=["status"]
type ResourceRun walruscore.ResourceRun

var _ runtime.Object = (*ResourceRun)(nil)

// ResourceRunList holds the list of ResourceRun.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceRunList walruscore.ResourceRunList

var _ runtime.Object = (*ResourceRunList)(nil)
