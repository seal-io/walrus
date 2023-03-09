package server

import (
	"context"

	"github.com/seal-io/seal/pkg/bus"
)

func (r *Server) initSubscribers(ctx context.Context, opts initOptions) error {
	if err := r.setupBusSubscribers(ctx, opts); err != nil {
		return err
	}

	return nil
}

func (r *Server) setupBusSubscribers(ctx context.Context, opts initOptions) error {
	return bus.Setup(ctx, bus.SetupOptions{
		ModelClient: opts.ModelClient,
	})
}
