package kubereviewsubject

import (
	"context"
	"errors"
	"fmt"
	"slices"

	authz "k8s.io/api/authorization/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	authzcli "k8s.io/client-go/kubernetes/typed/authorization/v1"
	"k8s.io/klog/v2"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	// Review holds the attributes for advanced reviewing.
	Review = authz.SubjectAccessReviewSpec
	// Reviews is the list of Review.
	Reviews = []Review

	// DeniedError is an error indicate which Review target has been denied.
	DeniedError struct {
		// Review holds the
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
		klog.Error(err, "ignored subject review denied error, need fixing manually")
	}
	return nil
}

// CanDo checks if the given subject can do the specified actions.
func CanDo(ctx context.Context, cli authzcli.SubjectAccessReviewInterface, reviews Reviews) error {
	if len(reviews) == 0 {
		return errors.New("no self review to check")
	}

	allowed := true
	for i := range reviews {
		sar := &authz.SubjectAccessReview{
			Spec: reviews[i],
		}

		sar, err := cli.Create(ctx, sar, meta.CreateOptions{})
		if err != nil {
			return fmt.Errorf("create subject access review %s: %w", reviews[i], err)
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
		sar := &authz.SubjectAccessReview{
			Spec: reviews[i],
		}

		err := cli.Create(ctx, sar, &ctrlcli.CreateOptions{})
		if err != nil {
			return fmt.Errorf("create subject access review %s: %w", reviews[i], err)
		}

		allowed = allowed && sar.Status.Allowed
		if !allowed {
			return DeniedError{Review: reviews[i]}
		}
	}

	return nil
}

// CanRequestUserDo leverages CanDo to review the requesting subject can do the specified actions or not.
//
// CanRequestUserDo overrides all given reviews' user information
// after it successfully parse from the given context.
func CanRequestUserDo(ctx context.Context, cli authzcli.SubjectAccessReviewInterface, reviews Reviews) error {
	reviews, err := overrideUserInfo(ctx, reviews)
	if err != nil {
		return err
	}

	return CanDo(ctx, cli, reviews)
}

// CanRequestUserDoWithCtrlClient is similar to CanRequestUserDo, but uses the ctrl client.
func CanRequestUserDoWithCtrlClient(ctx context.Context, cli ctrlcli.Client, reviews Reviews) error {
	reviews, err := overrideUserInfo(ctx, reviews)
	if err != nil {
		return err
	}

	return CanDoWithCtrlClient(ctx, cli, reviews)
}

// overrideUserInfo extracts the user info from the given context,
// and then override the user information of the given Reviews with the extracted result.
func overrideUserInfo(ctx context.Context, reviews Reviews) (Reviews, error) {
	ui, ok := genericapirequest.UserFrom(ctx)
	if !ok {
		return nil, errors.New("cannot retrieve kubernetes request user information from context")
	}

	var (
		user  = ui.GetName()
		uid   = ui.GetUID()
		extra = func() (out map[string]authz.ExtraValue) {
			in := ui.GetExtra()
			if len(in) == 0 {
				return
			}
			out = make(map[string]authz.ExtraValue, len(in))
			for i := range in {
				if len(in[i]) == 0 {
					continue
				}
				out[i] = slices.Clone(in[i])
			}
			return
		}()
		groups = ui.GetGroups()
	)

	// Override user information.
	for i := range reviews {
		reviews[i].User = user
		reviews[i].UID = uid
		reviews[i].Extra = extra
		reviews[i].Groups = groups
	}

	return reviews, nil
}
