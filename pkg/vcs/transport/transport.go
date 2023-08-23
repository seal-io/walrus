package transport

import "net/http"

// UserAgentTransport is an http.RoundTripper that sets a User-Agent header.
type UserAgentTransport struct {
	Base http.RoundTripper

	// UserAgent string to set in the request header.
	UserAgent string
}

func (t *UserAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", t.UserAgent)
	return t.base().RoundTrip(r)
}

func (t *UserAgentTransport) base() http.RoundTripper {
	if t.Base == nil {
		return http.DefaultTransport
	}

	return t.Base
}
