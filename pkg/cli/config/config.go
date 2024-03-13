package config

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/henvic/httpretty"

	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/version"
)

const (
	timeout         = 15 * time.Second
	openAPIPath     = "/openapi"
	apiVersion      = "v1"
	versionPath     = "/debug/version"
	accountInfoPath = "/account/info"
)

// CommonConfig indicate the common CLI command config.
type CommonConfig struct {
	Debug bool `json:"debug,omitempty" yaml:"debug,omitempty"`
}

// Config include the common config and server context config.
type Config struct {
	CommonConfig  `json:",inline" yaml:",inline" mapstructure:",squash"`
	ServerContext `json:",inline" yaml:",inline" mapstructure:",squash"`
}

// Version include the version and commit.
type Version struct {
	Version string `json:"version" yaml:"version"`
	Commit  string `json:"commit" yaml:"commit"`
}

// DoRequestWithTimeout send request to server with timeout.
func (c *Config) DoRequestWithTimeout(req *http.Request, timeout time.Duration) (*http.Response, error) {
	if req.URL.Host == "" {
		ep, err := url.Parse(c.Server)
		if err != nil {
			return nil, err
		}

		resolved := ep.ResolveReference(req.URL)
		req.URL = resolved
	}

	c.SetHeaders(req)

	resp, err := c.HttpClient(timeout).Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DoRequest send request to server.
func (c *Config) DoRequest(req *http.Request) (*http.Response, error) {
	return c.DoRequestWithTimeout(req, timeout)
}

func (c *Config) CheckReachable() (err error) {
	// Reachable set.
	c.Reachable = true

	defer func() {
		if err != nil {
			c.Reachable = false
		}
	}()

	var u *url.URL

	u, err = url.Parse(c.Server)
	if err != nil {
		return err
	}

	if u.Port() == "" {
		if u.Scheme == "http" {
			u.Host += ":80"
		} else if u.Scheme == "https" {
			u.Host += ":443"
		}
	}

	var conn net.Conn

	conn, err = net.DialTimeout("tcp", u.Host, 3*time.Second)
	if err != nil {
		return fmt.Errorf("access %s failed %w", c.Server, err)
	}

	defer conn.Close()

	return nil
}

// ValidateAndSetup validate and setup the context.
func (c *Config) ValidateAndSetup() (err error) {
	// Reachable set.
	c.Reachable = true

	defer func() {
		if err != nil {
			c.Reachable = false
		}
	}()

	msg := `cli configuration is invalid: %v. You can configure cli by running "walrus login"`

	if c.Server == "" {
		return fmt.Errorf(msg, "server address is empty")
	}

	err = c.CheckReachable()
	if err != nil {
		return err
	}

	switch {
	case c.Project != "" && c.Environment != "":
		err = c.validateEnvironment()
	case c.Project != "":
		err = c.validateProject()
	default:
		err = c.validateAccountInfo()
	}

	if err != nil {
		return fmt.Errorf(msg, err)
	}

	return nil
}

// HttpClient generate http client base on context config.
func (c *Config) HttpClient(timeout time.Duration) *http.Client {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if c.Insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	var tp http.RoundTripper = &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: tlsConfig,
		DialContext: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
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
	err := c.validateResourceItem(address)

	return err
}

// validateEnvironment validate environment name base on server context, project name.
func (c *Config) validateEnvironment() error {
	if c.Environment == "" {
		return nil
	}

	address := path.Join(apiVersion, projectResource, c.Project, environmentResource, c.Environment)
	err := c.validateResourceItem(address)

	return err
}

// validateResourceItem send get resource request to server.
func (c *Config) validateResourceItem(address string) error {
	req, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		return err
	}

	resp, err := c.DoRequestWithTimeout(req, timeout)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return common.CheckResponseStatus(resp)
}

// validateAccountInfo send get account info request to server.
func (c *Config) validateAccountInfo() error {
	req, err := http.NewRequest(http.MethodGet, accountInfoPath, nil)
	if err != nil {
		return err
	}

	resp, err := c.DoRequestWithTimeout(req, timeout)
	if err != nil {
		return fmt.Errorf("access %s failed", c.Server)
	}
	defer resp.Body.Close()

	return common.CheckResponseStatus(resp)
}

func (c *Config) ServerVersion() (*Version, error) {
	req, err := http.NewRequest(http.MethodGet, versionPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.DoRequestWithTimeout(req, timeout)
	if err != nil {
		return nil, fmt.Errorf("access %s failed", c.Server)
	}
	defer resp.Body.Close()

	err = common.CheckResponseStatus(resp)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v Version
	err = json.Unmarshal(b, &v)

	return &v, err
}
