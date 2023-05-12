package req

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"math"
	"mime/multipart"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"

	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/version"
)

// HTTP returns a new http client.
func HTTP() *HttpClient {
	return &HttpClient{
		client: &fasthttp.Client{
			NoDefaultUserAgentHeader: false,
			Name:                     version.GetUserAgent(),
			ReadTimeout:              15 * time.Second,
			MaxConnDuration:          5 * time.Minute,
			MaxConnWaitTimeout:       10 * time.Second,
			MaxIdleConnDuration:      10 * time.Second,
			Dial:                     fasthttpproxy.FasthttpProxyHTTPDialerTimeout(5 * time.Second),
			DisablePathNormalizing:   true,
			// Respect the request retry backoff of HttpRequest.
			MaxIdemponentCallAttempts: 1,
		},
	}
}

var defaultHttpClient = HTTP()

// HTTPRequest returns a new request generated by default http client.
func HTTPRequest() *HttpRequest {
	return defaultHttpClient.Request()
}

type HttpClient struct {
	client *fasthttp.Client
}

// WithIf allows chain calling conditional parameterizing.
func (in *HttpClient) WithIf(fn func(cli *HttpClient)) *HttpClient {
	if fn != nil {
		fn(in)
	}
	return in
}

// WithReadTimeout specifies the maximum duration for full response reading,
// default duration is 15 seconds, unlimited if specifies with 0,
// and only be valid if the duration is less than context.Context deadline.
func (in *HttpClient) WithReadTimeout(duration time.Duration) *HttpClient {
	in.client.ReadTimeout = duration
	return in
}

// WithWriteTimeout specifies the maximum duration for full request writing,
// default duration is unlimited,
// and only be valid if the duration is less than context.Context deadline.
func (in *HttpClient) WithWriteTimeout(duration time.Duration) *HttpClient {
	in.client.WriteTimeout = duration
	return in
}

// WithMaxConnDuration specifies the maximum duration for underlay dialer connecting,
// default duration is 5 minutes, unlimited if specifies with 0.
func (in *HttpClient) WithMaxConnDuration(duration time.Duration) *HttpClient {
	in.client.MaxConnDuration = duration
	return in
}

// WithUserAgent specifies the user-agent for request.
func (in *HttpClient) WithUserAgent(ua string) *HttpClient {
	in.client.Name = ua
	return in
}

// WithDial specifies the dial for connecting,
// default is an environment proxy aware dialer which has 5 seconds timeout.
func (in *HttpClient) WithDial(dial func(addr string) (net.Conn, error)) *HttpClient {
	in.client.Dial = dial
	return in
}

// WithInsecureSkipVerifyEnabled specifies skipping server certificate verification.
func (in *HttpClient) WithInsecureSkipVerifyEnabled() *HttpClient {
	if in.client.TLSConfig == nil {
		in.client.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	in.client.TLSConfig.InsecureSkipVerify = true
	return in
}

// WithMalformedPathNormalizeEnabled specifies normalizing malformed path for request.
func (in *HttpClient) WithMalformedPathNormalizeEnabled() *HttpClient {
	in.client.DisablePathNormalizing = false
	return in
}

// WithDualStackEnabled specifies enabling dual-stack connecting.
func (in *HttpClient) WithDualStackEnabled() *HttpClient {
	in.client.DialDualStack = true
	return in
}

// Request creates new http request.
func (in *HttpClient) Request() *HttpRequest {
	return &HttpRequest{
		client:              in.client,
		request:             fasthttp.AcquireRequest(),
		requestRetryIf:      defaultHttpRequestRetryIf,
		requestRetryBackoff: createHttpRequestRetryBackoff(1*time.Second, 15*time.Second, 5),
	}
}

// HttpRequestRetryIf defines the condition for request retry.
type HttpRequestRetryIf func(statusCode int, respError error) bool

type httpRequestRetryBackoff func(attemptNum int, resp *fasthttp.Response) (time.Duration, bool)

type HttpRequest struct {
	err                 error
	client              *fasthttp.Client
	request             *fasthttp.Request
	requestRetryIf      HttpRequestRetryIf
	requestRetryBackoff httpRequestRetryBackoff
	redirect            bool
	insight             bool
}

// WithIf allows chain calling conditional parameterizing.
func (in *HttpRequest) WithIf(fn func(cli *HttpRequest)) *HttpRequest {
	if fn != nil {
		fn(in)
	}
	return in
}

// WithBody specifies the given reader as body for request, which leads the header with
// - Content-Type: application/octet-stream; charset=ISO-8859-1
// - Transfer-Encoding: chunked,
// it's able to customize with WithContentType, WithContentEncoding and WithContentDisposition.
func (in *HttpRequest) WithBody(r io.Reader) *HttpRequest {
	in.request.SetBodyStream(r, -1)
	in.request.Header.SetContentType("application/octet-stream; charset=ISO-8859-1")
	return in
}

// WithBodyString specifies the given string as body for request, which leads the header with
// - Content-Type: text/plain; charset=UTF-8
// - Content-Length: size of bytes (calculated by fasthttp),
// it's able to overwrite the Content-Type with WithContentType.
func (in *HttpRequest) WithBodyString(s string) *HttpRequest {
	in.request.SetBodyString(s)
	in.request.Header.SetContentType("text/plain; charset=UTF-8")
	return in
}

// WithBodyBytes specifies the given bytes as body for request, which leads the header with
// - Content-Type: application/octet-stream; charset=ISO-8859-1
// - Content-Length: size of bytes (calculated by fasthttp),
// it's able to customize with WithContentType, WithContentEncoding and WithContentDisposition.
func (in *HttpRequest) WithBodyBytes(bs []byte) *HttpRequest {
	in.request.SetBody(bs)
	in.request.Header.SetContentType("application/octet-stream; charset=ISO-8859-1")
	return in
}

// WithBodyJSON parses the given object as JSON bytes,
// and specifies the bytes as body for request, which leads the header with
// - Content-Type: application/json; charset=UTF-8
// - Content-Length: size of JSON bytes (calculated by fasthttp),
// it's able to overwrite the Content-Type with WithContentType.
func (in *HttpRequest) WithBodyJSON(object interface{}) *HttpRequest {
	var bs, err = json.Marshal(object)
	if err != nil {
		in.err = err
		return in
	}
	in.request.SetBody(bs)
	in.request.Header.SetContentType("application/json; charset=UTF-8")
	return in
}

// WithBodyForm parses the given string-string map as multiple parts data form,
// and specifies the form as body for request, which leads the header with
// - Content-Type: multipart/form-data; boundary=...
// - Content-Length: size of bytes (calculated by fasthttp),
// it supports to upload a local file with "@" prefix.
func (in *HttpRequest) WithBodyForm(formParams url.Values) *HttpRequest {
	var buff bytes.Buffer
	var w = multipart.NewWriter(&buff)

	var err error
	for k, v := range formParams {
		for _, iv := range v {
			if strings.HasPrefix(k, "@") { // File.
				err = addFile(w, k[1:], iv)
			} else { // Form value.
				err = w.WriteField(k, iv)
			}
			if err != nil {
				in.err = err
				return in
			}
		}
	}
	err = w.Close()
	if err != nil {
		in.err = err
		return in
	}

	in.request.SetBody(buff.Bytes())
	in.request.Header.SetContentType(w.FormDataContentType())
	return in
}

func addFile(w *multipart.Writer, fieldName, path string) error {
	var file, err = os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	part, err := w.CreateFormFile(fieldName, filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	return err
}

// HttpCookie alias the fasthttp.Cookie.
type HttpCookie = fasthttp.Cookie

// WithCookies specifies the cookies for request.
func (in *HttpRequest) WithCookies(cs ...*HttpCookie) *HttpRequest {
	for i := 0; i < len(cs); i++ {
		var c fasthttp.Cookie
		c.CopyTo(cs[i])
		c.SetKeyBytes([]byte{}) // Clean key to make a correct bytes.
		in.request.Header.SetCookieBytesKV(cs[i].Key(), c.Cookie())
	}
	return in
}

// WithAccept specifies the Accept header for request.
func (in *HttpRequest) WithAccept(a string) *HttpRequest {
	in.request.Header.Set(fasthttp.HeaderAccept, a)
	return in
}

// WithAcceptJSON specifies the Accept header as "Accept: application/json" for request.
func (in *HttpRequest) WithAcceptJSON() *HttpRequest {
	in.WithAccept("application/json")
	return in
}

// WithContentType specifies the Content-Type header for request.
func (in *HttpRequest) WithContentType(ct string) *HttpRequest {
	in.request.Header.SetContentType(ct)
	return in
}

// WithContentEncoding specifies the Content-Encoding header for request.
func (in *HttpRequest) WithContentEncoding(ce string) *HttpRequest {
	in.request.Header.SetContentEncoding(ce)
	return in
}

// WithContentDisposition specifies the Content-Disposition header for request.
func (in *HttpRequest) WithContentDisposition(cd string) *HttpRequest {
	in.request.Header.Set(fasthttp.HeaderContentDisposition, cd)
	return in
}

// WithConnReuseSkipEnabled specifies skipping reuse the underlay connection after requesting.
func (in *HttpRequest) WithConnReuseSkipEnabled() *HttpRequest {
	in.request.SetConnectionClose()
	return in
}

// WithHeaders specifies the headers for request.
func (in *HttpRequest) WithHeaders(hs map[string]string) *HttpRequest {
	for k, v := range hs {
		in.request.Header.Set(k, v)
	}
	return in
}

// WithHeader specifies the header for request.
func (in *HttpRequest) WithHeader(k, v string) *HttpRequest {
	in.request.Header.Set(k, v)
	return in
}

// WithBearerAuthToken specifies the bearer auth token for request.
func (in *HttpRequest) WithBearerAuthToken(t string) *HttpRequest {
	return in.WithHeader("Authorization", "Bearer "+t)
}

// WithBasicAuth specifies the basic auth for request.
func (in *HttpRequest) WithBasicAuth(u, p string) *HttpRequest {
	var t = base64.StdEncoding.EncodeToString([]byte(u + ":" + p))
	return in.WithHeader("Authorization", "Basic "+t)
}

// WithRetryBackoff specifies the retry-backoff mechanism for request,
// default retry 5 times within 1s, 2s, 4s, 8s, 15s waiting,
// if retry-if == true, the whole time-cost is below,
//   - failed to dial server, default dialing timeout is 5s
//     5s -> [1s] -> 5s -> [2s] -> 5s -> [4s] -> 5s -> [8s] -> 5s -> [15s] -> 5s == 30s+[30s]
//   - dial server success, but failed to read, default read timeout is 15s
//     0s+15s -> [1s] -> 0s+15s -> [2s] -> 0s+15s -> [4s] -> 0s+15s -> [8s] -> 0s+15s -> [15s] -> 0s+15s == 90s+[30s].
func (in *HttpRequest) WithRetryBackoff(waitMin, waitMax time.Duration, attemptMax int) *HttpRequest {
	in.requestRetryBackoff = createHttpRequestRetryBackoff(waitMin, waitMax, attemptMax)
	return in
}

// WithRetryIf specifies the if-condition of retry operation for request,
// or stops retrying if setting with `nil`, default retry rules is below,
// - receiving TLS handshake timeout, the duration is same as HttpClient.WithWriteTimeout
// - receiving dial timeout, the duration is same as HttpClient.WithDial
// - receiving connection closed
// - receiving hang-on
// - receiving rate-limited of server, 429
// - receiving unexpected status code, not in the range of (1, 499).
func (in *HttpRequest) WithRetryIf(retryIf HttpRequestRetryIf) *HttpRequest {
	in.requestRetryIf = retryIf
	return in
}

// WithRedirect specifies to redirect automatically if receiving any redirect-able response.
func (in *HttpRequest) WithRedirect() *HttpRequest {
	in.redirect = true
	return in
}

// WithInsight specifies to insight into the request.
func (in *HttpRequest) WithInsight() *HttpRequest {
	in.insight = true
	return in
}

// Response do actual requesting.
func (in *HttpRequest) Response(ctx context.Context, url string, method string) *HttpResponse {
	var request = fasthttp.AcquireRequest()
	in.request.CopyTo(request)
	request.SetRequestURI(url)
	request.Header.SetMethod(method)

	var response = fasthttp.AcquireResponse()
	var resp = &HttpResponse{
		err:      in.err,
		request:  request,
		response: response,
	}
	if in.err != nil {
		return resp
	}
	for i := 0; ; i++ {
		if ctx.Err() != nil {
			// Allow faster failure.
			resp.err = ctx.Err()
			break
		}
		if in.insight {
			log.Debugf("requesting %s", request.Header.String())
		}
		var respErrChan = make(chan error)
		gopool.Go(func() {
			if in.redirect {
				respErrChan <- in.client.DoRedirects(request, response, 16)
				return
			}
			if deadline, ok := ctx.Deadline(); ok {
				respErrChan <- in.client.DoDeadline(request, response, deadline)
			} else {
				respErrChan <- in.client.Do(request, response)
			}
		})
		select {
		case <-ctx.Done():
			// Allow cancelling from requesting.
			resp.err = ctx.Err()
			break
		case err := <-respErrChan:
			resp.err = err
		}
		if in.requestRetryIf == nil || !in.requestRetryIf(response.StatusCode(), resp.err) {
			break
		}
		var waitDuration, shouldWait = in.requestRetryBackoff(i+1, response)
		if !shouldWait {
			log.Warnf("reached limitation of retry requesting %s %s", method, url)
			break
		}
		log.Debugf("retry requesting %s %s after %v", method, url, waitDuration)
		var waitTimer = time.NewTimer(waitDuration)
		select {
		case <-ctx.Done():
			// Allow cancelling from retry waiting.
			waitTimer.Stop()
			resp.err = ctx.Err()
			break
		case <-waitTimer.C:
		}
	}
	return resp
}

func (in *HttpRequest) PostWithContext(ctx context.Context, url string) *HttpResponse {
	return in.Response(ctx, url, fasthttp.MethodPost)
}

func (in *HttpRequest) Post(url string) *HttpResponse {
	return in.PostWithContext(context.Background(), url)
}

func (in *HttpRequest) DeleteWithContext(ctx context.Context, url string) *HttpResponse {
	return in.Response(ctx, url, fasthttp.MethodDelete)
}

func (in *HttpRequest) Delete(url string) *HttpResponse {
	return in.DeleteWithContext(context.Background(), url)
}

func (in *HttpRequest) PutWithContext(ctx context.Context, url string) *HttpResponse {
	return in.Response(ctx, url, fasthttp.MethodPut)
}

func (in *HttpRequest) Put(url string) *HttpResponse {
	return in.PutWithContext(context.Background(), url)
}

func (in *HttpRequest) PatchWithContext(ctx context.Context, url string) *HttpResponse {
	return in.Response(ctx, url, fasthttp.MethodPatch)
}

func (in *HttpRequest) Patch(url string) *HttpResponse {
	return in.PatchWithContext(context.Background(), url)
}

func (in *HttpRequest) GetWithContext(ctx context.Context, url string) *HttpResponse {
	return in.Response(ctx, url, fasthttp.MethodGet)
}

func (in *HttpRequest) Get(url string) *HttpResponse {
	return in.GetWithContext(context.Background(), url)
}

func (in *HttpRequest) HeadWithContext(ctx context.Context, url string) *HttpResponse {
	return in.Response(ctx, url, fasthttp.MethodHead)
}

func (in *HttpRequest) Head(url string) *HttpResponse {
	return in.HeadWithContext(context.Background(), url)
}

func (in *HttpRequest) OptionsWithContext(ctx context.Context, url string) *HttpResponse {
	return in.Response(ctx, url, fasthttp.MethodOptions)
}

func (in *HttpRequest) Options(url string) *HttpResponse {
	return in.OptionsWithContext(context.Background(), url)
}

func defaultHttpRequestRetryIf(statusCode int, respErr error) bool {
	if respErr != nil {
		// Retry if receiving TLS handshake timeout.
		if errors.Is(respErr, fasthttp.ErrTLSHandshakeTimeout) {
			return true
		}
		// Retry if receiving dialing timeout.
		if errors.Is(respErr, fasthttp.ErrDialTimeout) {
			return true
		}
		// Retry if receiving connection closed.
		if errors.Is(respErr, fasthttp.ErrConnectionClosed) {
			return true
		}
		var respErrMsg = respErr.Error()
		// Retry if receiving requesting timeout (tcp connected but failed to send header),
		// it can cause by DNS error, server resource exhausted.
		return strings.Contains(respErrMsg, "Client.Timeout exceeded while awaiting headers")
	}
	// Retry if receiving rate-limited of server.
	if statusCode == fasthttp.StatusTooManyRequests {
		return true
	}
	// Retry if receiving unexpected responses.
	if statusCode == 0 || (statusCode >= 500 && statusCode != fasthttp.StatusNotImplemented) {
		return true
	}
	return false
}

func createHttpRequestRetryBackoff(waitMin, waitMax time.Duration, attemptMax int) httpRequestRetryBackoff {
	return func(attemptNum int, resp *fasthttp.Response) (time.Duration, bool) {
		if attemptNum > attemptMax {
			return 0, false
		}
		if resp != nil {
			switch resp.StatusCode() {
			case fasthttp.StatusTooManyRequests, fasthttp.StatusServiceUnavailable:
				if retryAfterBytes := resp.Header.Peek(fasthttp.HeaderRetryAfter); len(retryAfterBytes) != 0 {
					if v, err := strconv.ParseInt(string(retryAfterBytes), 10, 64); err == nil {
						return time.Duration(v) * time.Second, true
					}
				}
			}
		}
		var v = math.Pow(2, float64(attemptNum)) * float64(waitMin)
		var sleep = time.Duration(v)
		if sleep > waitMax {
			sleep = waitMax
		}
		return sleep, true
	}
}

type HttpResponse struct {
	once     sync.Once
	err      error
	request  *fasthttp.Request
	response *fasthttp.Response
}

// Error returns error during chain calling,
// and releases the underlay request and response if error occurring.
func (in *HttpResponse) Error() error {
	in.once.Do(func() {
		if in.err == nil {
			var sc = in.StatusCode()
			if !(199 < sc && sc < 300) {
				var msg strings.Builder
				msg.WriteString(strconv.FormatInt(int64(sc), 10))
				var message = in.StatusMessage()
				if message != "" {
					msg.WriteString(" ")
					msg.WriteString(message)
				}
				var body = in.BodyStringOnly()
				if body != "" {
					msg.WriteString(": ")
					msg.WriteString(body)
				}
				in.err = errors.New(msg.String())
			}
		}
		// Releases the underlay request and response.
		if in.err != nil {
			fasthttp.ReleaseRequest(in.request)
			fasthttp.ReleaseResponse(in.response)
		}
	})
	return in.err
}

// StatusCode returns status code of response,
// it must call before Release called.
func (in *HttpResponse) StatusCode() int {
	return in.response.StatusCode()
}

// StatusMessage returns status message of response,
// it must call before Release called.
func (in *HttpResponse) StatusMessage() string {
	var m = in.response.Header.StatusMessage()
	if len(m) == 0 {
		return ""
	}
	return string(m)
}

// ResponseHeader alias the fasthttp.ResponseHeader.
type ResponseHeader = fasthttp.ResponseHeader

// Headers returns headers of response,
// it must call before Release called.
func (in *HttpResponse) Headers() *ResponseHeader {
	var header ResponseHeader
	in.response.Header.CopyTo(&header)
	return &header
}

// Header returns header of response with the given name,
// it must call before Release called.
func (in *HttpResponse) Header(k string) string {
	return string(in.response.Header.Peek(k))
}

// Cookies returns cookies of response,
// it must call before Release called.
func (in *HttpResponse) Cookies(ks ...string) (cs []*HttpCookie) {
	if len(ks) != 0 {
		for _, k := range ks {
			var cbytes = in.response.Header.PeekCookie(k)
			if len(cbytes) == 0 {
				continue
			}
			var c HttpCookie
			if err := c.ParseBytes(cbytes); err == nil {
				cs = append(cs, &c)
			}
		}
		return
	}
	in.response.Header.VisitAllCookie(func(key, cbytes []byte) {
		var c HttpCookie
		if err := c.ParseBytes(cbytes); err == nil {
			cs = append(cs, &c)
		}
	})
	return
}

// BodyBytesOnly returns the body bytes of response,
// it must call before Release called.
func (in *HttpResponse) BodyBytesOnly() []byte {
	return in.response.Body()
}

// BodyBytes returns the body bytes of response,
// and the error if occurring,
// it must call before Release called.
func (in *HttpResponse) BodyBytes() ([]byte, error) {
	if err := in.Error(); err != nil {
		return nil, err
	}
	return in.BodyBytesOnly(), nil
}

// BodyStringOnly returns the body string of response.
func (in *HttpResponse) BodyStringOnly() string {
	return string(in.response.Body())
}

// BodyString returns the body string of response
// and releases if the error occurs.
func (in *HttpResponse) BodyString() (string, error) {
	if err := in.Error(); err != nil {
		return "", err
	}
	return in.BodyStringOnly(), nil
}

// BodyOnly returns the body reader of response.
func (in *HttpResponse) BodyOnly() io.Reader {
	return bytes.NewBuffer(in.response.Body())
}

// Body returns the body reader of response,
// and releases if the error occurs.
func (in *HttpResponse) Body() (io.Reader, error) {
	if err := in.Error(); err != nil {
		return nil, err
	}
	return in.BodyOnly(), nil
}

// BodyJSON unmarshals the body bytes of response into the given receiver,
// and the error if occurring,
// it must call before Release called.
func (in *HttpResponse) BodyJSON(ptr interface{}) error {
	if err := in.Error(); err != nil {
		return err
	}
	return json.Unmarshal(in.response.Body(), ptr)
}
