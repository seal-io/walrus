package serviceresource

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/pkg/operator"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/pkg/serviceresources"
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
	startTs := time.Now()

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	// NB(thxCode): connectors are usually less in number,
	// in case of reuse the connection built from a connector,
	// we can treat each connector as a task group,
	// group 100 resources of each connector into one task unit,
	// and then process resources labeling in task unit.
	cs, err := listCandidateConnectors(ctx, in.modelClient)
	if err != nil {
		return fmt.Errorf("cannot list all connectors: %w", err)
	}

	if len(cs) == 0 {
		return nil
	}
	wg := gopool.Group()

	for i := range cs {
		at := in.buildApplyTasks(ctx, cs[i])
		wg.Go(at)
	}

	return wg.Wait()
}

func (in *LabelApplyTask) buildApplyTasks(ctx context.Context, c *model.Connector) func() error {
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
			Count(ctx)
		if err != nil {
			return fmt.Errorf("cannot count resources of connector %q: %w", c.ID, err)
		}

		if cnt == 0 {
			return nil
		}

		const bks = 100

		bkc := cnt/bks + 1
		if bkc == 1 {
			at := in.buildApplyTask(ctx, op, c.ID, 0, bks)
			return at()
		}
		wg := gopool.Group()

		for bk := 0; bk < bkc; bk++ {
			at := in.buildApplyTask(ctx, op, c.ID, bk*bks, bks)
			wg.Go(at)
		}

		return wg.Wait()
	}
}

func (in *LabelApplyTask) buildApplyTask(
	ctx context.Context,
	op optypes.Operator,
	connectorID object.ID,
	offset,
	limit int,
) func() error {
	return func() error {
		entities, err := serviceresources.ListCandidatesPageByConnector(
			ctx, in.modelClient, connectorID, offset, limit)
		if err != nil {
			return fmt.Errorf("error listing label candidates: %w", err)
		}

		return serviceresources.Label(ctx, op, entities)
	}
}
