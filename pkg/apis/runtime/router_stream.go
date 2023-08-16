package runtime

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

// isStreamRequest returns true if the incoming request is a stream request.
func isStreamRequest(c *gin.Context) bool {
	return IsUnidiStreamRequest(c) || IsBidiStreamRequest(c)
}

// IsUnidiStreamRequest returns true if the incoming request is a watching request.
func IsUnidiStreamRequest(c *gin.Context) bool {
	return c.Request.Method == http.MethodGet &&
		strings.EqualFold(c.Query("watch"), "true")
}

// doUnidiStreamRequest handles the unidirectional stream request.
func doUnidiStreamRequest(c *gin.Context, route Route, routeInput reflect.Value) {
	logger := log.WithName("api")

	// Ensure chunked request.
	protoMajor, protoMinor := c.Request.ProtoMajor, c.Request.ProtoMinor
	if protoMajor == 1 && protoMinor == 0 {
		// Do not support http/1.0.
		c.AbortWithStatus(http.StatusUpgradeRequired)
		return
	}

	// Flush response headers.
	c.Header("Cache-Control", "no-store")
	c.Header("Content-Type", "application/octet-stream; charset=UTF-8")
	c.Header("X-Content-Type-Options", "nosniff")

	if protoMajor == 1 {
		c.Header("Transfer-Encoding", "chunked")
	}

	c.Writer.Flush()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Minute)
	defer cancel()

	stream := RequestUnidiStream{
		ctx:       ctx,
		ctxCancel: cancel,
		conn:      c.Writer,
	}

	// Discard the underlay context in avoid of conflict using.
	if route.RequestAttributes.HasAll(RequestWithGinContext) {
		routeInput.Interface().(ginContextAdviceReceiver).SetGinContext(nil)
	}

	// Inject request with stream.
	routeInput.Interface().(unidiStreamAdviceReceiver).SetStream(stream)

	// Handle stream request.
	if route.RequestType.Kind() != reflect.Pointer {
		routeInput = routeInput.Elem()
	}
	routeOutputs := route.GoCaller.Call([]reflect.Value{routeInput})

	// Handle error if found.
	if errObj := routeOutputs[len(routeOutputs)-1].Interface(); errObj != nil {
		err := errObj.(error)
		if !isUnidiDownstreamCloseError(err) {
			logger.Errorf("error processing unidirectional stream request: %v", err)
		}
	}
}

// isUnidiDownstreamCloseError returns true if the error is caused by the downstream closing the connection.
func isUnidiDownstreamCloseError(err error) bool {
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	errMsg := err.Error()

	return strings.Contains(errMsg, "client disconnected") ||
		strings.Contains(errMsg, "stream closed")
}

// IsBidiStreamRequest returns true if the incoming request is a websocket request.
func IsBidiStreamRequest(c *gin.Context) bool {
	return c.Request.Method == http.MethodGet &&
		c.IsWebsocket()
}

// doBidiStreamRequest handles the bidirectional stream request.
func doBidiStreamRequest(c *gin.Context, route Route, routeInput reflect.Value) {
	logger := log.WithName("api")

	const (
		// Time allowed to read the next pong message from the peer.
		pongWait = 5 * time.Second
		// Send pings to peer with this period, must be less than `pongWait`,
		// it is also the timeout to write a ping message to the peer.
		pingPeriod = (pongWait * 9) / 10
	)

	// Ensure websocket request.
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Minute)
	defer cancel()

	// In order to avoid downstream connection leaking,
	// we need configuring a handler to close the upstream context.
	// To trigger the close handler,
	// we have to cut out a goroutine to received downstream,
	// if downstream closes, the close handler will be triggered.
	conn.SetCloseHandler(func(int, string) (err error) {
		cancel()
		return
	})

	frc := make(chan struct {
		t int
		r io.Reader
		e error
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

	stream := RequestBidiStream{
		firstReadOnce:  &sync.Once{},
		firstReadChan:  frc,
		ctx:            ctx,
		ctxCancel:      cancel,
		conn:           conn,
		connReadBytes:  &atomic.Int64{},
		connWriteBytes: &atomic.Int64{},
	}

	defer func() {
		c.Set("request_size", stream.connReadBytes.Load())
		c.Set("response_size", stream.connWriteBytes.Load())
	}()

	// Discard the underlay context in avoid of conflict using.
	if route.RequestAttributes.HasAll(RequestWithGinContext) {
		routeInput.Interface().(ginContextAdviceReceiver).SetGinContext(nil)
	}

	// Inject request with stream.
	routeInput.Interface().(bidiStreamAdviceReceiver).SetStream(stream)

	// Handle stream request.
	if route.RequestType.Kind() != reflect.Pointer {
		routeInput = routeInput.Elem()
	}
	routeOutputs := route.GoCaller.Call([]reflect.Value{routeInput})

	// Handle error if found.
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closed")

	if errObj := routeOutputs[len(routeOutputs)-1].Interface(); errObj != nil {
		err = errObj.(error)
		if !isBidiDownstreamCloseError(err) {
			var we *websocket.CloseError
			if errors.As(err, &we) {
				closeMsg = websocket.FormatCloseMessage(
					we.Code, we.Text)

				c.Set("response_status", we.Code)
			} else {
				logger.Errorf("error processing bidirectional stream request: %v", err)

				if ue := errors.Unwrap(err); ue != nil {
					err = ue
				}
				closeMsg = websocket.FormatCloseMessage(
					websocket.CloseInternalServerErr, err.Error())

				c.Set("response_status", websocket.CloseInternalServerErr)
			}
		}
	}

	_ = conn.WriteControl(websocket.CloseMessage, closeMsg, getDeadline(pingPeriod))
}

// isBidiDownstreamCloseError returns true if the error is caused by the downstream closing the connection.
func isBidiDownstreamCloseError(err error) bool {
	return errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) ||
		websocket.IsCloseError(err,
			websocket.CloseAbnormalClosure,
			websocket.CloseProtocolError,
			websocket.CloseGoingAway)
}

// getDeadline returns a deadline with the given duration.
func getDeadline(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}
