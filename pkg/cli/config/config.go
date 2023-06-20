package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/henvic/httpretty"

	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/version"
)

const (
	timeout     = 15 * time.Second
	openAPIPath = "/openapi"
)

// CommonConfig indicate the common CLI command config.
type CommonConfig struct {
	Debug  bool   `json:"debug,omitempty" yaml:"debug,omitempty"`
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
}

// Config include the common config and server context config.
type Config struct {
	CommonConfig  `json:",inline" yaml:",inline" mapstructure:",squash"`
	ServerContext `json:",inline" yaml:",inline" mapstructure:",squash"`
}

// DoRequest send request to server.
func (c *Config) DoRequest(req *http.Request) (*http.Response, error) {
	if req.URL.Host == "" {
		ep, err := url.Parse(c.Endpoint)
		if err != nil {
			return nil, err
		}

		resolved := ep.ResolveReference(req.URL)
		req.URL = resolved
	}

	c.setHeaders(req)

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ValidateAndSetup validate and setup the context.
func (c *Config) ValidateAndSetup() error {
	if c.Endpoint == "" {
		return errors.New("endpoint is required")
	}

	if c.Token == "" {
		return errors.New("token is required")
	}

	err := c.setProjectID()
	if err != nil {
		return err
	}

	err = c.setEnvironmentID()
	if err != nil {
		return err
	}

	return nil
}

// httpClient generate http client base on context config.
func (c *Config) httpClient() *http.Client {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if c.Insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	var tp http.RoundTripper = &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: tlsConfig,
	}

	if c.Debug {
		tp = c.logger().RoundTripper(tp)
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: tp,
	}
}

// logger use to record request.
func (c *Config) logger() *httpretty.Logger {
	return &httpretty.Logger{
		Time:           true,
		TLS:            true,
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
		Colors:         true,
		Formatters:     []httpretty.Formatter{&httpretty.JSONFormatter{}},
	}
}

// setHeaders set default headers.
func (c *Config) setHeaders(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("User-Agent", fmt.Sprintf("seal.io/seal-cli; version=%s", version.Get()))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
}

const (
	projectResource     = "projects"
	environmentResource = "environments"
)

// resourceItems represent the common resource list response.
type resourceItems struct {
	Items []resourceItem `json:"items"`
}

// resourceItem represent the resource.
type resourceItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// setProjectID set project id base on server context and project name.
func (c *Config) setProjectID() error {
	if c.ProjectName == "" {
		c.ProjectName = "default"
	}

	projectQuery := map[string]string{
		"query": c.ProjectName,
	}

	projectID, err := c.getResourceIDByName(projectResource, c.ProjectName, projectQuery)
	if err != nil {
		return err
	}
	c.ProjectID = projectID

	return nil
}

// setEnvironmentID set environment id base on server context, project id and environment name.
func (c *Config) setEnvironmentID() error {
	if c.EnvironmentName == "" {
		return nil
	}

	envQuery := map[string]string{
		"projectID": c.ProjectID,
	}

	envID, err := c.getResourceIDByName(environmentResource, c.EnvironmentName, envQuery)
	if err != nil {
		return err
	}

	c.EnvironmentID = envID

	return nil
}

// getResourceIDByName send request to server and get resource id by name.
func (c *Config) getResourceIDByName(resource, resourceName string, queries map[string]string) (string, error) {
	items, err := c.listResource(resource, queries)
	if err != nil {
		return "", err
	}

	for _, v := range items {
		if v.Name == resourceName {
			return v.ID, nil
		}
	}

	return "", fmt.Errorf("%s %s not found", resource, resourceName)
}

// listResource send list resource request to server.
func (c *Config) listResource(resourceName string, queries map[string]string) ([]resourceItem, error) {
	req, err := http.NewRequest(http.MethodGet, path.Join(c.APIVersion, resourceName), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range queries {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var items resourceItems
	body, err := io.ReadAll(resp.Body)

	defer func() { _ = resp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &items)
	if err != nil {
		return nil, err
	}

	return items.Items, nil
}
