package runtime

import (
	"context"
	"errors"
	"io"
	"reflect"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/seal-io/seal/utils/bytespool"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

func isUpgradeStreamRequest(c *gin.Context) bool {
	return websocket.IsWebSocketUpgrade(c.Request)
}

func doUpgradeStreamRequest(c *gin.Context, mr reflect.Value) {
	var logger = log.WithName("restful")

	const (
		// Time allowed to write a message to the peer.
		writeWait = 10 * time.Second
		// Time allowed to read the next pong message from the peer.
		pongWait = 60 * time.Second
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
		_ = ws.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			getDeadline(writeWait))
		_ = ws.Close()
	}()

	// prepare request stream
	var g, ctx = gopool.WithContext(c.Request.Context())
	var sc = make(chan io.Reader, runtime.NumCPU()*2)
	var rc = make(chan any, runtime.NumCPU()*2)
	defer close(sc)
	defer close(rc)
	var rs = RequestStream{
		ctx: ctx,
		sc:  sc,
		rc:  rc,
	}

	g.Go(func() error {
		var inputs = make([]reflect.Value, 0, 1)
		inputs = append(inputs, reflect.ValueOf(rs))
		_ = mr.Call(inputs)
		return context.Canceled // escape
	})
	g.Go(func() error {
		var t = time.NewTicker(pingPeriod)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				_ = ws.SetReadDeadline(getDeadline(pongWait))
				ws.SetPongHandler(func(string) error {
					return ws.SetReadDeadline(getDeadline(pongWait))
				})
				var err = ws.WriteControl(websocket.PingMessage,
					[]byte{}, getDeadline(writeWait))
				if err != nil {
					return err
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			var msgType, msgReader, err = ws.NextReader()
			if err != nil {
				return err
			}
			switch msgType {
			default:
				logger.Warnf("received unexpected message: %d", msgType)
				continue
			case websocket.TextMessage:
			}
			var t = time.NewTimer(1 * time.Second)
			select {
			case rc <- msgReader:
				t.Stop()
			case <-t.C:
				t.Stop()
				return errors.New("timeout receiving: blocked buffer")
			}
		}
	})
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			var msgReader = <-sc
			var msgWriter, err = ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return err
			}
			var buff = bytespool.GetBytes(0)
			err = func() error {
				defer func() {
					_ = msgWriter.Close()
					bytespool.Put(buff)
				}()
				var _, err = io.CopyBuffer(msgWriter, msgReader, buff)
				return err
			}()
			if err != nil {
				return err
			}
		}
	})
	err = g.Wait()
	if err != nil {
		rc <- err
		if !errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded) &&
			!websocket.IsCloseError(err, websocket.CloseNoStatusReceived, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			logger.Warnf("failed to keep stream request: %v", err)
		}
	}
}

func getDeadline(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}
