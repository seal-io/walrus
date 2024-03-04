package v1

import (
	"k8s.io/apimachinery/pkg/runtime"

	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
)

// ResourceDefinition is the schema for the resource definitions API.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"],shortName=["resdef"],subResources=["status"]
type ResourceDefinition walruscore.ResourceDefinition

var _ runtime.Object = (*ResourceDefinition)(nil)

// ResourceDefinitionList holds the list of ResourceDefinition.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceDefinitionList walruscore.ResourceDefinitionList

var _ runtime.Object = (*ResourceDefinitionList)(nil)
