package server

import (
	"context"

	"github.com/seal-io/seal/pkg/cron"
)

// setupCronSpecSyncer set the syncer to sync cronjob spec with setting.
func (r *Server) setupCronSpecSyncer(ctx context.Context, opts cron.SetupSyncerOptions) error {
	return cron.SetupSyncer(ctx, opts)
}
