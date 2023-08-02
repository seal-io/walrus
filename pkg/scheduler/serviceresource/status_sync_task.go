package serviceresource

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operator"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/pkg/serviceresources"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type StatusSyncTask struct {
	mu sync.Mutex

	modelClient   model.ClientSet
	minBucketSize int
	logger        log.Logger
}

func NewStatusSyncTask(mc model.ClientSet, minBucketSizes map[string]int) (*StatusSyncTask, error) {
	in := &StatusSyncTask{
		modelClient:   mc,
		minBucketSize: 100,
	}

	if v, exist := minBucketSizes[in.Name()]; exist {
		in.minBucketSize = v
	}

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
	startTs := time.Now()

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	// NB(thxCode): in case of aggregate the status for service,
	// we group several services as a process,
	// and synchronize their statuses in one process.

	cnt, err := in.modelClient.Services().Query().
		Count(ctx)
	if err != nil {
		return fmt.Errorf("cannot count services: %w", err)
	}

	if cnt == 0 {
		return nil
	}

	// Index none custom connectors for reusing.
	cs, err := listCandidateConnectors(ctx, in.modelClient)
	if err != nil {
		return fmt.Errorf("cannot list connectors: %w", err)
	}

	if len(cs) == 0 {
		return nil
	}

	ops := make(map[oid.ID]optypes.Operator, len(cs))

	for i := range cs {
		op, err := operator.Get(ctx, optypes.CreateOptions{
			Connector: *cs[i],
		})
		if err != nil {
			// Warn out without breaking the whole syncing.
			in.logger.Warnf("cannot get operator of connector %q: %v", cs[i].ID, err)
			continue
		}

		if err = op.IsConnected(ctx); err != nil {
			// Warn out without breaking the whole syncing.
			in.logger.Warnf("unreachable connector %q", cs[i].ID)
			// NB(thxCode): replace disconnected connector with unknown connector.
			op = operator.UnReachable()
		}
		ops[cs[i].ID] = op
	}

	bkc, bks := getBucket(cnt, in.minBucketSize)
	in.logger.Debugf("processing within %d buckets, maximum %d items per bucket",
		bkc, bks)

	wg := gopool.Group()

	for bk := 0; bk < bkc; bk++ {
		p := in.buildProcess(ctx, bk*bks, bks, ops)
		// NB(thxCode): we generally assume that the target sources of the connectors are all inconsistent,
		// if the target sources of multiple connectors point to the same address,
		// this may reach the bursting limit of the operator's client,
		// and harm the target source connected to the connector,
		// finally, result in a higher latency of the status syncing.
		wg.Go(p)
	}

	return wg.Wait()
}

func (in *StatusSyncTask) buildProcess(
	ctx context.Context,
	offset,
	limit int,
	ops map[oid.ID]optypes.Operator,
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
			return fmt.Errorf("error listing services in pagination %d,%d: %w",
				offset, limit, err)
		}

		if len(svcs) == 0 {
			return nil
		}

		for i := range svcs {
			// Don't return directly when error occurs,
			// but records it and continue to handle the next connector,
			// the final error collect all errors,
			// and reports this time task running as failure at observing.
			err = multierr.Append(err, in.process(ctx, svcs[i], ops))

			if multierr.AppendInto(&err, ctx.Err()) {
				// Give up the loop if the context is canceled.
				break
			}
		}

		return err
	}
}

func (in *StatusSyncTask) process(
	ctx context.Context,
	svc *model.Service,
	ops map[oid.ID]optypes.Operator,
) (berr error) {
	rs, err := svc.QueryResources().
		Order(model.Desc(serviceresource.FieldCreateTime)).
		Unique(false).
		Select(
			serviceresource.FieldID,
			serviceresource.FieldStatus,
			serviceresource.FieldServiceID,
			serviceresource.FieldConnectorID,
			serviceresource.FieldType,
			serviceresource.FieldName,
			serviceresource.FieldDeployerType).
		Where(serviceresource.Shape(types.ServiceResourceShapeInstance)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("error listing service resources: %w", err)
	}

	ids := make(map[oid.ID][]*model.ServiceResource)
	for y := range rs {
		// Group resources by connector.
		ids[rs[y].ConnectorID] = append(ids[rs[y].ConnectorID],
			rs[y])
	}

	var sr serviceresources.StateResult

	for cid, crs := range ids {
		op, exist := ops[cid]
		if !exist {
			// Ignore if not found operator.
			continue
		}

		nsr, err := serviceresources.State(ctx, op, in.modelClient, crs)
		if multierr.AppendInto(&berr, err) {
			// Mark error as transitioning,
			// which doesn't flip the status.
			nsr.Transitioning = true
		}

		sr.Merge(nsr)
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

	update, err := dao.ServiceStatusUpdate(in.modelClient, svc)
	if err != nil {
		return multierr.Append(berr, err)
	}

	err = update.Exec(ctx)
	if err != nil {
		berr = multierr.Append(berr, err)
	}

	return berr
}
