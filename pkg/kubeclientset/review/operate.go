package review

import (
	"context"
	"slices"

	authz "k8s.io/api/authorization/v1"
	authzcli "k8s.io/client-go/kubernetes/typed/authorization/v1"

	"github.com/seal-io/walrus/pkg/kubereviewself"
)

type (
	// Simple holds the attributes for simply reviewing.
	Simple struct {
		// Namespace is the namespace of the action being requested.
		// Currently, there is no distinction between no namespace and all namespaces。
		Namespace string
		// Group is the API Group of the Resource.
		// "*" means all.
		Group string
		// Version is the API Version of the Resource.
		// "*" means all.
		Version string
		// Resource is one of the existing resource types.
		// "*" means all.
		Resource string
		// Name is the name of the resource being requested for a "get" or deleted for a "delete".
		// "" (empty) means all.
		Name string
	}
	// Simples is the list of Simple.
	Simples = []Simple
)

type ApplyOption func([]string, string) ([]string, string)

// WithStatusOnly checks if the subject can do kubeclientset.Apply
// with kubeclientset.WithStatusOnly for something.
func WithStatusOnly() ApplyOption {
	return func(verbs []string, _ string) ([]string, string) {
		return slices.DeleteFunc(verbs, func(s string) bool {
			return s == "create"
		}), "status"
	}
}

// CanDoApply checks if the current subject can do kubeclientset.Apply for something.
func CanDoApply(ctx context.Context,
	cli authzcli.SelfSubjectAccessReviewInterface, reviews Simples,
	opts ...ApplyOption,
) error {
	verbs := []string{"get", "create", "patch"}
	subres := ""
	for i := range opts {
		verbs, subres = opts[i](verbs, subres)
	}

	return canDoSimple(ctx, cli, reviews, verbs, subres)
}

type CreateOption func([]string, string) ([]string, string)

// WithUpdateIfExisted checks if the subject can do kubeclientset.Create
// with kubeclientset.WithUpdateIfExisted for something.
func WithUpdateIfExisted() CreateOption {
	return func(verbs []string, subres string) ([]string, string) {
		return append(verbs, "update"), subres
	}
}

// WithRecreateIfDuplicated checks if the subject can do kubeclientset.Create
// with kubeclientset.WithRecreateIfDuplicated for something.
func WithRecreateIfDuplicated() CreateOption {
	return func(verbs []string, subres string) ([]string, string) {
		return append(verbs, "delete"), subres
	}
}

// CanDoCreate checks if the current subject can do kubeclientset.Create for something.
func CanDoCreate(ctx context.Context,
	cli authzcli.SelfSubjectAccessReviewInterface, reviews Simples,
	opts ...CreateOption,
) error {
	verbs := []string{"get", "create"}
	subres := ""
	for i := range opts {
		verbs, subres = opts[i](verbs, subres)
	}

	return canDoSimple(ctx, cli, reviews, verbs, subres)
}

type UpdateOption func([]string, string) ([]string, string)

// WithCreateIfNotExisted checks if the subject can do kubeclientset.Update
// with kubeclientset.WithCreateIfNotExisted for something.
func WithCreateIfNotExisted() UpdateOption {
	return func(verbs []string, subres string) ([]string, string) {
		return append(verbs, "create"), subres
	}
}

// CanDoUpdate checks if the current subject can do kubeclientset.Update for something.
func CanDoUpdate(ctx context.Context,
	cli authzcli.SelfSubjectAccessReviewInterface, reviews Simples,
	opts ...UpdateOption,
) error {
	verbs := []string{"get", "update"}
	subres := ""
	for i := range opts {
		verbs, subres = opts[i](verbs, subres)
	}

	return canDoSimple(ctx, cli, reviews, verbs, subres)
}

// CanDoDelete checks if the current subject can do kubeclientset.Delete for something.
func CanDoDelete(ctx context.Context,
	cli authzcli.SelfSubjectAccessReviewInterface, reviews Simples,
) error {
	verbs := []string{"delete"}
	subres := ""

	return canDoSimple(ctx, cli, reviews, verbs, subres)
}

func canDoSimple(
	ctx context.Context,
	cli authzcli.SelfSubjectAccessReviewInterface, reviews Simples,
	verbs []string, subresource string,
) error {
	advances := make(kubereviewself.Reviews, 0, len(reviews)*len(verbs))
	for i := range reviews {
		simple := &reviews[i]
		for _, verb := range verbs {
			advances = append(advances, kubereviewself.Review{
				ResourceAttributes: &authz.ResourceAttributes{
					Verb:        verb,
					Group:       simple.Group,
					Version:     simple.Version,
					Resource:    simple.Resource,
					Subresource: subresource,
					Namespace:   simple.Namespace,
					Name:        simple.Name,
				},
			})
		}
	}

	return kubereviewself.CanDo(ctx, cli, advances)
}
