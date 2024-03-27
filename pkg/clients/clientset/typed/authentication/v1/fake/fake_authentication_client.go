// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package fake

import (
	v1 "github.com/seal-io/walrus/pkg/clients/clientset/typed/authentication/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeAuthenticationV1 struct {
	*testing.Fake
}

func (c *FakeAuthenticationV1) SelfSubjectReviews() v1.SelfSubjectReviewInterface {
	return &FakeSelfSubjectReviews{c}
}

func (c *FakeAuthenticationV1) TokenReviews() v1.TokenReviewInterface {
	return &FakeTokenReviews{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeAuthenticationV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}