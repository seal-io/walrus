package serviceresource

import (
	"context"
	"io"
	"sync"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/utils/json"
)

func asTermStream(proxy *runtime.RequestBidiStream, initWidth, initHeight int32) termStream {
	resizeCh := make(chan termSize, 2)
	resizeCh <- termSize{Width: initWidth, Height: initHeight}

	return termStream{
		Context: context.Background(),
		once:    &sync.Once{},
		proxy:   proxy,
		resize:  resizeCh,
	}
}

type termSize struct {
	Width  int32 `json:"width"`
	Height int32 `json:"height"`
}

type termStream struct {
	context.Context

	once   *sync.Once
	proxy  *runtime.RequestBidiStream
	resize chan termSize
}

func (h termStream) Read(p []byte) (n int, err error) {
	for {
		n, err = h.proxy.Read(p)
		if err != nil {
			// Send exit to upstream in case of session leaking.
			h.once.Do(func() {
				n = copy(p, "exit 0\n")
				err = nil
			})

			// Prevent useless error message being reported,
			// for example, `use of closed network connection`, `websocket close`.
			if err != nil {
				err = io.EOF
			}

			return
		}

		// Resize command is something like `#{"width":100,"height":100}#`
		// without ending \n(line feed) and \r(carriage return) chars.
		if n >= 24 && p[0] == '#' && p[1] == '{' && p[n-2] == '}' && p[n-1] == '#' {
			var ts termSize
			if err = json.Unmarshal(p[1:n-1], &ts); err == nil && ts.Width > 0 && ts.Height > 0 {
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
	select {
	case <-h.proxy.Done():
		return 0, 0, false
	case t, ok := <-h.resize:
		return uint16(t.Width), uint16(t.Height), ok
	}
}
