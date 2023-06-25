package server

import (
	"context"

	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/cron"
	"github.com/seal-io/seal/pkg/dao/model"
	connskd "github.com/seal-io/seal/pkg/scheduler/connector"
	serviceskd "github.com/seal-io/seal/pkg/scheduler/service"
	appresskd "github.com/seal-io/seal/pkg/scheduler/serviceresource"
	tokenskd "github.com/seal-io/seal/pkg/scheduler/token"
	"github.com/seal-io/seal/pkg/settings"
)

func (r *Server) initBackgroundJobs(ctx context.Context, opts initOptions) error {
	cs := cron.JobCreators{
		settings.ConnectorCostCollectCronExpr.Name(): buildConnectorCostCollectJobCreator(opts.ModelClient),
		settings.ConnectorStatusSyncCronExpr.Name():  buildConnectorStatusSyncJobCreator(opts.ModelClient),
		settings.ResourceStatusSyncCronExpr.Name():   buildResourceStatusSyncJobCreator(opts.ModelClient),
		settings.ResourceLabelApplyCronExpr.Name():   buildResourceLabelApplyJobCreator(opts.ModelClient),
		settings.ResourceComponentsDiscoverCronExpr.Name(): buildResourceComponentsDiscoverJobCreator(
			opts.ModelClient,
		),
		settings.TokenDeploymentExpiredCleanCronExpr.Name(): buildTokenDeploymentExpireCleanJobCreator(
			opts.ModelClient,
		),
		settings.ServiceDependencyCheckCronExpr.Name(): buildServiceDependencyCheckJobCreator(
			opts.ModelClient,
			opts.K8sConfig,
			opts.SkipTLSVerify,
		),
	}

	return cron.Register(ctx, opts.ModelClient, cs)
}

func buildConnectorCostCollectJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := connskd.NewCollectTask(mc)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildConnectorStatusSyncJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := connskd.NewStatusSyncTask(mc)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildResourceStatusSyncJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := appresskd.NewStatusSyncTask(mc)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildResourceLabelApplyJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := appresskd.NewLabelApplyTask(mc)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildResourceComponentsDiscoverJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := appresskd.NewComponentsDiscoverTask(mc)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildTokenDeploymentExpireCleanJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := tokenskd.NewDeploymentExpiredCleanTask(mc)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildServiceDependencyCheckJobCreator(mc model.ClientSet, kc *rest.Config, skipTLSVerify bool) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := serviceskd.NewServiceDependencyCheckTask(mc, kc, skipTLSVerify)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}
}
