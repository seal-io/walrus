package runtime

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/seal-io/walrus/utils/json"
)

// RequestUnidiStream holds the request for single-direction stream.
type RequestUnidiStream struct {
	ctx       context.Context
	ctxCancel func()
	conn      gin.ResponseWriter
}

// Write implements io.Writer.
func (r RequestUnidiStream) Write(p []byte) (n int, err error) {
	n, err = r.conn.Write(p)
	if err != nil {
		return
	}

	r.conn.Flush()

	return
}

// SendMsg sends the given data to client.
func (r RequestUnidiStream) SendMsg(data []byte) error {
	_, err := r.Write(data)
	return err
}

// SendJSON marshals the given object as JSON and sends to client.
func (r RequestUnidiStream) SendJSON(i any) error {
	bs, err := json.Marshal(i)
	if err != nil {
		return err
	}

	return r.SendMsg(bs)
}

// Cancel cancels the underlay context.Context.
func (r RequestUnidiStream) Cancel() {
	r.ctxCancel()
}

// Deadline implements context.Context.
func (r RequestUnidiStream) Deadline() (deadline time.Time, ok bool) {
	return r.ctx.Deadline()
}

// Done implements context.Context.
func (r RequestUnidiStream) Done() <-chan struct{} {
	return r.ctx.Done()
}

// Err implements context.Context.
func (r RequestUnidiStream) Err() error {
	return r.ctx.Err()
}

// Value implements context.Context.
func (r RequestUnidiStream) Value(key any) any {
	return r.ctx.Value(key)
}

// RequestBidiStream holds the request for dual-directions stream.
type RequestBidiStream struct {
	firstReadOnce *sync.Once
	firstReadChan <-chan struct {
		t int
		r io.Reader
		e error
	}
	ctx            context.Context
	ctxCancel      func()
	conn           *websocket.Conn
	connReadBytes  *atomic.Int64
	connWriteBytes *atomic.Int64
}

// Read implements io.Reader.
func (r RequestBidiStream) Read(p []byte) (n int, err error) {
	var (
		firstRead bool
		msgType   int
		msgReader io.Reader
	)

	r.firstReadOnce.Do(func() {
		fr, ok := <-r.firstReadChan
		if !ok {
			return
		}
		firstRead = true
		msgType, msgReader, err = fr.t, fr.r, fr.e
	})

	if !firstRead {
		msgType, msgReader, err = r.conn.NextReader()
	}

	if err != nil {
		return
	}

	switch msgType {
	default:
		err = &websocket.CloseError{
			Code: websocket.CloseUnsupportedData,
			Text: "unresolved message type: binary",
		}

		return
	case websocket.TextMessage:
	}

	n, err = msgReader.Read(p)
	if err == nil {
		// Measure read bytes.
		r.connReadBytes.Add(int64(n))
	}

	return
}

// Write implements io.Writer.
func (r RequestBidiStream) Write(p []byte) (n int, err error) {
	msgWriter, err := r.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}

	defer func() { _ = msgWriter.Close() }()

	n, err = msgWriter.Write(p)
	if err == nil {
		// Measure write bytes.
		r.connWriteBytes.Add(int64(n))
	}

	return
}

// RecvMsg receives message from client.
func (r RequestBidiStream) RecvMsg() ([]byte, error) {
	return io.ReadAll(r)
}

// SendMsg sends the given data to client.
func (r RequestBidiStream) SendMsg(data []byte) error {
	_, err := r.Write(data)
	return err
}

// RecvJSON receives JSON message from client and unmarshals into the given object.
func (r RequestBidiStream) RecvJSON(i any) error {
	bs, err := r.RecvMsg()
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, i)
}

// SendJSON marshals the given object as JSON and sends to client.
func (r RequestBidiStream) SendJSON(i any) error {
	bs, err := json.Marshal(i)
	if err != nil {
		return err
	}

	return r.SendMsg(bs)
}

// Cancel cancels the underlay context.Context.
func (r RequestBidiStream) Cancel() {
	r.ctxCancel()
}

// Deadline implements context.Context.
func (r RequestBidiStream) Deadline() (deadline time.Time, ok bool) {
	return r.ctx.Deadline()
}

// Done implements context.Context.
func (r RequestBidiStream) Done() <-chan struct{} {
	return r.ctx.Done()
}

// Err implements context.Context.
func (r RequestBidiStream) Err() error {
	return r.ctx.Err()
}

// Value implements context.Context.
func (r RequestBidiStream) Value(key any) any {
	return r.ctx.Value(key)
}
