package resourcecomponent

import (
	"context"
	"fmt"

	"go.uber.org/multierr"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/resourcecomponents"
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

	// Count the total size of the resources,
	// skip if no resources or error raising.
	total, err := in.modelClient.Resources().Query().
		Count(ctx)
	if err != nil || total == 0 {
		return err
	}

	wg := gopool.Group()

	// Divide the resources in multiple batches according to the gopool burst size.
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
		ress, err := in.modelClient.Resources().Query().
			Order(model.Desc(resource.FieldCreateTime)).
			Offset(offset).
			Limit(limit).
			Unique(false).
			Select(
				resource.FieldID,
				resource.FieldStatus).
			All(ctx)
		if err != nil {
			return fmt.Errorf("error listing resources with offset %d, limit %d: %w",
				offset, limit, err)
		}

		// Merge the errors to return them all at once,
		// instead of returning the first error.
		var berr error

		for i := range ress {
			berr = multierr.Append(berr, in.process(ctx, ress[i], opIndexer, opLimiter))
		}

		return berr
	}
}

func (in *StatusSyncTask) process(
	ctx context.Context,
	res *model.Resource,
	opIndexer operatorIndexer,
	opLimiter operatorLimiter,
) error {
	rs, err := res.QueryComponents().
		Where(
			resourcecomponent.Shape(types.ResourceComponentShapeInstance),
			resourcecomponent.ModeNEQ(types.ResourceComponentModeData)).
		Order(model.Desc(resourcecomponent.FieldCreateTime)).
		Unique(false).
		Select(
			resourcecomponent.FieldResourceID,
			resourcecomponent.FieldConnectorID,
			resourcecomponent.FieldShape,
			resourcecomponent.FieldMode,
			resourcecomponent.FieldStatus,
			resourcecomponent.FieldID,
			resourcecomponent.FieldDeployerType,
			resourcecomponent.FieldType,
			resourcecomponent.FieldName).
		All(ctx)
	if err != nil {
		return fmt.Errorf("error listing resource components: %w", err)
	}

	// Group resources by connector ID.
	connRess := make(map[object.ID][]*model.ResourceComponent)
	for y := range rs {
		connRess[rs[y].ConnectorID] = append(connRess[rs[y].ConnectorID], rs[y])
	}

	var sr resourcecomponents.StateResult

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

		nsr, err := resourcecomponents.State(ctx, op, in.modelClient, crs)
		if multierr.AppendInto(&berr, err) {
			// Mark error as transitioning,
			// which doesn't flip the status.
			nsr.Transitioning = true
		}

		sr.Merge(nsr)

		opLimiter.Release(op.ID())
	}

	// As the components state task may take a long time to complete,
	// it is necessary to get the resource again to avoid overriding the status.
	res, err = in.modelClient.Resources().Query().
		Select(
			resource.FieldID,
			resource.FieldStatus).
		Where(resource.ID(res.ID)).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("error getting resource %s: %w", res.ID, err)
	}

	// State resource.
	// NB(alex): If the resource is transitioning, it should be skipped,
	// as the resource transitioning status may be changed by other tasks.
	if status.ResourceStatusUnDeployed.IsTrue(res) ||
		status.ResourceStatusDeployed.IsUnknown(res) ||
		status.ResourceStatusDeleted.Exist(res) ||
		status.ResourceStatusStopped.Exist(res) {
		// Skip if the resource is undeployed, on deploying, deleting or stopping.
		return berr
	}

	switch {
	case sr.Error:
		status.ResourceStatusReady.False(res, "")
	case sr.Transitioning:
		status.ResourceStatusReady.Unknown(res, "")
	default:
		status.ResourceStatusReady.True(res, "")
	}

	res.Status.SetSummary(status.WalkResource(&res.Status))

	err = in.modelClient.Resources().UpdateOne(res).
		SetStatus(res.Status).
		Exec(ctx)
	if err != nil {
		berr = multierr.Append(berr, err)
	}

	return berr
}
