package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"

	settingbus "github.com/seal-io/walrus/pkg/bus/setting"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/setting"
	"github.com/seal-io/walrus/utils/cron"
	"github.com/seal-io/walrus/utils/log"
)

type StartCorrecterOptions struct {
	ModelClient model.ClientSet
	Interval    time.Duration
}

// StartCorrecter starts a corrector to calibrate the jobs.
func StartCorrecter(ctx context.Context, opts StartCorrecterOptions) error {
	c := &corrector{
		logger:      log.WithName("task").WithName("corrector"),
		modelClient: opts.ModelClient,
	}

	return c.Calibrate(ctx, opts.Interval)
}

type corrector struct {
	mu sync.Mutex

	logger      log.Logger
	modelClient model.ClientSet
}

// Calibrate calibrates the all jobs according to the related cron expression periodically.
func (in *corrector) Calibrate(ctx context.Context, interval time.Duration) error {
	return wait.PollUntilContextCancel(ctx, interval, false,
		func(ctx context.Context) (bool, error) {
			// Warn out if error raised,
			// don't stop the loop.
			if err := in.doCalibrate(ctx); err != nil {
				in.logger.Warnf("error calibrating job: %v", err)
			}

			return false, nil
		})
}

func (in *corrector) doCalibrate(ctx context.Context) (err error) {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}

	defer func() {
		in.mu.Unlock()
	}()

	// Entities represents the settings of cron jobs.
	entities, err := in.modelClient.Settings().Query().
		Where(setting.NameIn(sets.KeySet(js).UnsortedList()...)).
		All(ctx)
	if err != nil {
		return err
	}

	// Statues holds the statuses of running cron jobs.
	statuses := cron.State()

	// Candidates holds the settings of cron jobs which need to be updated.
	candidates := make([]*model.Setting, 0, len(entities))

	// Merge the errors to return them all at once,
	// instead of returning the first error.
	var berr error

	for i := range entities {
		name := entities[i].Name

		// NB(thxCode): the job is managed(stop -> restart) by other process,
		// ignore it directly.
		if _, exist := statuses[name]; !exist {
			continue
		}

		expr, err := cron.ParseCronExpr(string(entities[i].Value), false)
		if err != nil {
			berr = multierr.Append(berr,
				fmt.Errorf("invalid cron expression of setting %s: %w", name, err))
			continue
		}

		if !statuses[name].LastRun.IsZero() &&
			!statuses[name].NextRun.Equal(expr.Next(statuses[name].LastRun)) {
			in.logger.Infof("job expression %s changed to %s, update later", name, expr)

			candidates = append(candidates, entities[i])
		}
	}

	if len(candidates) == 0 {
		if berr == nil {
			in.logger.Info("nothing changed")
		}

		return berr
	}

	// Notify syncer to update the cron jobs with settings.
	return multierr.Append(berr, settingbus.Notify(ctx, in.modelClient, candidates))
}
