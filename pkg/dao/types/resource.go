package types

import "sort"

const (
	// ResourceRelationshipTypeImplicit indicates the resource dependency is auto created by resource reference.
	ResourceRelationshipTypeImplicit = "Implicit"
	// ResourceRelationshipTypeExplicit indicates the resource dependency is manually created by user.
	ResourceRelationshipTypeExplicit = "Explicit"
)

type (
	// ResourceEndpoint holds the endpoint definition of resource.
	ResourceEndpoint struct {
		// Name indicates the name of endpoint.
		Name string `json:"name"`
		// URL indicates the URL of endpoint, must be a valid URL.
		URL string `json:"url"`
	}

	// ResourceEndpoints holds a list of the endpoints definitions of resource.
	ResourceEndpoints []ResourceEndpoint
)

// ResourceEndpointsFromMap converts a map to ResourceEndpoints.
func ResourceEndpointsFromMap(m map[string]string) ResourceEndpoints {
	eps := make(ResourceEndpoints, 0, len(m))
	for name, url := range m {
		eps = append(eps, ResourceEndpoint{
			Name: name,
			URL:  url,
		})
	}

	return eps
}

func (in ResourceEndpoints) Len() int {
	return len(in)
}

func (in ResourceEndpoints) Less(i, j int) bool {
	return in[i].Name < in[j].Name
}

func (in ResourceEndpoints) Swap(i, j int) {
	in[i], in[j] = in[j], in[i]
}

func (in ResourceEndpoints) Sort() ResourceEndpoints {
	sort.Sort(in)
	return in
}
