// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testing "k8s.io/client-go/testing"
)

// FakeTokenReviews implements TokenReviewInterface
type FakeTokenReviews struct {
	Fake *FakeAuthenticationV1
}

var tokenreviewsResource = v1.SchemeGroupVersion.WithResource("tokenreviews")

var tokenreviewsKind = v1.SchemeGroupVersion.WithKind("TokenReview")

// Create takes the representation of a tokenReview and creates it.  Returns the server's representation of the tokenReview, and an error, if there is any.
func (c *FakeTokenReviews) Create(ctx context.Context, tokenReview *v1.TokenReview, opts metav1.CreateOptions) (result *v1.TokenReview, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(tokenreviewsResource, tokenReview), &v1.TokenReview{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.TokenReview), err
}
