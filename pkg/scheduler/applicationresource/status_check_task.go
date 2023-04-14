package applicationresource

import (
	"context"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type StatusCheckTask struct {
	mu sync.Mutex

	modelClient model.ClientSet
	logger      log.Logger
}

func NewStatusCheckTask(modelClient model.ClientSet) (*StatusCheckTask, error) {
	return &StatusCheckTask{
		modelClient: modelClient,
		logger:      log.WithName("task").WithName("application-resource").WithName("status-check"),
	}, nil
}

func (in *StatusCheckTask) Process(ctx context.Context, args ...interface{}) error {
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
		Count(ctx)
	if err != nil {
		return err
	}

	// divide processing buckets with count.
	const bks = 100
	var bkc = cnt / bks
	if bkc == 0 {
		var st = in.buildStateTask(ctx, 0, bks)
		return st()
	}
	var wg = gopool.Group()
	for bk := 0; bk < bkc; bk++ {
		var st = in.buildStateTask(ctx, bk, bks)
		wg.Go(st)
	}
	return wg.Wait()
}

func (in *StatusCheckTask) buildStateTask(ctx context.Context, offset, limit int) func() error {
	return func() (berr error) {
		var entities, err = in.modelClient.ApplicationResources().Query().
			Order(model.Desc(applicationresource.FieldCreateTime)).
			Offset(offset).
			Limit(limit).
			Unique(false).
			Select(
				applicationresource.FieldID,
				applicationresource.FieldStatus,
				applicationresource.FieldInstanceID,
				applicationresource.FieldConnectorID,
				applicationresource.FieldType,
				applicationresource.FieldName,
				applicationresource.FieldDeployerType).
			WithConnector(func(cq *model.ConnectorQuery) {
				cq.Select(
					connector.FieldName,
					connector.FieldType,
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
			// get status of the application resource.
			st, err := op.GetStatus(ctx, entities[i])
			if err != nil {
				berr = multierr.Append(berr, err)
			}
			// get endpoints of the application resource.
			eps, err := op.GetEndpoints(ctx, entities[i])
			if err != nil {
				berr = multierr.Append(berr, err)
			}
			// new application resource status.
			newStatus := types.ApplicationResourceStatus{
				Status:            *st,
				ResourceEndpoints: eps,
			}
			if entities[i].Status.Equal(newStatus) {
				// do not update if the status is same as previous.
				continue
			}

			err = in.modelClient.ApplicationResources().UpdateOne(entities[i]).
				SetStatus(newStatus).
				Exec(ctx)
			if err != nil {
				if model.IsNotFound(err) {
					// application resource has been deleted by other thread processing.
					continue
				}
				berr = multierr.Append(berr, err)
			}
		}
		return
	}
}
