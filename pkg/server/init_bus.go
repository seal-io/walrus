package server

import (
	"context"

	"github.com/seal-io/seal/pkg/bus"
)

// setupBusSubscribers launches the synchronous subscribers provided by the bus.
func (r *Server) setupBusSubscribers(ctx context.Context, opts initOptions) error {
	return bus.Setup(ctx, bus.SetupOptions{
		ModelClient: opts.ModelClient,
	})
}
