package types

import (
	"crypto/md5" // #nosec
	"encoding/hex"
	"reflect"
	"sort"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

const (
	// ServiceResourceModeManaged indicates the resource created to target platform,
	// it is writable(update or delete).
	ServiceResourceModeManaged = "managed"

	// ServiceResourceModeData indicates the resource read from target platform,
	// it is read-only.
	ServiceResourceModeData = "data"

	// ServiceResourceModeDiscovered indicates the resource discovered from target platform,
	// it inherits its composition's characteristic to be writable or not.
	ServiceResourceModeDiscovered = "discovered"
)

const (
	// ServiceResourceShapeClass indicates the resource is a class.
	ServiceResourceShapeClass = "class"
	// ServiceResourceShapeInstance indicates the resource is an instance that implements a class.
	ServiceResourceShapeInstance = "instance"
)

// ServiceResourceRelationshipTypeDependency indicates the relationship between service resource and its dependencies.
const ServiceResourceRelationshipTypeDependency = "Dependency"

type ServiceResourceEndpoint struct {
	// EndpointType is the extra info for service resource type, like nodePort, loadBalance.
	EndpointType string `json:"endpointType,omitempty"`
	// Endpoints are the access endpoints.
	Endpoints []string `json:"endpoints,omitempty"`
}

type ServiceResourceStatus struct {
	status.Status `json:",inline"`

	ResourceEndpoints ServiceResourceEndpoints `json:"resourceEndpoints,omitempty"`
}

func (a ServiceResourceStatus) Equal(newArs ServiceResourceStatus) bool {
	// Status.
	if !a.Status.Equal(newArs.Status) {
		return false
	}

	// Endpoints.
	return a.ResourceEndpoints.Equal(newArs.ResourceEndpoints)
}

type ServiceResourceEndpoints []ServiceResourceEndpoint

func (a ServiceResourceEndpoints) Equal(eps ServiceResourceEndpoints) bool {
	if len(a) != len(eps) {
		return false
	}
	sortEndpoints := func(eps ServiceResourceEndpoints) {
		for i := range eps {
			sort.Strings(eps[i].Endpoints)
		}

		sort.SliceStable(eps, func(i, j int) bool {
			if eps[i].EndpointType != eps[j].EndpointType {
				return eps[i].EndpointType < eps[j].EndpointType
			}

			if len(eps[i].Endpoints) != len(eps[j].Endpoints) {
				return len(eps[i].Endpoints) < len(eps[j].Endpoints)
			}

			if len(eps[i].Endpoints) != len(eps[j].Endpoints) {
				return len(eps[i].Endpoints) < len(eps[j].Endpoints)
			}

			return hashAddrs(eps[i].Endpoints) < hashAddrs(eps[j].Endpoints)
		})
	}

	sortEndpoints(a)
	sortEndpoints(eps)

	return reflect.DeepEqual(a, eps)
}

func hashAddrs(addrs []string) string {
	h := md5.New() // #nosec
	for _, v := range addrs {
		h.Write([]byte(v))
	}

	return hex.EncodeToString(h.Sum(nil))
}

type ServiceResourceOperationKeys struct {
	// Labels stores label of layer,
	// its length means each key contains levels with the same value as level.
	Labels []string `json:"labels,omitempty"`
	// Keys stores key in tree.
	Keys []ServiceResourceOperationKey `json:"keys,omitempty"`
}

// ServiceResourceOperationKey holds hierarchy query keys.
type ServiceResourceOperationKey struct {
	// Keys indicates the subordinate keys,
	// usually, it should not be valued in leaves.
	Keys []ServiceResourceOperationKey `json:"keys,omitempty"`
	// Name indicates the name of the key.
	Name string `json:"name"`
	// Value indicates the value of the key,
	// usually, it should be valued in leaves.
	Value string `json:"value,omitempty"`
	// Loggable indicates whether to be able to get log.
	Loggable *bool `json:"loggable,omitempty"`
	// Executable indicates whether to be able to execute remote command.
	Executable *bool `json:"executable,omitempty"`
}
