package types

import (
	"crypto/md5" // #nosec
	"encoding/hex"
	"encoding/json"
	"reflect"
	"sort"
	"time"

	"github.com/seal-io/walrus/pkg/dao/types/status"
)

const (
	// ResourceComponentModeManaged indicates the resource created to target platform,
	// it is writable(update or delete).
	ResourceComponentModeManaged = "managed"

	// ResourceComponentModeData indicates the resource read from target platform,
	// it is read-only.
	ResourceComponentModeData = "data"

	// ResourceComponentModeDiscovered indicates the resource discovered from target platform,
	// it inherits its composition's characteristic to be writable or not.
	ResourceComponentModeDiscovered = "discovered"
)

const (
	// ResourceComponentShapeClass indicates the resource is a class.
	ResourceComponentShapeClass = "class"
	// ResourceComponentShapeInstance indicates the resource is an instance that implements a class.
	ResourceComponentShapeInstance = "instance"
)

// ResourceComponentRelationshipTypeDependency indicates the relationship between resource component
// and its dependencies.
const ResourceComponentRelationshipTypeDependency = "Dependency"

type ResourceComponentEndpoint struct {
	// EndpointType is the extra info for resource component type, like nodePort, loadBalance.
	EndpointType string `json:"endpointType,omitempty"`
	// Endpoints are the access endpoints.
	Endpoints []string `json:"endpoints,omitempty"`
}

type ResourceComponentStatus struct {
	status.Status `json:",inline"`

	ResourceEndpoints ResourceComponentEndpoints `json:"resourceEndpoints,omitempty"`
}

func (a ResourceComponentStatus) Equal(newArs ResourceComponentStatus) bool {
	// Status.
	if !a.Status.Equal(newArs.Status) {
		return false
	}

	// Endpoints.
	return a.ResourceEndpoints.Equal(newArs.ResourceEndpoints)
}

type ResourceComponentEndpoints []ResourceComponentEndpoint

func (a ResourceComponentEndpoints) Equal(eps ResourceComponentEndpoints) bool {
	if len(a) != len(eps) {
		return false
	}
	sortEndpoints := func(eps ResourceComponentEndpoints) {
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

type ResourceComponentOperationKeys struct {
	// Labels stores label of layer,
	// its length means each key contains levels with the same value as level.
	Labels []string `json:"labels,omitempty"`
	// Keys stores key in tree.
	Keys []ResourceComponentOperationKey `json:"keys,omitempty"`
}

// ResourceComponentOperationKey holds hierarchy query keys.
type ResourceComponentOperationKey struct {
	// Keys indicates the subordinate keys,
	// usually, it should not be valued in leaves.
	Keys []ResourceComponentOperationKey `json:"keys,omitempty"`
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

type ResourceComponentDriftDetection struct {
	Drifted bool `json:"drifted"`
	// Time indicates the time when the resource component is drifted.
	Time time.Time `json:"time"`
	// Result indicates the drift result of resource component.
	Result *ResourceComponentDrift `json:"drift"`
}

type ResourceComponentDrift struct {
	Address       string  `json:"address"`
	ModuleAddress string  `json:"module_address"`
	Mode          string  `json:"mode"`
	Type          string  `json:"type"`
	Name          string  `json:"name"`
	ProviderName  string  `json:"provider_name"`
	Change        *Change `json:"change"`
}

type Change struct {
	Actions         []string        `json:"actions"`
	Before          json.RawMessage `json:"before"`
	After           json.RawMessage `json:"after"`
	AfterUnknown    json.RawMessage `json:"after_unknown"`
	BeforeSensitive json.RawMessage `json:"before_sensitive"`
	AfterSensitive  json.RawMessage `json:"after_sensitive"`
}
