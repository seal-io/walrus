package runtime

import (
	"context"
	"errors"
	"reflect"
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
	var logger = log.WithName("restful")

	const (
		// Time allowed to write a message to the peer.
		writeWait = 10 * time.Second
		// Time allowed to read the next pong message from the peer.
		pongWait = 5 * time.Second
		// Send pings to peer with this period, must be less than `pongWait`.
		pingPeriod = (pongWait * 9) / 10
	)

	var up = websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
	}

	var ws, err = up.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("error upgrading stream request: %v", err)
		return
	}
	defer func() {
		_ = ws.Close()
	}()

	var downCtx = c.Request.Context()
	var upCtx, upCtxCancel = context.WithCancel(context.Background())
	defer upCtxCancel()
	var proxy = RequestStream{
		ctx:       upCtx,
		ctxCancel: upCtxCancel,
		ws:        ws,
	}

	gopool.Go(func() {
		var ping = func() error {
			_ = ws.SetReadDeadline(getDeadline(pongWait))
			ws.SetPongHandler(func(string) error {
				return ws.SetReadDeadline(getDeadline(pongWait))
			})
			return ws.WriteControl(websocket.PingMessage,
				[]byte{},
				getDeadline(writeWait))
		}
		var t = time.NewTicker(pingPeriod)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				if ping() != nil {
					// cancel upstream if failed to touch downstream.
					upCtxCancel()
					return
				}
			case <-downCtx.Done():
				// cancel upstream if downstream has been closed explicitly.
				upCtxCancel()
				return
			case <-upCtx.Done():
				// close by main progress.
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
		if !errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded) {
			logger.Errorf("error processing stream request: %v", err)
			var we *websocket.CloseError
			if errors.As(err, &we) {
				closeMsg = websocket.FormatCloseMessage(we.Code, we.Text)
			} else {
				if ue := errors.Unwrap(err); ue != nil {
					err = ue
				}
				closeMsg = websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())
			}
		}
	}
	_ = ws.WriteControl(websocket.CloseMessage, closeMsg, getDeadline(writeWait))
}

func getDeadline(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}
