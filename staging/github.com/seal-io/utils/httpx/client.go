package httpx

import (
	"context"
	"io"
	"net/http"

	"github.com/henvic/httpretty"

	"github.com/seal-io/utils/pools/bytespool"
)

// DefaultClient is the default http.Client used by the package.
//
// It is used for requests pooling.
var DefaultClient = http.DefaultClient

// Client returns a new http.Client with the given options,
// the result http.Client is used for fast-consuming requests.
//
// If you want a requests pool management, use DefaultClient instead.
func Client(opts ...*ClientOption) *http.Client {
	var o *ClientOption
	if len(opts) > 0 {
		o = opts[0]
	} else {
		o = ClientOptions()
	}

	root := http.DefaultTransport
	if o.transport != nil {
		root = o.transport
	}

	if o.debug {
		pretty := &httpretty.Logger{
			Time:           true,
			TLS:            true,
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
			Formatters:     []httpretty.Formatter{&JSONFormatter{}},
		}
		root = pretty.RoundTripper(root)
	}

	transport := RoundTripperChain{
		Next: root,
	}
	for i := range o.roundTrips {
		transport = RoundTripperChain{
			Do:   o.roundTrips[i],
			Next: transport,
		}
	}

	return &http.Client{
		Transport: transport,
		Timeout:   o.timeout,
	}
}

// NewGetRequestWithContext returns a new http.MethodGet request,
// which is saving your life from http.NewRequestWithContext.
func NewGetRequestWithContext(ctx context.Context, uri string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// NewGetRequest returns a new http.MethodGet request,
// which is saving your life from http.NewRequest.
func NewGetRequest(uri string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, uri, nil)
}

// NewHeadRequestWithContext returns a new http.MethodHead request,
// which is saving your life from http.NewRequestWithContext.
func NewHeadRequestWithContext(ctx context.Context, uri string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// NewHeadRequest returns a new http.MethodHead request,
// which is saving your life from http.NewRequest.
func NewHeadRequest(uri string) (*http.Request, error) {
	return http.NewRequest(http.MethodHead, uri, nil)
}

// NewPostRequestWithContext returns a new http.MethodPost request with the given context,
// which is saving your life from http.NewRequestWithContext.
func NewPostRequestWithContext(ctx context.Context, uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
}

// NewPostRequest returns a new http.MethodPost request,
// which is saving your life from http.NewRequest.
func NewPostRequest(uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(http.MethodPost, uri, body)
}

// NewPutRequestWithContext returns a new http.MethodPut request with the given context,
// which is saving your life from http.NewRequestWithContext.
func NewPutRequestWithContext(ctx context.Context, uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodPut, uri, body)
}

// NewPutRequest returns a new http.MethodPut request,
// which is saving your life from http.NewRequest.
func NewPutRequest(uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(http.MethodPut, uri, body)
}

// NewPatchRequestWithContext returns a new http.MethodPatch request with the given context,
// which is saving your life from http.NewRequestWithContext.
func NewPatchRequestWithContext(ctx context.Context, uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodPatch, uri, body)
}

// NewPatchRequest returns a new http.MethodPatch request,
// which is saving your life from http.NewRequest.
func NewPatchRequest(uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(http.MethodPatch, uri, body)
}

// NewDeleteRequestWithContext returns a new http.MethodDelete request with the given context,
// which is saving your life from http.NewRequestWithContext.
func NewDeleteRequestWithContext(ctx context.Context, uri string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// NewDeleteRequest returns a new http.MethodDelete request,
// which is saving your life from http.NewRequest.
func NewDeleteRequest(uri string) (*http.Request, error) {
	return http.NewRequest(http.MethodDelete, uri, nil)
}

// NewConnectRequestWithContext returns a new http.MethodConnect request with the given context,
// which is saving your life from http.NewRequestWithContext.
func NewConnectRequestWithContext(ctx context.Context, uri string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodConnect, uri, nil)
}

// NewConnectRequest returns a new http.MethodConnect request,
// which is saving your life from http.NewRequest.
func NewConnectRequest(uri string) (*http.Request, error) {
	return http.NewRequest(http.MethodConnect, uri, nil)
}

// NewOptionsRequestWithContext returns a new http.MethodOptions request with the given context,
// which is saving your life from http.NewRequestWithContext.
func NewOptionsRequestWithContext(ctx context.Context, uri string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodOptions, uri, nil)
}

// NewOptionsRequest returns a new http.MethodOptions request,
// which is saving your life from http.NewRequest.
func NewOptionsRequest(uri string) (*http.Request, error) {
	return http.NewRequest(http.MethodOptions, uri, nil)
}

// NewTraceRequestWithContext returns a new http.MethodTrace request with the given context,
// which is saving your life from http.NewRequestWithContext.
func NewTraceRequestWithContext(ctx context.Context, uri string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodTrace, uri, nil)
}

// NewTraceRequest returns a new http.MethodTrace request,
// which is saving your life from http.NewRequest.
func NewTraceRequest(uri string) (*http.Request, error) {
	return http.NewRequest(http.MethodTrace, uri, nil)
}

// Error is similar to http.Error,
// but it can get the error message by the given code.
func Error(rw http.ResponseWriter, code int) {
	http.Error(rw, http.StatusText(code), code)
}

// Close closes the http response body without error.
func Close(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
}

// BodyBytes returns the body of the http response as a byte slice.
func BodyBytes(resp *http.Response) []byte {
	buf := bytespool.GetBytes()
	defer bytespool.Put(buf)

	w := bytespool.GetBuffer()
	_, _ = io.CopyBuffer(w, resp.Body, buf)
	return w.Bytes()
}

// BodyString returns the body of the http response as a string.
func BodyString(resp *http.Response) string {
	return string(BodyBytes(resp))
}
