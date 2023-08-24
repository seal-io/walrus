package serviceresource

import (
	"context"
	"fmt"
	"runtime"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/operator"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/utils/log"
)

type (
	operatorIndexer map[object.ID]optypes.Operator
	operatorLimiter map[string]chan struct{}
)

func (t operatorLimiter) Acquire(id string) {
	if _, ok := t[id]; !ok {
		return
	}

	t[id] <- struct{}{}
}

func (t operatorLimiter) Release(id string) {
	if _, ok := t[id]; !ok {
		return
	}

	<-t[id]
	runtime.Gosched()
}

// retrieveOperators fetches all operators and constructs two structures as results,
// one is an indexer that indexing the operator by the related connector ID,
// another one is a limiter that assisting the caller to control the operations' parallelism of the operator.
func retrieveOperators(
	ctx context.Context,
	modelClient model.ClientSet,
	logger log.Logger,
) (
	operatorIndexer,
	operatorLimiter,
	error,
) {
	cs, err := modelClient.Connectors().Query().
		Select(
			connector.FieldID,
			connector.FieldName,
			connector.FieldType,
			connector.FieldCategory,
			connector.FieldConfigVersion,
			connector.FieldConfigData).
		Where(
			connector.CategoryNEQ(types.ConnectorCategoryCustom),
			connector.CategoryNEQ(types.ConnectorCategoryVersionControl)).
		All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot list connectors: %w", err)
	}

	if len(cs) == 0 {
		return nil, nil, nil
	}

	// Index the operators by the related connector ID.
	indexer := make(map[object.ID]optypes.Operator, len(cs))

	// Record the operators' parallelism by the corresponding ID.
	limiter := make(operatorLimiter)

	for i := range cs {
		var op optypes.Operator

		op, err = operator.Get(ctx, optypes.CreateOptions{
			Connector: *cs[i],
		})
		if err != nil {
			// Warn out without breaking the whole syncing.
			logger.Warnf("cannot get operator of connector %q: %v", cs[i].ID, err)
			continue
		}

		if err = op.IsConnected(ctx); err != nil {
			// Warn out without breaking the whole syncing.
			logger.Warnf("unreachable connector %q", cs[i].ID)
			// Replace disconnected connector with unknown connector.
			op = operator.UnReachable()
		}

		// Index operators by connector ID.
		indexer[cs[i].ID] = op

		// Group connector ID list by operator ID.
		if id, burst := op.ID(), op.Burst(); id != "" && burst > 0 {
			if _, ok := limiter[id]; !ok {
				limiter[id] = make(chan struct{}, burst)
			}
		}
	}

	return indexer, limiter, nil
}

// getBatches returns the size of a batch and the count of batches.
//
// The total indicates the total number of units to be processed.
// The burst indicates the maximum number of parallel batch operations.
// The minSize can avoid excessively small pages that result in excessive database reads.
//
// The returning size indicates how many units should be processed in a batch serially,
// and the returning count indicates how many batches should be processed in parallel.
//
// The multiple of returning size and count must be greater than or equal to the given total,
// and the returning count must be less than or equal to the given burst.
func getBatches(total, burst, minSize int) (size, count int) {
	if total < 0 {
		total = 0
	}

	if burst <= 0 {
		burst = 1
	}

	if minSize <= 0 {
		minSize = burst
	}

	size = total / burst
	count = burst

	// Avoid excessively small pages that result in excessive database reads,
	// let the bucket size be the minimum value and recalculate the bucket count.
	if size < minSize {
		size = minSize
		count = total / size
	}

	// Avoid missing the last page.
	if count*size < total {
		// Turn up the count at first.
		count++

		// Turn down the count and turn up the size if the count is greater than the burst.
		if count > burst {
			count--
			size++
		}
	}

	return
}
