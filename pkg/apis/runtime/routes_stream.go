package runtime

import (
	"context"
	"errors"
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

func isUpgradeStreamRequest(c *gin.Context) bool {
	return websocket.IsWebSocketUpgrade(c.Request)
}

func doUpgradeStreamRequest(c *gin.Context, mr reflect.Value, ri reflect.Value) {
	var logger = log.WithName("apis")

	const (
		// Time allowed to read the next pong message from the peer.
		pongWait = 5 * time.Second
		// Send pings to peer with this period, must be less than `pongWait`,
		// it is also the timeout to write a ping message to the peer.
		pingPeriod = (pongWait * 9) / 10
	)

	var up = websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
	}

	var conn, err = up.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("error upgrading stream request: %v", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	var ctx, cancel = context.WithCancel(c)
	defer cancel()
	var (
		fro sync.Once
		frc = make(chan struct {
			t int
			r io.Reader
			e error
		})
		proxy = RequestStream{
			firstReadOnce: &fro,
			firstReadChan: frc,
			ctx:           ctx,
			ctxCancel:     cancel,
			conn:          conn,
		}
	)

	// in order to avoid downstream connection leaking,
	// we need configuring a handler to close the upstream context.
	// to trigger the close handler,
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
		frc <- fr
	})

	// ping downstream asynchronously.
	gopool.Go(func() {
		var ping = func() error {
			_ = conn.SetReadDeadline(getDeadline(pongWait))
			conn.SetPongHandler(func(string) error {
				return conn.SetReadDeadline(getDeadline(pongWait))
			})
			return conn.WriteControl(websocket.PingMessage,
				[]byte{},
				getDeadline(pingPeriod))
		}
		var t = time.NewTicker(pingPeriod)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				if ping() != nil {
					// cancel upstream if failed to touch downstream.
					cancel()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	})

	var closeMsg = websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closed")
	var inputs = make([]reflect.Value, 0, 2)
	inputs = append(inputs, reflect.ValueOf(proxy))
	inputs = append(inputs, ri)
	var outputs = mr.Call(inputs)
	var errInterface = outputs[len(outputs)-1].Interface()
	if errInterface != nil {
		err = errInterface.(error)
		if !isDownstreamCloseError(err) {
			var we *websocket.CloseError
			if errors.As(err, &we) {
				closeMsg = websocket.FormatCloseMessage(we.Code, we.Text)
			} else {
				logger.Errorf("error processing stream request: %v", err)
				if ue := errors.Unwrap(err); ue != nil {
					err = ue
				}
				closeMsg = websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())
			}
		}
	}
	_ = conn.WriteControl(websocket.CloseMessage, closeMsg, getDeadline(pingPeriod))
}

func isDownstreamCloseError(err error) bool {
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
