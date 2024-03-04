package review

import (
	"context"
	"errors"
	"fmt"
	"slices"

	authz "k8s.io/api/authorization/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	authzcli "k8s.io/client-go/kubernetes/typed/authorization/v1"
	"k8s.io/klog/v2"
)

type (
	Advanced  = authz.SelfSubjectAccessReviewSpec
	Advanceds = []Advanced

	DeniedError struct {
		Review Advanced
	}
)

func (e DeniedError) Error() string {
	return fmt.Sprintf("denied %s", e.Review)
}

// IsDeniedError checks if the error is a DeniedError.
func IsDeniedError(err error) bool {
	return errors.As(err, &DeniedError{})
}

// Try ignores the DeniedError.
func Try(err error) error {
	if err != nil {
		if !IsDeniedError(err) {
			return err
		}
		klog.Error(err, "ignored, need fixing manually")
	}
	return nil
}

// CanDo checks if the current subject can do the specified actions.
func CanDo(ctx context.Context, cli authzcli.SelfSubjectAccessReviewInterface, reviews Advanceds) error {
	if len(reviews) == 0 {
		return errors.New("no review to check")
	}

	allowed := true
	for i := range reviews {
		sar := &authz.SelfSubjectAccessReview{
			Spec: reviews[i],
		}

		sar, err := cli.Create(ctx, sar, meta.CreateOptions{})
		if err != nil {
			return fmt.Errorf("create SelfSubjectAccessReview %s: %w", reviews[i], err)
		}

		allowed = allowed && sar.Status.Allowed
		if !allowed {
			return DeniedError{Review: reviews[i]}
		}
	}

	return nil
}

type (
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
	Simples = []Simple
)

type ApplyOption func([]string, string) ([]string, string)

// WithStatusOnly checks if the subject can do kubeclientset.Apply
// with kubeclientset.ApplyStatusOnly for something.
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

// WithUpdate checks if the subject can do kubeclientset.Create
// with kubeclientset.CreateOrUpdate for something.
func WithUpdate() CreateOption {
	return func(verbs []string, subres string) ([]string, string) {
		return append(verbs, "update"), subres
	}
}

// WithRecreate checks if the subject can do kubeclientset.Create
// with kubeclientset.CreateOrRecreate for something.
func WithRecreate() CreateOption {
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

// WithCreate checks if the subject can do kubeclientset.Update
// with kubeclientset.UpdateOrCreate for something.
func WithCreate() UpdateOption {
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
	advances := make(Advanceds, 0, len(reviews)*len(verbs))
	for i := range reviews {
		simple := &reviews[i]
		for _, verb := range verbs {
			advances = append(advances, Advanced{
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

	return CanDo(ctx, cli, advances)
}
