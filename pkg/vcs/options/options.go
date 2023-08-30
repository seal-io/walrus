package options

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/drone/go-scm/scm"
	scmtransport "github.com/drone/go-scm/scm/transport"

	"github.com/seal-io/walrus/utils/version"
)

type ClientOption func(*http.Client)

// WithUserAgent sets the user agent string to use for requests.
func WithUserAgent(userAgent string) ClientOption {
	return func(cli *http.Client) {
		if userAgent == "" {
			return
		}

		const headerKey = "User-Agent"
		cli.Transport = &scmtransport.Custom{
			Base: cli.Transport,
			Before: func(r *http.Request) {
				if r.Header.Get(headerKey) == "" {
					r.Header.Set(headerKey, userAgent)
				}
			},
		}
	}
}

// WithInsecureSkipVerify disables SSL certificate verification.
func WithInsecureSkipVerify() ClientOption {
	return func(cli *http.Client) {
		tr, ok := cli.Transport.(*scmtransport.Custom)
		if !ok {
			return
		}

		for b := tr.Base; b != nil; {
			switch v := b.(type) {
			case *scmtransport.Custom:
				b = v.Base
				continue
			case *http.Transport:
				if v.TLSClientConfig == nil {
					v.TLSClientConfig = &tls.Config{
						MinVersion: tls.VersionTLS12,
					}
				}
				v.TLSClientConfig.InsecureSkipVerify = true
			}

			return
		}
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(cli *http.Client) {
		cli.Timeout = timeout
	}
}

// WithToken sets the OAuth token.
func WithToken(token string) ClientOption {
	return func(cli *http.Client) {
		if token == "" {
			return
		}

		const headerKey = "Authorization"
		cli.Transport = &scmtransport.Custom{
			Base: cli.Transport,
			Before: func(r *http.Request) {
				if r.Header.Get(headerKey) == "" {
					r.Header.Set(headerKey, "Bearer "+token)
				}
			},
		}
	}
}

// SetClientOptions sets the client options.
func SetClientOptions(client *scm.Client, opts ...ClientOption) {
	httpCli := client.Client
	if httpCli == nil {
		httpCli = &http.Client{}
		client.Client = httpCli
	}

	// Default timeout.
	if httpCli.Timeout == 0 {
		httpCli.Timeout = 15 * time.Second
	}

	// Default transport.
	if httpCli.Transport == nil {
		httpCli.Transport = &scmtransport.Custom{
			Base: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
				},
			},
			Before: func(r *http.Request) {},
		}
	} else {
		httpCli.Transport = &scmtransport.Custom{
			Base:   httpCli.Transport,
			Before: func(r *http.Request) {},
		}
	}

	// Default user agent.
	WithUserAgent(version.GetUserAgent())(httpCli)

	// Apply options.
	for i := range opts {
		if opts[i] == nil {
			continue
		}

		opts[i](httpCli)
	}
}
