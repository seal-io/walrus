package serviceresource

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
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/operator"
	optypes "github.com/seal-io/seal/pkg/operator/types"
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
	startTs := time.Now()

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	// NB(thxCode): connectors are usually less in number,
	// in case of reuse the connection built from a connector,
	// we can treat each connector as a task group,
	// group 100 resources of each connector into one task unit,
	// and then process sub resources syncing in task unit.
	cs, err := listCandidateConnectors(ctx, in.modelClient)
	if err != nil {
		return fmt.Errorf("cannot list all connectors: %w", err)
	}

	if len(cs) == 0 {
		return nil
	}
	wg := gopool.Group()

	for i := range cs {
		st := in.buildSyncTasks(ctx, cs[i])
		wg.Go(st)
	}

	return wg.Wait()
}

func (in *ComponentsDiscoverTask) buildSyncTasks(ctx context.Context, c *model.Connector) func() error {
	return func() error {
		op, err := operator.Get(ctx, optypes.CreateOptions{
			Connector: *c,
		})
		if err != nil {
			return err
		}

		if err = op.IsConnected(ctx); err != nil {
			// Warn out without breaking the whole syncing.
			in.logger.Warnf("unreachable connector %q", c.ID)
			// NB(thxCode): replace disconnected connector with unknown connector.
			op = operator.UnReachable()
		}

		cnt, err := c.QueryResources().
			Where(serviceresource.ModeNEQ(types.ServiceResourceModeDiscovered)).
			Count(ctx)
		if err != nil {
			return fmt.Errorf("cannot count not discovered resources of connector %q: %w", c.ID, err)
		}

		if cnt == 0 {
			return nil
		}

		const bks = 100

		bkc := cnt/bks + 1
		if bkc == 1 {
			at := in.buildSyncTask(ctx, op, c.ID, 0, bks)
			return at()
		}
		wg := gopool.Group()

		for bk := 0; bk < bkc; bk++ {
			at := in.buildSyncTask(ctx, op, c.ID, bk*bks, bks)
			wg.Go(at)
		}

		return wg.Wait()
	}
}

func (in *ComponentsDiscoverTask) buildSyncTask(
	ctx context.Context,
	op optypes.Operator,
	connectorID oid.ID,
	offset,
	limit int,
) func() error {
	return func() (berr error) {
		entities, err := in.modelClient.ServiceResources().Query().
			Where(
				serviceresource.ModeNEQ(types.ServiceResourceModeDiscovered),
				serviceresource.ConnectorID(connectorID)).
			Order(model.Desc(serviceresource.FieldCreateTime)).
			Offset(offset).
			Limit(limit).
			Unique(false).
			Select(
				serviceresource.FieldID,
				serviceresource.FieldProjectID,
				serviceresource.FieldType,
				serviceresource.FieldServiceID,
				serviceresource.FieldConnectorID,
				serviceresource.FieldName,
				serviceresource.FieldDeployerType).
			All(ctx)
		if err != nil {
			return err
		}

		for i := range entities {
			// Get observed components from remote.
			observedComps, err := op.GetComponents(ctx, entities[i])
			if multierr.AppendInto(&berr, err) {
				continue
			}

			if observedComps == nil {
				continue
			}

			// Get record components from local.
			recordComps, err := in.modelClient.ServiceResources().Query().
				Where(serviceresource.CompositionID(entities[i].ID)).
				All(ctx)
			if err != nil {
				return
			}

			// Calculate creating list and deleting list.
			observedCompsIndex := make(map[string]*model.ServiceResource, len(observedComps))

			for j := range observedComps {
				c := observedComps[j]
				observedCompsIndex[strs.Join("/", c.Type, c.Name)] = c
			}
			deleteCompIDs := make([]oid.ID, 0, len(recordComps))

			for _, c := range recordComps {
				k := strs.Join("/", c.Type, c.Name)
				if observedCompsIndex[k] != nil {
					delete(observedCompsIndex, k)
					continue
				}

				deleteCompIDs = append(deleteCompIDs, c.ID)
			}

			createComps := make([]*model.ServiceResource, 0, len(observedCompsIndex))
			for k := range observedCompsIndex {
				createComps = append(createComps, observedCompsIndex[k])
			}

			// Create new components.
			if len(createComps) != 0 {
				creates, err := dao.ServiceResourceCreates(in.modelClient, createComps...)
				if !multierr.AppendInto(&berr, err) {
					_, err = in.modelClient.ServiceResources().CreateBulk(creates...).
						Save(ctx)
					if err != nil {
						berr = multierr.Append(berr, err)
					}
				}
			}

			// Delete stale components.
			if len(deleteCompIDs) != 0 {
				_, err = in.modelClient.ServiceResources().Delete().
					Where(serviceresource.IDIn(deleteCompIDs...)).
					Exec(ctx)
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					berr = multierr.Append(berr, err)
				}
			}
		}

		return
	}
}