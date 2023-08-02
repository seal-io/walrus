package serviceresource

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/operator"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/pkg/serviceresources"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type LabelApplyTask struct {
	mu sync.Mutex

	modelClient   model.ClientSet
	minBucketSize int
	logger        log.Logger
}

func NewLabelApplyTask(mc model.ClientSet, minBucketSizes map[string]int) (*LabelApplyTask, error) {
	in := &LabelApplyTask{
		modelClient:   mc,
		minBucketSize: 100,
	}

	if v, exist := minBucketSizes[in.Name()]; exist {
		in.minBucketSize = v
	}

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
	// we treat one connector as a process group,
	// and label several service resources in one process.

	cs, err := listCandidateConnectors(ctx, in.modelClient)
	if err != nil {
		return fmt.Errorf("cannot list all connectors: %w", err)
	}

	if len(cs) == 0 {
		return nil
	}

	wg := gopool.Group()

	for i := range cs {
		// Don't return directly when error occurs,
		// but records it and continue to handle the next connector,
		// the final error collect all errors,
		// and reports this time task running as failure at observing.
		pg := in.buildProcessGroup(ctx, cs[i], wg)
		err = multierr.Append(err, pg())

		if multierr.AppendInto(&err, ctx.Err()) {
			// Give up the loop if the context is canceled.
			break
		}
	}

	return multierr.Append(err, wg.Wait())
}

func (in *LabelApplyTask) buildProcessGroup(
	ctx context.Context,
	c *model.Connector,
	wg gopool.IWaitGroup,
) func() error {
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
			// Replace disconnected connector with unknown connector.
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

		bkc, bks := getBucket(cnt, in.minBucketSize)
		in.logger.Debugf("processing group %q within %d buckets, maximum %d items per bucket",
			c.ID, bkc, bks)

		for bk := 0; bk < bkc; bk++ {
			p := in.buildProcess(ctx, op, c.ID, bk*bks, bks)
			// NB(thxCode): we generally assume that the target sources of the connectors are all inconsistent,
			// if the target sources of multiple connectors point to the same address,
			// this may reach the bursting limit of the operator's client,
			// and harm the target source connected to the connector,
			// finally, result in a higher latency of the label applying.
			wg.Go(p)
		}

		return nil
	}
}

func (in *LabelApplyTask) buildProcess(
	ctx context.Context,
	op optypes.Operator,
	connectorID oid.ID,
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
