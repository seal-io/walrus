package server

import (
	"context"

	"github.com/seal-io/seal/pkg/connectors"
)

func (r *Server) initSubscribes(ctx context.Context, opts initOptions) error {
	return connectors.AddSubscriber("connector-cost-subscriber", connectors.EnsureCostTools)
}
