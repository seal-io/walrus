package cron

import (
	"context"
	"fmt"

	"go.uber.org/multierr"

	settingbus "github.com/seal-io/walrus/pkg/bus/setting"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/cron"
	"github.com/seal-io/walrus/utils/log"
)

// Sync observes the cron expr setting changes,
// and replaces the jobs that have changed expression.
func Sync(ctx context.Context, m settingbus.BusMessage) error {
	logger := log.WithName("task")

	// Merge the errors to return them all at once,
	// instead of returning the first error.
	var berr error

	for i := 0; i < len(m.Refers); i++ {
		if m.Refers[i] == nil || m.Refers[i].Name == "" {
			continue
		}

		n := m.Refers[i].Name

		c, exist := js[n]
		if !exist {
			continue
		}

		s := settings.Index(n)
		if s == nil {
			continue
		}

		// Get cron expression value of the job from transactional model client.
		v, err := s.Value(ctx, m.TransactionalModelClient)
		if err != nil {
			berr = multierr.Append(berr, fmt.Errorf("error getting job cron expression: %w", err))
			continue
		}

		// Create a job with the new cron expression value.
		expr, job, err := c(logger.WithValues("createdBy", n), v)
		if err != nil {
			berr = multierr.Append(berr, fmt.Errorf("error creating job for %s: %w", n, err))
			continue
		}

		// Reschedule the job with the new cron expression.
		err = cron.Schedule(n, expr, job)
		if err != nil {
			berr = multierr.Append(berr, fmt.Errorf("error rescheduling job for %s: %w", n, err))
			continue
		}

		logger.Infof("rescheduled job %q with expression %q", n, expr)
	}

	return berr
}
