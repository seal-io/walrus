package runtime

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

func isStreamRequest(c *gin.Context) bool {
	return IsUnidiStreamRequest(c) || IsBidiStreamRequest(c)
}

// IsUnidiStreamRequest returns true if the incoming request is a watching request.
func IsUnidiStreamRequest(c *gin.Context) bool {
	return c.Request.Method == http.MethodGet &&
		strings.EqualFold(c.Query("watch"), "true")
}

// IsBidiStreamRequest returns true if the incoming request is a websocket request.
func IsBidiStreamRequest(c *gin.Context) bool {
	return c.Request.Method == http.MethodGet &&
		c.IsWebsocket()
}

func doStreamRequest(c *gin.Context, mr, ri reflect.Value) {
	switch {
	case IsUnidiStreamRequest(c):
		doUnidiStreamRequest(c, mr, ri)
	case IsBidiStreamRequest(c):
		doBidiStreamRequest(c, mr, ri)
	default:
		// Unreachable.
		panic("cannot process as stream request")
	}
}

func doUnidiStreamRequest(c *gin.Context, mr, ri reflect.Value) {
	protoMajor, protoMinor := c.Request.ProtoMajor, c.Request.ProtoMinor
	if protoMajor == 1 && protoMinor == 0 {
		// Do not support http/1.0.
		c.AbortWithStatus(http.StatusUpgradeRequired)
		return
	}

	logger := log.WithName("api")

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()
	proxy := RequestUnidiStream{
		ctx:       ctx,
		ctxCancel: cancel,
		conn:      c.Writer,
	}

	c.Header("Cache-Control", "no-store")
	c.Header("Content-Type", "application/octet-stream; charset=ISO-8859-1")
	c.Header("X-Content-Type-Options", "nosniff")

	if protoMajor == 1 {
		c.Header("Transfer-Encoding", "chunked")
	}

	c.Writer.Flush()

	inputs := make([]reflect.Value, 0, 2)
	inputs = append(inputs, reflect.ValueOf(proxy))
	inputs = append(inputs, ri)
	outputs := mr.Call(inputs)

	errInterface := outputs[len(outputs)-1].Interface()
	if errInterface != nil {
		err := errInterface.(error)
		if !isUnidiDownstreamCloseError(err) {
			logger.Errorf("error processing unidirectional stream request: %v", err)
		}
	}
}

func doBidiStreamRequest(c *gin.Context, mr, ri reflect.Value) {
	logger := log.WithName("api")

	const (
		// Time allowed to read the next pong message from the peer.
		pongWait = 5 * time.Second
		// Send pings to peer with this period, must be less than `pongWait`,
		// it is also the timeout to write a ping message to the peer.
		pingPeriod = (pongWait * 9) / 10
	)

	up := websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
	}

	conn, err := up.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("error upgrading bidirectional stream request: %v", err)
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	var (
		frc = make(chan struct {
			t int
			r io.Reader
			e error
		})
		proxy = RequestBidiStream{
			firstReadOnce: &sync.Once{},
			firstReadChan: frc,
			ctx:           ctx,
			ctxCancel:     cancel,
			conn:          conn,
		}
	)

	// In order to avoid downstream connection leaking,
	// we need configuring a handler to close the upstream context.
	// To trigger the close handler,
	// we have to cut out a goroutine to received downstream,
	// if downstream closes, the close handler will be triggered.
	conn.SetCloseHandler(func(int, string) (err error) {
		cancel()
		return
	})
	gopool.Go(func() {
		var fr struct {
			t int
			r io.Reader
			e error
		}
		fr.t, fr.r, fr.e = conn.NextReader()
		select {
		case frc <- fr:
		case <-ctx.Done():
			close(frc)
		}
	})

	// Ping downstream asynchronously.
	gopool.Go(func() {
		ping := func() error {
			_ = conn.SetReadDeadline(getDeadline(pongWait))
			conn.SetPongHandler(func(string) error {
				return conn.SetReadDeadline(getDeadline(pongWait))
			})

			return conn.WriteControl(websocket.PingMessage,
				[]byte{},
				getDeadline(pingPeriod))
		}

		t := time.NewTicker(pingPeriod)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				if ping() != nil {
					// Cancel upstream if failed to touch downstream.
					cancel()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	})

	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closed")
	inputs := make([]reflect.Value, 0, 2)
	inputs = append(inputs, reflect.ValueOf(proxy))
	inputs = append(inputs, ri)
	outputs := mr.Call(inputs)

	errInterface := outputs[len(outputs)-1].Interface()
	if errInterface != nil {
		err = errInterface.(error)
		if !isBidiDownstreamCloseError(err) {
			var we *websocket.CloseError
			if errors.As(err, &we) {
				closeMsg = websocket.FormatCloseMessage(we.Code, we.Text)
			} else {
				logger.Errorf("error processing bidirectional stream request: %v", err)
				if ue := errors.Unwrap(err); ue != nil {
					err = ue
				}
				closeMsg = websocket.FormatCloseMessage(
					websocket.CloseInternalServerErr, err.Error())
			}
		}
	}
	_ = conn.WriteControl(websocket.CloseMessage, closeMsg, getDeadline(pingPeriod))
}

func isUnidiDownstreamCloseError(err error) bool {
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	errMsg := err.Error()

	return strings.Contains(errMsg, "client disconnected") ||
		strings.Contains(errMsg, "stream closed")
}

func isBidiDownstreamCloseError(err error) bool {
	return errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) ||
		websocket.IsCloseError(err,
			websocket.CloseAbnormalClosure,
			websocket.CloseProtocolError,
			websocket.CloseGoingAway)
}

func getDeadline(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}
