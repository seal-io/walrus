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
	"github.com/seal-io/walrus/pkg/operator"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
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

func (in *StatusSyncTask) Name() string {
	return "resource-status-sync"
}

func (in *StatusSyncTask) Process(ctx context.Context, args ...any) error {
	// NB(thxCode): in case of aggregate the status for service,
	// we group 10 services into one task group,
	// treat each service as a task unit,
	// and then process resources stating in task unit.
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
	ops := make(map[object.ID]optypes.Operator, len(cs))

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
	// Execute tasks.
	const bks = 10

	bkc := cnt/bks + 1
	if bkc == 1 {
		st := in.buildStateTasks(ctx, 0, bks, ops)
		return st()
	}
	wg := gopool.Group()

	for bk := 0; bk < bkc; bk++ {
		st := in.buildStateTasks(ctx, bk*bks, bks, ops)
		wg.Go(st)
	}

	return wg.Wait()
}

func (in *StatusSyncTask) buildStateTasks(
	ctx context.Context,
	offset,
	limit int,
	ops map[object.ID]optypes.Operator,
) func() error {
	return func() error {
		is, err := in.modelClient.Services().Query().
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
		wg := gopool.Group()

		for i := range is {
			at := in.buildStateTask(ctx, is[i], ops)
			wg.Go(at)
		}

		return wg.Wait()
	}
}

func (in *StatusSyncTask) buildStateTask(
	ctx context.Context,
	i *model.Service,
	ops map[object.ID]optypes.Operator,
) func() error {
	return func() (berr error) {
		rs, err := i.QueryResources().
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

		ids := make(map[object.ID][]*model.ServiceResource)
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
		if status.ServiceStatusDeleted.Exist(i) {
			// Skip if the service is on deleting.
			return
		}

		switch {
		case sr.Error:
			status.ServiceStatusReady.False(i, "")
		case sr.Transitioning:
			status.ServiceStatusReady.Unknown(i, "")
		default:
			status.ServiceStatusReady.True(i, "")
		}

		i.Status.SetSummary(status.WalkService(&i.Status))

		err = in.modelClient.Services().UpdateOne(i).
			SetStatus(i.Status).
			Exec(ctx)
		if err != nil {
			berr = multierr.Append(berr, err)
		}

		return
	}
}
