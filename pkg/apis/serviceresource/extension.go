package serviceresource

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/operator"
	optypes "github.com/seal-io/seal/pkg/operator/types"
)

func (h Handler) RouteGetKeys(req RouteGetKeysRequest) (RouteGetKeysResponse, error) {
	res := req.Entity

	op, err := operator.Get(req.Context, optypes.CreateOptions{
		Connector: *res.Edges.Connector,
	})
	if err != nil {
		return nil, err
	}

	if err = op.IsConnected(req.Context); err != nil {
		return nil, fmt.Errorf("unreachable connector: %w", err)
	}

	return op.GetKeys(req.Context, res)
}

func (h Handler) RouteLog(req RouteLogRequest) error {
	var (
		ctx context.Context
		out io.Writer
	)

	if req.Stream != nil {
		// In stream.
		ctx = req.Stream
		out = req.Stream
	} else {
		ctx = req.Context
		out = req.Context.Writer
	}

	res := req.Entity

	op, err := operator.Get(ctx, optypes.CreateOptions{
		Connector: *res.Edges.Connector,
	})
	if err != nil {
		return err
	}

	if err = op.IsConnected(ctx); err != nil {
		return fmt.Errorf("unreachable connector: %w", err)
	}

	opts := optypes.LogOptions{
		Out:           out,
		WithoutFollow: req.Stream == nil,
		Previous:      req.Previous,
		SinceSeconds:  req.SinceSeconds,
		Timestamps:    req.Timestamps,
		TailLines:     req.TailLines,
	}

	return op.Log(ctx, req.Key, opts)
}

func (h Handler) RouteExec(req RouteExecRequest) error {
	// Only allow stream request.
	if req.Stream == nil {
		return runtime.Errorf(http.StatusBadRequest, "stream request required")
	}

	op, err := operator.Get(req.Stream, optypes.CreateOptions{
		Connector: *req.Entity.Edges.Connector,
	})
	if err != nil {
		return err
	}

	if err = op.IsConnected(req.Stream); err != nil {
		return fmt.Errorf("unreachable connector: %w", err)
	}

	ts := asTermStream(*req.Stream, req.Width, req.Height)
	defer func() { _ = ts.Close() }()

	opts := optypes.ExecOptions{
		Out:     ts,
		In:      ts,
		Shell:   req.Shell,
		Resizer: ts,
	}

	err = op.Exec(ts, req.Key, opts)
	if err != nil {
		if strings.Contains(err.Error(), "OCI runtime exec failed: exec failed:") {
			return &websocket.CloseError{
				Code: websocket.CloseUnsupportedData,
				Text: "unresolved exec shell: " + req.Shell,
			}
		}

		return err
	}

	return nil
}
