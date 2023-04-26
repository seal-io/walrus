package applicationresource

import (
	"context"
	"sync"
	"time"

	"github.com/seal-io/seal/pkg/applicationresources"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type StatusSyncTask struct {
	mu sync.Mutex

	modelClient model.ClientSet
	logger      log.Logger
}

func NewStatusSyncTask(mc model.ClientSet) (*StatusSyncTask, error) {
	in := &StatusSyncTask{}
	in.modelClient = mc
	in.logger = log.WithName("task").WithName(in.Name())
	return in, nil
}

func (in *StatusSyncTask) Name() string {
	return "resource-status-sync"
}

func (in *StatusSyncTask) Process(ctx context.Context, args ...interface{}) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}
	var startTs = time.Now()
	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	var cnt, err = in.modelClient.ApplicationResources().Query().
		Count(ctx)
	if err != nil {
		return err
	}

	// divide processing buckets with count.
	const bks = 100
	var bkc = cnt / bks
	if bkc == 0 {
		var st = in.buildStateTask(ctx, 0, bks)
		return st()
	}
	var wg = gopool.Group()
	for bk := 0; bk < bkc; bk++ {
		var st = in.buildStateTask(ctx, bk, bks)
		wg.Go(st)
	}
	return wg.Wait()
}

func (in *StatusSyncTask) buildStateTask(ctx context.Context, offset, limit int) func() error {
	return func() error {
		var entities, err = applicationresources.ListStateCandidatesByPage(
			ctx, in.modelClient, offset, limit)
		if err != nil {
			return err
		}
		return applicationresources.State(ctx, in.modelClient, entities)
	}
}
