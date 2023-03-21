package types

import (
	"crypto/md5" // #nosec
	"encoding/hex"
	"reflect"
	"sort"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

type ApplicationResourceEndpoint struct {
	// EndpointType is the extra info for application resource type, like nodePort, loadBalance.
	EndpointType string `json:"endpointType,omitempty"`
	// Endpoints are the access endpoints.
	Endpoints []string `json:"endpoints,omitempty"`
}

type ApplicationResourceStatus struct {
	status.Status `json:",inline"`

	ResourceEndpoints ApplicationResourceEndpoints `json:"resourceEndpoints,omitempty"`
}

func (a ApplicationResourceStatus) Equal(newArs ApplicationResourceStatus) bool {
	// status
	if !a.Status.Equal(newArs.Status) {
		return false
	}

	// endpoints
	return a.ResourceEndpoints.Equal(newArs.ResourceEndpoints)
}

type ApplicationResourceEndpoints []ApplicationResourceEndpoint

func (a ApplicationResourceEndpoints) Equal(eps ApplicationResourceEndpoints) bool {
	if len(a) != len(eps) {
		return false
	}
	sortEndpoints := func(eps ApplicationResourceEndpoints) {
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
