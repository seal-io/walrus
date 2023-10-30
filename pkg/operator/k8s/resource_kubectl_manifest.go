package k8s

import (
	"regexp"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/seal-io/walrus/pkg/dao/model"
)

// parseResourcesOfKubectlManifest parses the given `kubectl_manifest` model.ResourceComponent,
// and keeps resource item which matches.
func parseResourcesOfKubectlManifest(
	sr *model.ResourceComponent,
	match func(versionResource schema.GroupVersionResource) bool,
) ([]resource, error) {
	if match == nil {
		return nil, nil
	}

	res := getResourceFromAPIPath(sr.Name)
	if res == nil || !match(res.GroupVersionResource) {
		return nil, nil
	}

	return []resource{*res}, nil
}

// getResourceFromAPIPath get resource from kubernetes api path. E.g., /api/v1/namespaces/foo/services/bar.
func getResourceFromAPIPath(apiPath string) *resource {
	namespacedCoreRegex := regexp.MustCompile(`/api/v1/namespaces/([^/]+)/([^/]+)/([^/]+)`)
	namespacedApisRegex := regexp.MustCompile(`/apis/([^/]+)/([^/]+)/namespaces/([^/]+)/([^/]+)/([^/]+)`)
	nonNamespacedCoreRegex := regexp.MustCompile(`/api/v1/([^/]+)/([^/]+)`)
	nonNamespacedApisRegex := regexp.MustCompile(`/apis/([^/]+)/([^/]+)/([^/]+)/([^/]+)`)

	switch {
	case namespacedCoreRegex.MatchString(apiPath):
		matches := namespacedCoreRegex.FindStringSubmatch(apiPath)

		return &resource{
			GroupVersionResource: schema.GroupVersionResource{
				Version:  corev1.SchemeGroupVersion.Version,
				Resource: matches[2],
			},
			Namespace: matches[1],
			Name:      matches[3],
		}
	case namespacedApisRegex.MatchString(apiPath):
		matches := namespacedApisRegex.FindStringSubmatch(apiPath)

		return &resource{
			GroupVersionResource: schema.GroupVersionResource{
				Group:    matches[1],
				Version:  matches[2],
				Resource: matches[4],
			},
			Namespace: matches[3],
			Name:      matches[5],
		}
	case nonNamespacedCoreRegex.MatchString(apiPath):
		matches := nonNamespacedCoreRegex.FindStringSubmatch(apiPath)

		return &resource{
			GroupVersionResource: schema.GroupVersionResource{
				Version:  corev1.SchemeGroupVersion.Version,
				Resource: matches[1],
			},
			Name: matches[2],
		}
	case nonNamespacedApisRegex.MatchString(apiPath):
		matches := nonNamespacedApisRegex.FindStringSubmatch(apiPath)

		return &resource{
			GroupVersionResource: schema.GroupVersionResource{
				Group:    matches[1],
				Version:  matches[2],
				Resource: matches[3],
			},
			Name: matches[4],
		}
	}

	return nil
}
