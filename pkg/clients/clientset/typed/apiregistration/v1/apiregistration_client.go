// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

import (
	"net/http"

	"github.com/seal-io/walrus/pkg/clients/clientset/scheme"
	rest "k8s.io/client-go/rest"
	v1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
)

type ApiregistrationV1Interface interface {
	RESTClient() rest.Interface
	APIServicesGetter
}

// ApiregistrationV1Client is used to interact with features provided by the apiregistration.k8s.io group.
type ApiregistrationV1Client struct {
	restClient rest.Interface
}

func (c *ApiregistrationV1Client) APIServices() APIServiceInterface {
	return newAPIServices(c)
}

// NewForConfig creates a new ApiregistrationV1Client for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*ApiregistrationV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

// NewForConfigAndClient creates a new ApiregistrationV1Client for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*ApiregistrationV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &ApiregistrationV1Client{client}, nil
}

// NewForConfigOrDie creates a new ApiregistrationV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *ApiregistrationV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new ApiregistrationV1Client for the given RESTClient.
func New(c rest.Interface) *ApiregistrationV1Client {
	return &ApiregistrationV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *ApiregistrationV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}