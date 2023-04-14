package server

import (
	"context"

	"github.com/seal-io/seal/pkg/cron"
	"github.com/seal-io/seal/pkg/dao/model"
	appresskd "github.com/seal-io/seal/pkg/scheduler/applicationresource"
	connskd "github.com/seal-io/seal/pkg/scheduler/connector"
	costskd "github.com/seal-io/seal/pkg/scheduler/cost"

	"github.com/seal-io/seal/pkg/settings"
)

func (r *Server) initBackgroundTasks(ctx context.Context, opts initOptions) error {
	var cs = cron.JobCreators{
		settings.CostCollectCronExpr.Name():         buildCostCollectJobCreator(opts.ModelClient),
		settings.ConnectorCheckCronExpr.Name():      buildConnectorCheckJobCreator(opts.ModelClient),
		settings.ResourceStatusCheckCronExpr.Name(): buildResourceStatusCheckJobCreator(opts.ModelClient),
		settings.ResourceLabelApplyCronExpr.Name():  buildResourceLabelApplyJobCreator(opts.ModelClient),
	}
	return cron.Register(ctx, opts.ModelClient, cs)
}

func buildCostCollectJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		var task, err = costskd.NewCollectTask(mc)
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

func buildResourceStatusCheckJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		var task, err = appresskd.NewStatusCheckTask(mc)
		if err != nil {
			return nil, nil, err
		}
		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildResourceLabelApplyJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		var task, err = appresskd.NewLabelApplyTask(mc)
		if err != nil {
			return nil, nil, err
		}
		return cron.ImmediateExpr(expr), task, nil
	}
}
