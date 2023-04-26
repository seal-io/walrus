package server

import (
	"context"

	"github.com/seal-io/seal/pkg/bus"
)

func (r *Server) initSubscribers(ctx context.Context, opts initOptions) (err error) {
	err = bus.Setup(ctx, bus.SetupOptions{ModelClient: opts.ModelClient})
	if err != nil {
		return err
	}

	return
}
