package applicationresource

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/applicationresource/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
)

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
	}
}

type Handler struct {
	modelClient model.ClientSet
}

func (h Handler) Kind() string {
	return "ApplicationResource"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

// Batch APIs

// Extensional APIs

func (h Handler) GetKeys(ctx *gin.Context, req view.GetKeysRequest) (*view.GetKeysResponse, error) {
	var res = req.Entity

	var op, err = platform.GetOperator(ctx, operator.CreateOptions{Connector: *res.Edges.Connector})
	if err != nil {
		return nil, err
	}
	ok, err := op.IsConnected(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("cannot connect %s", res.Edges.Connector.Name)
	}

	return op.GetKeys(ctx, *res)
}

func (h Handler) StreamLog(ctx runtime.RequestStream, req view.StreamLogRequest) error {
	var res = req.Entity

	var op, err = platform.GetOperator(ctx, operator.CreateOptions{Connector: *res.Edges.Connector})
	if err != nil {
		return err
	}
	ok, err := op.IsConnected(ctx)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("cannot connect %s", res.Edges.Connector.Name)
	}

	var opts = operator.LogOptions{
		Key:          req.Key,
		Out:          ctx,
		Previous:     req.Previous,
		Tail:         req.Tail,
		SinceSeconds: req.SinceSeconds,
		Timestamps:   req.Timestamps,
	}
	return op.Log(ctx, *res, opts)
}

func (h Handler) StreamExec(ctx runtime.RequestStream, req view.StreamExecRequest) error {
	var res = req.Entity

	op, err := platform.GetOperator(ctx, operator.CreateOptions{Connector: *res.Edges.Connector})
	if err != nil {
		return err
	}
	ok, err := op.IsConnected(ctx)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("cannot connect %s", res.Edges.Connector.Name)
	}

	var ts = asTermStream(ctx, req.Width, req.Height)
	defer func() { _ = ts.Close() }()
	var opts = operator.ExecOptions{
		Key:     req.Key,
		Out:     ts,
		In:      ts,
		Shell:   req.Shell,
		Resizer: ts,
	}
	return op.Exec(ctx, *res, opts)
}
