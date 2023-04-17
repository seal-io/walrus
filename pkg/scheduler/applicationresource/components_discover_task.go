package applicationresource

import (
	"context"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
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

	var cnt, err = in.modelClient.ApplicationResources().Query().
		Where(applicationresource.ModeNEQ(types.ApplicationResourceModeDiscovered)).
		Count(ctx)
	if err != nil {
		return err
	}

	// divide processing buckets with count.
	const bks = 100
	var bkc = cnt / bks
	if bkc == 0 {
		var st = in.buildSyncTask(ctx, 0, bks)
		return st()
	}
	var wg = gopool.Group()
	for bk := 0; bk < bkc; bk++ {
		var st = in.buildSyncTask(ctx, bk, bks)
		wg.Go(st)
	}
	return wg.Wait()
}

func (in *ComponentsDiscoverTask) buildSyncTask(ctx context.Context, offset, limit int) func() error {
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
			WithConnector(func(cq *model.ConnectorQuery) {
				cq.Select(
					connector.FieldName,
					connector.FieldType,
					connector.FieldCategory,
					connector.FieldConfigVersion,
					connector.FieldConfigData)
			}).
			All(ctx)
		if err != nil {
			return err
		}

		for i := 0; i < len(entities); i++ {
			var op, err = platform.GetOperator(ctx, operator.CreateOptions{
				Connector: *entities[i].Edges.Connector,
			})
			if multierr.AppendInto(&berr, err) {
				continue
			}

			// get observed components from remote.
			observedComps, err := op.GetComponents(ctx, entities[i])
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
			for j := range deleteCompIDs {
				err = in.modelClient.ApplicationResources().DeleteOneID(deleteCompIDs[j]).
					Exec(ctx)
				if err != nil {
					berr = multierr.Append(berr, err)
				}
			}
		}

		return
	}
}
