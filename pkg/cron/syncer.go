package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	settingbus "github.com/seal-io/walrus/pkg/bus/setting"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	cronutil "github.com/seal-io/walrus/utils/cron"
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

		ce, err := cronutil.ParseCronExpr(cs, false)
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
