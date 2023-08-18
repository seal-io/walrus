package k8s

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestGetResourceFromAPIPath(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected *resource
	}{
		{
			name: "namespaced resource in core api",
			path: "/api/v1/namespaces/default/services/example",
			expected: &resource{
				GroupVersionResource: schema.GroupVersionResource{
					Version:  "v1",
					Resource: "services",
				},
				Namespace: "default",
				Name:      "example",
			},
		},
		{
			name: "namespaced resource in apis",
			path: "/apis/batch/v1/namespaces/default/jobs/example",
			expected: &resource{
				GroupVersionResource: schema.GroupVersionResource{
					Group:    "batch",
					Version:  "v1",
					Resource: "jobs",
				},
				Namespace: "default",
				Name:      "example",
			},
		},
		{
			name: "non-namespaced resource in core api",
			path: "/api/v1/nodes/example",
			expected: &resource{
				GroupVersionResource: schema.GroupVersionResource{
					Version:  "v1",
					Resource: "nodes",
				},
				Namespace: "",
				Name:      "example",
			},
		},
		{
			name: "non-namespaced resource in apis",
			path: "/apis/rbac.authorization.k8s.io/v1/clusterroles/example",
			expected: &resource{
				GroupVersionResource: schema.GroupVersionResource{
					Group:    "rbac.authorization.k8s.io",
					Version:  "v1",
					Resource: "clusterroles",
				},
				Namespace: "",
				Name:      "example",
			},
		},
		{
			name:     "not a valid k8s api path",
			path:     "/something/else",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		actual := getResourceFromAPIPath(tc.path)
		assert.Equal(t, tc.expected, actual, fmt.Sprintf("unexpected result in test case: %s", tc.name))
	}
}
