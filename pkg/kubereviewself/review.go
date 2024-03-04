package kubereviewself

import (
	"context"
	"errors"
	"fmt"

	authz "k8s.io/api/authorization/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	authzcli "k8s.io/client-go/kubernetes/typed/authorization/v1"
	"k8s.io/klog/v2"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	// Review holds the attributes for advanced reviewing.
	Review = authz.SelfSubjectAccessReviewSpec
	// Reviews is the list of Review.
	Reviews = []Review

	// DeniedError is an error indicate which Review target has been denied.
	DeniedError struct {
		Review Review
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
		klog.Error(err, "ignored self review denied error, need fixing manually")
	}
	return nil
}

// CanDo checks if the current subject can do the specified actions.
func CanDo(ctx context.Context, cli authzcli.SelfSubjectAccessReviewInterface, reviews Reviews) error {
	if len(reviews) == 0 {
		return errors.New("no self review to check")
	}

	allowed := true
	for i := range reviews {
		sar := &authz.SelfSubjectAccessReview{
			Spec: reviews[i],
		}

		sar, err := cli.Create(ctx, sar, meta.CreateOptions{})
		if err != nil {
			return fmt.Errorf("create self subject access review %s: %w", reviews[i], err)
		}

		allowed = allowed && sar.Status.Allowed
		if !allowed {
			return DeniedError{Review: reviews[i]}
		}
	}

	return nil
}

// CanDoWithCtrlClient is similar to CanDo, but uses the ctrl client.
func CanDoWithCtrlClient(ctx context.Context, cli ctrlcli.Client, reviews Reviews) error {
	if len(reviews) == 0 {
		return errors.New("no self review to check")
	}

	allowed := true
	for i := range reviews {
		sar := &authz.SelfSubjectAccessReview{
			Spec: reviews[i],
		}

		err := cli.Create(ctx, sar, &ctrlcli.CreateOptions{})
		if err != nil {
			return fmt.Errorf("create self subject access review %s: %w", reviews[i], err)
		}

		allowed = allowed && sar.Status.Allowed
		if !allowed {
			return DeniedError{Review: reviews[i]}
		}
	}

	return nil
}
