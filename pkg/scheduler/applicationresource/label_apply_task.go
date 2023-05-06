package applicationresource

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/seal-io/seal/pkg/applicationresources"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/operatorunknown"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type LabelApplyTask struct {
	mu sync.Mutex

	modelClient model.ClientSet
	logger      log.Logger
}

func NewLabelApplyTask(mc model.ClientSet) (*LabelApplyTask, error) {
	in := &LabelApplyTask{}
	in.modelClient = mc
	in.logger = log.WithName("task").WithName(in.Name())
	return in, nil
}

func (in *LabelApplyTask) Name() string {
	return "resource-label-apply"
}

func (in *LabelApplyTask) Process(ctx context.Context, args ...interface{}) error {
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
	// and then process resources labeling in task unit.
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
		var at = in.buildApplyTasks(ctx, cs[i])
		wg.Go(at)
	}
	return wg.Wait()
}

func (in *LabelApplyTask) buildApplyTasks(ctx context.Context, c *model.Connector) func() error {
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
			Count(ctx)
		if err != nil {
			return fmt.Errorf("cannot count resources of connector %q: %w", c.ID, err)
		}
		if cnt == 0 {
			return nil
		}
		const bks = 100
		var bkc = cnt / bks
		if bkc == 0 {
			var at = in.buildApplyTask(ctx, op, 0, bks)
			return at()
		}
		var wg = gopool.Group()
		for bk := 0; bk < bkc; bk++ {
			var at = in.buildApplyTask(ctx, op, bk, bks)
			wg.Go(at)
		}
		return wg.Wait()
	}
}

func (in *LabelApplyTask) buildApplyTask(ctx context.Context, op operator.Operator, offset, limit int) func() error {
	return func() error {
		var entities, err = applicationresources.ListCandidatesByPage(
			ctx, in.modelClient, offset, limit)
		if err != nil {
			return fmt.Errorf("error listing label candidates: %w", err)
		}
		return applicationresources.Label(ctx, op, entities)
	}
}
