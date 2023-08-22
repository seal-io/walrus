package telemetry

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/telemetry"
	"github.com/seal-io/walrus/utils/log"
)

type PeriodicReportTask struct {
	logger      log.Logger
	modelClient model.ClientSet
}

func NewPeriodicReportTask(logger log.Logger, mc model.ClientSet) (in *PeriodicReportTask, err error) {
	in = &PeriodicReportTask{
		logger:      logger,
		modelClient: mc,
	}

	return
}

func (in *PeriodicReportTask) Name() string {
	return "telemetry-periodic-report"
}

func (in *PeriodicReportTask) Process(ctx context.Context, args ...any) error {
	err := telemetry.EnqueuePeriodicReportEvent(ctx, in.modelClient)
	if err != nil {
		return fmt.Errorf("error enqueue telemetry periodic report event: %w", err)
	}

	return nil
}
