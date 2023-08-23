package options

import (
	"crypto/tls"
	"net/http"
	"time"

	scmtransport "github.com/drone/go-scm/scm/transport"

	"github.com/drone/go-scm/scm"

	"github.com/seal-io/walrus/pkg/vcs/transport"
	"github.com/seal-io/walrus/utils/version"
)

type ClientOptions struct {
	// UserAgent is the user agent string to use for requests.
	UserAgent string

	// SkipVerify disables SSL certificate verification.
	SkipVerify bool

	// Timeout is the request timeout.
	Timeout time.Duration

	// Token is the OAuth token.
	Token string
}

type ClientOption func(*ClientOptions)

// WithUserAgent sets the user agent string to use for requests.
func WithUserAgent(userAgent string) ClientOption {
	return func(opts *ClientOptions) {
		opts.UserAgent = userAgent
	}
}

// WithSkipVerify disables SSL certificate verification.
func WithSkipVerify(skipVerify bool) ClientOption {
	return func(opts *ClientOptions) {
		opts.SkipVerify = skipVerify
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(opts *ClientOptions) {
		opts.Timeout = timeout
	}
}

// WithToken sets the OAuth token.
func WithToken(token string) ClientOption {
	return func(opts *ClientOptions) {
		opts.Token = token
	}
}

// SetClientOptions sets the client options.
func SetClientOptions(client *scm.Client, opts ...ClientOption) {
	options := ClientOptions{
		UserAgent:  version.GetUserAgent(),
		SkipVerify: false,
		Timeout:    15 * time.Second,
	}

	for _, opt := range opts {
		opt(&options)
	}

	if client.Client == nil {
		client.Client = &http.Client{}
	}

	// Set the request timeout.
	client.Client.Timeout = options.Timeout

	//nolint:gosec
	baseTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: options.SkipVerify,
		},
	}

	// Set the oauth bearer token.
	bearerTransport := &scmtransport.BearerToken{
		Token: options.Token,
		Base:  baseTransport,
	}

	// Set the user agent.
	client.Client.Transport = &transport.UserAgentTransport{
		UserAgent: options.UserAgent,
		Base:      bearerTransport,
	}
}
