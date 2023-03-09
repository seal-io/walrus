package server

import (
	"context"

	"github.com/seal-io/seal/pkg/bus"
	"github.com/seal-io/seal/pkg/topic"
)

func (r *Server) initSubscribers(ctx context.Context, opts initOptions) error {
	if err := bus.Setup(ctx, bus.SetupOptions{ModelClient: opts.ModelClient}); err != nil {
		return err
	}

	if err := topic.Setup(ctx, topic.SetupOptions{ModelClient: opts.ModelClient}); err != nil {
		return err
	}

	return nil
}
