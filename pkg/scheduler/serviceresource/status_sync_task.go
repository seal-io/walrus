package serviceresource

import (
	"context"
	"fmt"

	"go.uber.org/multierr"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/service"
	"github.com/seal-io/walrus/pkg/dao/model/serviceresource"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/serviceresources"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

type StatusSyncTask struct {
	modelClient model.ClientSet
	logger      log.Logger
}

func NewStatusSyncTask(logger log.Logger, mc model.ClientSet) (in *StatusSyncTask, err error) {
	in = &StatusSyncTask{
		logger:      logger,
		modelClient: mc,
	}

	return
}

func (in *StatusSyncTask) Process(ctx context.Context, args ...any) error {
	// Retrieve operators.
	opIndexer, opLimiter, err := retrieveOperators(ctx, in.modelClient, in.logger)
	if err != nil || len(opIndexer) == 0 {
		return err
	}

	// Count the total size of the services,
	// skip if no services or error raising.
	total, err := in.modelClient.Services().Query().
		Count(ctx)
	if err != nil || total == 0 {
		return err
	}

	wg := gopool.Group()

	// Divide the services in multiple batches according to the gopool burst size.
	bs, bc := getBatches(total, gopool.Burst(), 10)
	// Process the resources in batches concurrently.
	for b := 0; b < bc; b++ {
		p := in.buildProcess(ctx, opIndexer, opLimiter, b*bs, bs)
		wg.Go(p)
	}

	return wg.Wait()
}

func (in *StatusSyncTask) buildProcess(
	ctx context.Context,
	opIndexer operatorIndexer,
	opLimiter operatorLimiter,
	offset,
	limit int,
) func() error {
	return func() error {
		svcs, err := in.modelClient.Services().Query().
			Order(model.Desc(service.FieldCreateTime)).
			Offset(offset).
			Limit(limit).
			Unique(false).
			Select(
				service.FieldID,
				service.FieldStatus).
			All(ctx)
		if err != nil {
			return fmt.Errorf("error listing services with offset %d, limit %d: %w",
				offset, limit, err)
		}

		// Merge the errors to return them all at once,
		// instead of returning the first error.
		var berr error

		for i := range svcs {
			berr = multierr.Append(berr, in.process(ctx, svcs[i], opIndexer, opLimiter))
		}

		return berr
	}
}

func (in *StatusSyncTask) process(
	ctx context.Context,
	svc *model.Service,
	opIndexer operatorIndexer,
	opLimiter operatorLimiter,
) error {
	rs, err := svc.QueryResources().
		Where(
			serviceresource.Shape(types.ServiceResourceShapeInstance),
			serviceresource.ModeNEQ(types.ServiceResourceModeData)).
		Order(model.Desc(serviceresource.FieldCreateTime)).
		Unique(false).
		Select(
			serviceresource.FieldServiceID,
			serviceresource.FieldConnectorID,
			serviceresource.FieldShape,
			serviceresource.FieldMode,
			serviceresource.FieldStatus,
			serviceresource.FieldID,
			serviceresource.FieldDeployerType,
			serviceresource.FieldType,
			serviceresource.FieldName).
		All(ctx)
	if err != nil {
		return fmt.Errorf("error listing service resources: %w", err)
	}

	// Group resources by connector ID.
	connRess := make(map[object.ID][]*model.ServiceResource)
	for y := range rs {
		connRess[rs[y].ConnectorID] = append(connRess[rs[y].ConnectorID], rs[y])
	}

	var sr serviceresources.StateResult

	// Merge the errors to return them all at once,
	// instead of returning the first error.
	var berr error

	for cid, crs := range connRess {
		op, exist := opIndexer[cid]
		if !exist {
			// Ignore if not found operator.
			continue
		}

		// Controls the concurrency of operators with the same ID,
		// avoids server instability or throttling caused by creating too many client connections.
		opLimiter.Acquire(op.ID())

		nsr, err := serviceresources.State(ctx, op, in.modelClient, crs)
		if multierr.AppendInto(&berr, err) {
			// Mark error as transitioning,
			// which doesn't flip the status.
			nsr.Transitioning = true
		}

		sr.Merge(nsr)

		opLimiter.Release(op.ID())
	}

	// State service.
	if status.ServiceStatusDeleted.Exist(svc) {
		// Skip if the service is on deleting.
		return berr
	}

	switch {
	case sr.Error:
		status.ServiceStatusReady.False(svc, "")
	case sr.Transitioning:
		status.ServiceStatusReady.Unknown(svc, "")
	default:
		status.ServiceStatusReady.True(svc, "")
	}

	svc.Status.SetSummary(status.WalkService(&svc.Status))

	err = in.modelClient.Services().UpdateOne(svc).
		SetStatus(svc.Status).
		Exec(ctx)
	if err != nil {
		berr = multierr.Append(berr, err)
	}

	return berr
}
