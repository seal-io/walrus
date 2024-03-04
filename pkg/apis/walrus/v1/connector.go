package v1

import (
	"k8s.io/apimachinery/pkg/runtime"

	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
)

// Connector is the schema for the connectors API.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:apireg-gen:resource:scope="Namespaced",categories=["walrus"],shortName=["conn"],subResources=["status"]
type Connector walruscore.Connector

var _ runtime.Object = (*Connector)(nil)

// ConnectorList holds the list of Connector.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ConnectorList walruscore.ConnectorList

var _ runtime.Object = (*ConnectorList)(nil)
