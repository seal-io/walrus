package scheduler

import (
	"context"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type ResourceStatusCheckTask struct {
	mu sync.Mutex

	modelClient model.ClientSet
	logger      log.Logger
}

func NewResourceStatusCheckTask(modelClient model.ClientSet) (*ResourceStatusCheckTask, error) {
	return &ResourceStatusCheckTask{
		modelClient: modelClient,
		logger:      log.WithName("resource").WithName("state"),
	}, nil
}

func (in *ResourceStatusCheckTask) Process(ctx context.Context, args ...interface{}) error {
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

func (in *ResourceStatusCheckTask) buildStateTask(ctx context.Context, offset, limit int) func() error {
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
			if st == entities[i].Status {
				// do not update if the status is same as previous.
				continue
			}
			// update status of the application resource.
			// TODO(thxCode): dig out the detail of status,
			//   update to the status message.
			err = in.modelClient.ApplicationResources().UpdateOne(entities[i]).
				SetStatus(st).
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
