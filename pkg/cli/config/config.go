package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/henvic/httpretty"

	"github.com/seal-io/walrus/utils/strs"
	"github.com/seal-io/walrus/utils/version"
)

const (
	timeout     = 15 * time.Second
	openAPIPath = "/openapi"
	apiVersion  = "v1"
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
		ep, err := url.Parse(c.Server)
		if err != nil {
			return nil, err
		}

		resolved := ep.ResolveReference(req.URL)
		req.URL = resolved
	}

	c.SetHeaders(req)

	resp, err := c.HttpClient().Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ValidateAndSetup validate and setup the context.
func (c *Config) ValidateAndSetup() error {
	if c.Server == "" {
		return errors.New("endpoint is required")
	}

	if c.Token == "" {
		return errors.New("token is required")
	}

	err := c.validateProject()
	if err != nil {
		return err
	}

	err = c.validateEnvironment()
	if err != nil {
		return err
	}

	return nil
}

// HttpClient generate http client base on context config.
func (c *Config) HttpClient() *http.Client {
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

// SetHeaders set default headers.
func (c *Config) SetHeaders(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("User-Agent", version.GetUserAgentWith("walrus-cli"))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")
}

func (c *Config) SetHost(req *http.Request) error {
	if req.URL.Host != "" {
		return nil
	}

	ep, err := url.Parse(c.Server)
	if err != nil {
		return err
	}

	resolved := ep.ResolveReference(req.URL)
	req.URL = resolved

	return nil
}

const (
	projectResource     = "projects"
	environmentResource = "environments"
)

// validateProject validate project name base on server context.
func (c *Config) validateProject() error {
	if c.Project == "" {
		return nil
	}

	address := path.Join(apiVersion, projectResource, c.Project)
	err := c.validateResourceItem(projectResource, c.Project, address)

	return err
}

// validateEnvironment validate environment name base on server context, project name.
func (c *Config) validateEnvironment() error {
	if c.Environment == "" {
		return nil
	}

	address := path.Join(apiVersion, projectResource, c.Project, environmentResource, c.Environment)
	err := c.validateResourceItem(environmentResource, c.Environment, address)

	return err
}

// validateResourceItem send get resource request to server.
func (c *Config) validateResourceItem(resource, name, address string) error {
	req, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		return err
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return err
	}

	switch {
	default:
	case resp.StatusCode == http.StatusUnauthorized:
		return fmt.Errorf("unauthorized, please check the validity of the token")
	case resp.StatusCode == http.StatusNotFound:
		return fmt.Errorf("%s %s not found", strs.Singularize(resource), name)
	case resp.StatusCode < 200 || resp.StatusCode >= 300:
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}
