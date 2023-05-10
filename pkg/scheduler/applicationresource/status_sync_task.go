package applicationresource

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/applicationresources"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operatorunknown"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/pkg/topic/datamessage"
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

	// NB(thxCode): in case of aggregate the status for application instance,
	// we group 10 application instances into one task group,
	// treat each application instance as a task unit,
	// and then process resources stating in task unit.
	var cnt, err = in.modelClient.ApplicationInstances().Query().
		Count(ctx)
	if err != nil {
		return fmt.Errorf("cannot count application instances: %w", err)
	}
	if cnt == 0 {
		return nil
	}
	// index none custom connectors for reusing.
	cs, err := listCandidateConnectors(ctx, in.modelClient)
	if err != nil {
		return fmt.Errorf("cannot list connectors: %w", err)
	}
	if len(cs) == 0 {
		return nil
	}
	var ops = make(map[types.ID]operator.Operator, len(cs))
	for i := range cs {
		var op, err = platform.GetOperator(ctx, operator.CreateOptions{
			Connector: *cs[i],
		})
		if err != nil {
			// warn out without breaking the whole syncing.
			in.logger.Warnf("cannot get operator of connector %q: %v", cs[i].ID, err)
			continue
		}
		if err = op.IsConnected(ctx); err != nil {
			// warn out without breaking the whole syncing.
			in.logger.Warnf("unreachable connector %q", cs[i].ID)
			// NB(thxCode): replace disconnected connector with unknown connector.
			op = operatorunknown.Operator{}
		}
		ops[cs[i].ID] = op
	}
	// execute tasks.
	const bks = 10
	var bkc = cnt/bks + 1
	if bkc == 1 {
		var st = in.buildStateTasks(ctx, 0, bks, ops)
		return st()
	}
	var wg = gopool.Group()
	for bk := 0; bk < bkc; bk++ {
		var st = in.buildStateTasks(ctx, bk*bks, bks, ops)
		wg.Go(st)
	}
	return wg.Wait()
}

func (in *StatusSyncTask) buildStateTasks(ctx context.Context, offset, limit int, ops map[types.ID]operator.Operator) func() error {
	return func() error {
		var is, err = in.modelClient.ApplicationInstances().Query().
			Order(model.Desc(applicationinstance.FieldCreateTime)).
			Offset(offset).
			Limit(limit).
			Unique(false).
			Select(
				applicationinstance.FieldID,
				applicationinstance.FieldApplicationID,
				applicationinstance.FieldStatus).
			All(ctx)
		if err != nil {
			return fmt.Errorf("error listing application instances in pagination %d,%d: %w",
				offset, limit, err)
		}
		var wg = gopool.Group()
		for i := range is {
			var at = in.buildStateTask(ctx, is[i], ops)
			wg.Go(at)
		}
		return wg.Wait()
	}
}

func (in *StatusSyncTask) buildStateTask(ctx context.Context, i *model.ApplicationInstance, ops map[types.ID]operator.Operator) func() error {
	return func() (berr error) {
		var rs, err = i.QueryResources().
			Order(model.Desc(applicationresource.FieldCreateTime)).
			Unique(false).
			Select(
				applicationresource.FieldID,
				applicationresource.FieldStatus,
				applicationresource.FieldInstanceID,
				applicationresource.FieldConnectorID,
				applicationresource.FieldType,
				applicationresource.FieldName,
				applicationresource.FieldDeployerType).
			All(ctx)
		if err != nil {
			return fmt.Errorf("error listing application resources: %w", err)
		}

		var ids = make(map[types.ID][]*model.ApplicationResource)
		for y := range rs {
			// group resources by connector.
			ids[rs[y].ConnectorID] = append(ids[rs[y].ConnectorID],
				rs[y])
		}

		var sr applicationresources.StateResult
		for cid, crs := range ids {
			var op, exist = ops[cid]
			if !exist {
				// ignore if not found operator.
				continue
			}
			nsr, err := applicationresources.State(ctx, op, in.modelClient, crs)
			if multierr.AppendInto(&berr, err) {
				// mark error as transitioning,
				// which doesn't flip the status.
				nsr.Transitioning = true
			}
			sr.Merge(nsr)
		}

		// state application instance.
		if status.ApplicationInstanceStatusDeleted.Exist(i) {
			// skip if the instance is on deleting.
			return
		}
		switch {
		case sr.Error:
			status.ApplicationInstanceStatusReady.False(i, "")
		case sr.Transitioning:
			status.ApplicationInstanceStatusReady.Unknown(i, "")
		default:
			status.ApplicationInstanceStatusReady.True(i, "")
		}
		update, err := dao.ApplicationInstanceStatusUpdate(in.modelClient, i)
		if err != nil {
			berr = multierr.Append(berr, err)
			return
		}

		err = update.Exec(ctx)
		if err != nil {
			berr = multierr.Append(berr, err)
		}

		// FIXME(thxCode): there are several places call the following publishing,
		//   we need an approach to unify them.
		err = datamessage.Publish(ctx, string(datamessage.Application), model.OpUpdate, []types.ID{i.ApplicationID})
		if err != nil {
			berr = multierr.Append(berr, err)
		}

		return
	}
}
