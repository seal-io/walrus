package applicationresource

import (
	"context"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type ResourceLabelApplyTask struct {
	mu sync.Mutex

	modelClient model.ClientSet
	logger      log.Logger
}

func NewResourceLabelApplyTask(modelClient model.ClientSet) (*ResourceLabelApplyTask, error) {
	return &ResourceLabelApplyTask{
		modelClient: modelClient,
		logger:      log.WithName("resource").WithName("label-apply"),
	}, nil
}

func (in *ResourceLabelApplyTask) Process(ctx context.Context, args ...interface{}) error {
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
		var st = in.applyLabels(ctx, 0, bks)
		return st()
	}
	var wg = gopool.Group()
	for bk := 0; bk < bkc; bk++ {
		var st = in.applyLabels(ctx, bk, bks)
		wg.Go(st)
	}
	return wg.Wait()
}

func (in *ResourceLabelApplyTask) applyLabels(ctx context.Context, offset, limit int) func() error {
	return func() (berr error) {
		var entities, err = in.modelClient.ApplicationResources().Query().
			Order(model.Desc(applicationresource.FieldCreateTime)).
			Offset(offset).
			Limit(limit).
			Unique(false).
			Where(
				applicationresource.TypeIn(intercept.TFLabeledTypes...),
			).
			Select(
				applicationresource.FieldID,
				applicationresource.FieldInstanceID,
				applicationresource.FieldConnectorID,
				applicationresource.FieldType,
				applicationresource.FieldName,
				applicationresource.FieldDeployerType).
			WithInstance(func(aiq *model.ApplicationInstanceQuery) {
				aiq.Select(
					applicationinstance.FieldApplicationID,
					applicationinstance.FieldEnvironmentID,
				).WithEnvironment(func(
					eq *model.EnvironmentQuery) {
					eq.Select(
						environment.FieldID,
						environment.FieldName,
					)
				}).WithApplication(func(aq *model.ApplicationQuery) {
					aq.Select(
						application.FieldID,
						application.FieldName,
						application.FieldProjectID,
					).WithProject(func(pq *model.ProjectQuery) {
						pq.Select(
							project.FieldID,
							project.FieldName,
						)
					})
				})
			}).
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

		var connAppRes = make(map[*model.Connector][]*model.ApplicationResource)
		for i := 0; i < len(entities); i++ {
			conn := entities[i].Edges.Connector
			connAppRes[conn] = append(connAppRes[conn], entities[i])
		}

		for conn, ress := range connAppRes {
			var op, err = platform.GetOperator(ctx, operator.CreateOptions{
				Connector: *conn,
			})
			if multierr.AppendInto(&berr, err) {
				continue
			}

			for _, res := range ress {
				var (
					appName     string
					projectName string
					envName     string
				)
				instance := res.Edges.Instance
				if instance == nil {
					continue
				}

				// application
				app := instance.Edges.Application
				if app != nil {
					appName = app.Name

					// project
					if app.Edges.Project != nil {
						projectName = app.Edges.Project.Name
					}
				}

				// environment
				if instance.Edges.Environment != nil {
					envName = instance.Edges.Environment.Name
				}

				labels := map[string]string{
					types.LabelSealEnvironment: envName,
					types.LabelSealProject:     projectName,
					types.LabelSealApplication: appName,
				}
				err = op.Label(ctx, res, labels)
				if multierr.AppendInto(&berr, err) {
					continue
				}
			}
		}
		return
	}
}
