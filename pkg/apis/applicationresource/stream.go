package applicationresource

import (
	"sync"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/utils/json"
)

func asTermStream(proxy runtime.RequestStream, initWidth, initHeight int32) termStream {
	var resizeCh = make(chan termSize, 2)
	resizeCh <- termSize{Width: initWidth, Height: initHeight}
	return termStream{
		once:   &sync.Once{},
		proxy:  proxy,
		resize: resizeCh,
	}
}

type termSize struct {
	Width  int32 `json:"width"`
	Height int32 `json:"height"`
}

type termStream struct {
	once   *sync.Once
	proxy  runtime.RequestStream
	resize chan termSize
}

func (h termStream) Close() error {
	close(h.resize)
	return nil
}

func (h termStream) Read(p []byte) (n int, err error) {
	for {
		n, err = h.proxy.Read(p)
		if err != nil {
			if runtime.IsRequestStreamCloseError(err) {
				// send exit to upstream if proxy exit unexpectedly.
				h.once.Do(func() {
					n = copy(p, "exit 0\n")
					err = nil
				})
			}
			return
		}
		// resize command is something like `#{"Width":100,"Height":100}#`.
		if n >= 24 && p[0] == '#' && p[1] == '{' && p[n-3] == '}' && p[n-2] == '#' {
			var ts termSize
			if err = json.Unmarshal(p[1:n-2], &ts); err == nil && ts.Width > 0 && ts.Height > 0 {
				h.resize <- ts
			}
			continue
		}
		return
	}
}

func (h termStream) Write(p []byte) (n int, err error) {
	return h.proxy.Write(p)
}

func (h termStream) Next() (uint16, uint16, bool) {
	var t, ok = <-h.resize
	return uint16(t.Width), uint16(t.Height), ok
}
