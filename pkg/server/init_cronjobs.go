package server

import (
	"context"

	"github.com/seal-io/seal/pkg/cron"
	"github.com/seal-io/seal/pkg/dao/model"
	connskd "github.com/seal-io/seal/pkg/scheduler/connector"
	costskd "github.com/seal-io/seal/pkg/scheduler/cost"

	"github.com/seal-io/seal/pkg/settings"
)

func (r *Server) initBackgroundTasks(ctx context.Context, opts initOptions) error {
	var cs = cron.JobCreators{
		settings.CostCollectCronExpr.Name():    buildCostCollectJobCreator(opts.ModelClient),
		settings.CostToolsCheckCronExpr.Name(): buildConnectorCheckJobCreator(opts.ModelClient),
	}
	return cron.Register(ctx, opts.ModelClient, cs)
}

func buildCostCollectJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		var task, err = costskd.NewCostSyncTask(mc)
		if err != nil {
			return nil, nil, err
		}
		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildConnectorCheckJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		var task, err = connskd.NewStatusCheckTask(mc)
		if err != nil {
			return nil, nil, err
		}
		return cron.ImmediateExpr(expr), task, nil
	}
}
