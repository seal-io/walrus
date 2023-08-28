package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	settingbus "github.com/seal-io/walrus/pkg/bus/setting"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/cron"
	"github.com/seal-io/walrus/utils/log"
)

type StartSyncerOptions struct {
	ModelClient model.ClientSet
	Interval    time.Duration
}

func SetupSyncer(ctx context.Context, opts StartSyncerOptions) error {
	syncer := NewSyncer(opts.ModelClient)

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := syncer.Sync(ctx)
			if err != nil {
				syncer.logger.Warnf("error sync cronjob spec: %v", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func NewSyncer(mc model.ClientSet) *Syncer {
	in := &Syncer{}
	in.modelClient = mc
	in.logger = log.WithName("task").WithName("cronjob-spec-sync")

	return in
}

type Syncer struct {
	mu sync.Mutex

	logger      log.Logger
	modelClient model.ClientSet
}

func (in *Syncer) Sync(ctx context.Context) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}
	startTs := time.Now()

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	all, err := in.modelClient.Settings().Query().All(ctx)
	if err != nil {
		return err
	}

	jm := make(map[string]crypto.String, len(all))
	for _, v := range all {
		jm[v.Name] = v.Value
	}

	cj := CurrentJobs()
	for name, j := range cj {
		cst, ok := jm[name]
		if !ok {
			continue
		}
		cs := string(cst)

		ce, err := cron.ParseCronExpr(cs, false)
		if err != nil {
			in.logger.Warnf("cron spec %s is invalid in setting %s: %v", cs, name, err)
			continue
		}

		var (
			lastRun = j.LastRun()
			nextRun = j.NextRun()
		)

		shouldNext := ce.Next(lastRun)

		in.logger.V(6).Infof("cronjob %s last run at %s, next run at %s, calculate next run at %s",
			name, lastRun, nextRun, shouldNext)

		if !lastRun.IsZero() && shouldNext != nextRun {
			in.logger.Infof("cronjob spec %s change to %s, update cronjob", name, cs)

			err = settingbus.Notify(ctx, in.modelClient, model.Settings{
				&model.Setting{
					Name:  name,
					Value: cst,
				},
			})
			if err != nil {
				return fmt.Errorf("error notify cronjob spec %s changed: %w", name, err)
			}
		}
	}

	return nil
}

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
