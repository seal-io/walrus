package telemetry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/telemetry"
	"github.com/seal-io/walrus/utils/log"
)

type PeriodicReportTask struct {
	mu sync.Mutex

	modelClient model.ClientSet
	logger      log.Logger
}

func NewPeriodicReportTask(mc model.ClientSet) (*PeriodicReportTask, error) {
	in := &PeriodicReportTask{}
	in.modelClient = mc
	in.logger = log.WithName("task").WithName(in.Name())

	return in, nil
}

func (in *PeriodicReportTask) Name() string {
	return "telemetry-periodic-report"
}

func (in *PeriodicReportTask) Process(ctx context.Context, args ...any) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}
	startTs := time.Now()

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	err := telemetry.EnqueuePeriodicReportEvent(ctx, in.modelClient)
	if err != nil {
		return fmt.Errorf("error enqueue telemetry periodic report event: %w", err)
	}

	return nil
}
