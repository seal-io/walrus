package server

import (
	"context"

	costskd "github.com/seal-io/seal/pkg/costs/scheduler"
	"github.com/seal-io/seal/pkg/cron"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/settings"
)

func (r *Server) initBackgroundTasks(ctx context.Context, opts initOptions) error {
	var cs = cron.JobCreators{
		settings.CostCollectCronExpr.Name():    buildCostCollectJobCreator(opts.ModelClient),
		settings.CostToolsCheckCronExpr.Name(): buildCostToolsCheckJobCreator(opts.ModelClient),
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

func buildCostToolsCheckJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		var task, err = costskd.NewToolsCheckTask(mc)
		if err != nil {
			return nil, nil, err
		}
		return cron.ImmediateExpr(expr), task, nil
	}
}
