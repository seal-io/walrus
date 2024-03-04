package httpx

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type ClientOption struct {
	transport  *http.Transport
	timeout    time.Duration
	debug      bool
	roundTrips []func(req *http.Request) error
}

func ClientOptions() *ClientOption {
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSClientConfig:     &tls.Config{MinVersion: tls.VersionTLS12},
		DialContext:         (&net.Dialer{KeepAlive: -1}).DialContext,
		ForceAttemptHTTP2:   true,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &ClientOption{
		transport: transport,
		timeout:   30 * time.Second,
	}
}

// WithTimeout sets the request timeout.
func (o *ClientOption) WithTimeout(timeout time.Duration) *ClientOption {
	if timeout >= 0 {
		o.timeout = timeout
	}
	return o
}

// WithDebug sets the debug mode.
func (o *ClientOption) WithDebug() *ClientOption {
	o.debug = true
	return o
}

// WithoutProxy disables the proxy.
func (o *ClientOption) WithoutProxy() *ClientOption {
	o.transport.Proxy = nil
	return o
}

// SkipInsecureVerify skips the insecure verify.
func (o *ClientOption) SkipInsecureVerify() *ClientOption {
	o.transport.TLSClientConfig.InsecureSkipVerify = true
	return o
}

// WithDial sets the dial function.
func (o *ClientOption) WithDial(dial func(context.Context, string, string) (net.Conn, error)) *ClientOption {
	o.transport.DialContext = dial
	return o
}

// WithRoundTrip sets the round trip function.
func (o *ClientOption) WithRoundTrip(rt func(req *http.Request) error) *ClientOption {
	if rt == nil {
		return o
	}
	o.roundTrips = append(o.roundTrips, rt)
	return o
}

// WithUserAgent sets the user agent.
func (o *ClientOption) WithUserAgent(ua string) *ClientOption {
	return o.WithRoundTrip(func(req *http.Request) error {
		req.Header.Set("User-Agent", ua)
		return nil
	})
}

// WithBearerAuth sets the bearer token.
func (o *ClientOption) WithBearerAuth(token string) *ClientOption {
	return o.WithRoundTrip(func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	})
}

// WithBasicAuth sets the basic authentication.
func (o *ClientOption) WithBasicAuth(username, password string) *ClientOption {
	return o.WithRoundTrip(func(req *http.Request) error {
		req.SetBasicAuth(username, password)
		return nil
	})
}

// WithHeader sets the header.
func (o *ClientOption) WithHeader(key, value string) *ClientOption {
	return o.WithRoundTrip(func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	})
}

// WithHeaders sets the headers.
func (o *ClientOption) WithHeaders(headers map[string]string) *ClientOption {
	return o.WithRoundTrip(func(req *http.Request) error {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		return nil
	})
}
