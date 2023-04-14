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

type LabelApplyTask struct {
	mu sync.Mutex

	modelClient model.ClientSet
	logger      log.Logger
}

func NewLabelApplyTask(mc model.ClientSet) (*LabelApplyTask, error) {
	in := &LabelApplyTask{}
	in.modelClient = mc
	in.logger = log.WithName("task").WithName(in.Name())
	return in, nil
}

func (in *LabelApplyTask) Name() string {
	return "resource-label-apply"
}

func (in *LabelApplyTask) Process(ctx context.Context, args ...interface{}) error {
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
		var st = applicationresources.ApplyLabels(ctx, in.modelClient, 0, bks, nil)
		return st()
	}
	var wg = gopool.Group()
	for bk := 0; bk < bkc; bk++ {
		var st = applicationresources.ApplyLabels(ctx, in.modelClient, bk, bks, nil)
		wg.Go(st)
	}
	return wg.Wait()
}
