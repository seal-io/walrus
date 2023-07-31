package server

import (
	"context"

	"github.com/seal-io/seal/pkg/databaselistener"
)

// setupDatabaseListener set the listener to listen database table changes event.
func (r *Server) setupDatabaseListener(ctx context.Context, opts databaselistener.SetupOptions) error {
	return databaselistener.Setup(ctx, opts)
}
