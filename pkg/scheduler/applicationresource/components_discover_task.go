package applicationresource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/operatorunknown"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

type ComponentsDiscoverTask struct {
	mu sync.Mutex

	modelClient model.ClientSet
	logger      log.Logger
}

func NewComponentsDiscoverTask(mc model.ClientSet) (*ComponentsDiscoverTask, error) {
	in := &ComponentsDiscoverTask{}
	in.modelClient = mc
	in.logger = log.WithName("task").WithName(in.Name())
	return in, nil
}

func (in *ComponentsDiscoverTask) Name() string {
	return "resource-components-discover"
}

func (in *ComponentsDiscoverTask) Process(ctx context.Context, args ...interface{}) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}
	var startTs = time.Now()
	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	// NB(thxCode): connectors are usually less in number,
	// in case of reuse the connection built from a connector,
	// we can treat each connector as a task group,
	// group 100 resources of each connector into one task unit,
	// and then process sub resources syncing in task unit.
	var cs, err = in.modelClient.Connectors().Query().
		Select(
			connector.FieldID,
			connector.FieldName,
			connector.FieldType,
			connector.FieldCategory,
			connector.FieldConfigVersion,
			connector.FieldConfigData).
		Where(connector.CategoryNEQ(types.ConnectorCategoryCustom)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("cannot list all connectors: %w", err)
	}
	if len(cs) == 0 {
		return nil
	}
	var wg = gopool.Group()
	for i := range cs {
		var st = in.buildSyncTasks(ctx, cs[i])
		wg.Go(st)
	}
	return wg.Wait()
}

func (in *ComponentsDiscoverTask) buildSyncTasks(ctx context.Context, c *model.Connector) func() error {
	return func() error {
		var op, err = platform.GetOperator(ctx, operator.CreateOptions{
			Connector: *c,
		})
		if err != nil {
			return err
		}
		if err = op.IsConnected(ctx); err != nil {
			// warn out without breaking the whole syncing.
			in.logger.Warnf("unreachable connector %q", c.ID)
			// NB(thxCode): replace disconnected connector with unknown connector.
			op = operatorunknown.Operator{}
		}

		cnt, err := c.QueryResources().
			Where(applicationresource.ModeNEQ(types.ApplicationResourceModeDiscovered)).
			Count(ctx)
		if err != nil {
			return fmt.Errorf("cannot count not discovered resources of connector %q: %w", c.ID, err)
		}
		if cnt == 0 {
			return nil
		}
		const bks = 100
		var bkc = cnt / bks
		if bkc == 0 {
			var at = in.buildSyncTask(ctx, op, 0, bks)
			return at()
		}
		var wg = gopool.Group()
		for bk := 0; bk < bkc; bk++ {
			var at = in.buildSyncTask(ctx, op, bk, bks)
			wg.Go(at)
		}
		return wg.Wait()
	}
}

func (in *ComponentsDiscoverTask) buildSyncTask(ctx context.Context, op operator.Operator, offset, limit int) func() error {
	return func() (berr error) {
		var entities, err = in.modelClient.ApplicationResources().Query().
			Order(model.Desc(applicationresource.FieldCreateTime)).
			Offset(offset).
			Limit(limit).
			Unique(false).
			Select(
				applicationresource.FieldID,
				applicationresource.FieldType,
				applicationresource.FieldModule,
				applicationresource.FieldInstanceID,
				applicationresource.FieldConnectorID,
				applicationresource.FieldName,
				applicationresource.FieldDeployerType).
			All(ctx)
		if err != nil {
			return err
		}

		for i := range entities {
			// get observed components from remote.
			var observedComps, err = op.GetComponents(ctx, entities[i])
			if multierr.AppendInto(&berr, err) {
				continue
			}
			if observedComps == nil {
				continue
			}

			// get record components from local.
			recordComps, err := in.modelClient.ApplicationResources().Query().
				Where(applicationresource.CompositionID(entities[i].ID)).
				All(ctx)
			if err != nil {
				return
			}

			// calculate creating list and deleting list.
			var observedCompsIndex = make(map[string]*model.ApplicationResource, len(observedComps))
			for j := range observedComps {
				var c = observedComps[j]
				observedCompsIndex[strs.Join("/", c.Type, c.Name)] = c
			}
			var deleteCompIDs = make([]oid.ID, 0, len(recordComps))
			for _, c := range recordComps {
				var k = strs.Join("/", c.Type, c.Name)
				if observedCompsIndex[k] != nil {
					delete(observedCompsIndex, k)
					continue
				}
				deleteCompIDs = append(deleteCompIDs, c.ID)
			}
			var createComps = make([]*model.ApplicationResource, 0, len(observedCompsIndex))
			for k := range observedCompsIndex {
				createComps = append(createComps, observedCompsIndex[k])
			}

			// create new components.
			if len(createComps) != 0 {
				creates, err := dao.ApplicationResourceCreates(in.modelClient, createComps...)
				if !multierr.AppendInto(&berr, err) {
					_, err = in.modelClient.ApplicationResources().CreateBulk(creates...).
						Save(ctx)
					if err != nil {
						berr = multierr.Append(berr, err)
					}
				}
			}

			// delete stale components.
			_, err = in.modelClient.ApplicationResources().Delete().
				Where(applicationresource.IDIn(deleteCompIDs...)).
				Exec(ctx)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				berr = multierr.Append(berr, err)
			}
		}

		return
	}
}
