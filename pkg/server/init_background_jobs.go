package server

import (
	"context"

	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/cron"
	"github.com/seal-io/walrus/pkg/dao/model"
	catalogskd "github.com/seal-io/walrus/pkg/scheduler/catalog"
	connskd "github.com/seal-io/walrus/pkg/scheduler/connector"
	serviceskd "github.com/seal-io/walrus/pkg/scheduler/service"
	appresskd "github.com/seal-io/walrus/pkg/scheduler/serviceresource"
	telemetryskd "github.com/seal-io/walrus/pkg/scheduler/telemetry"
	tokenskd "github.com/seal-io/walrus/pkg/scheduler/token"
	"github.com/seal-io/walrus/pkg/settings"
)

// startBackgroundJobs starts the background jobs by Cron Expression to do something periodically.
func (r *Server) startBackgroundJobs(ctx context.Context, opts initOptions) error {
	cs := cron.JobCreators{
		settings.ConnectorCostCollectCronExpr.Name():    buildConnectorCostCollectJobCreator(opts.ModelClient),
		settings.ConnectorStatusSyncCronExpr.Name():     buildConnectorStatusSyncJobCreator(opts.ModelClient),
		settings.ResourceStatusSyncCronExpr.Name():      buildResourceStatusSyncJobCreator(opts.ModelClient),
		settings.ResourceLabelApplyCronExpr.Name():      buildResourceLabelApplyJobCreator(opts.ModelClient),
		settings.TelemetryPeriodicReportCronExpr.Name(): buildTelemetryPeriodicReportJobCreator(opts.ModelClient),
		settings.ResourceComponentsDiscoverCronExpr.Name(): buildResourceComponentsDiscoverJobCreator(
			opts.ModelClient,
		),
		settings.TokenDeploymentExpiredCleanCronExpr.Name(): buildTokenDeploymentExpireCleanJobCreator(
			opts.ModelClient,
		),
		settings.ServiceRelationshipCheckCronExpr.Name(): buildServiceRelationshipCheckJobCreator(
			opts.ModelClient,
			opts.K8sConfig,
			opts.SkipTLSVerify,
		),
		settings.CatalogTemplateSyncCronExpr.Name(): buildCatalogTemplateSyncJobCreator(opts.ModelClient),
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

func buildServiceRelationshipCheckJobCreator(mc model.ClientSet, kc *rest.Config, skipTLSVerify bool) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := serviceskd.NewServiceRelationshipCheckTask(mc, kc, skipTLSVerify)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildTelemetryPeriodicReportJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := telemetryskd.NewPeriodicReportTask(mc)
		if err != nil {
			return nil, nil, err
		}

		return cron.AwaitedExpr(expr), task, nil
	}
}

func buildCatalogTemplateSyncJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := catalogskd.NewCatalogTemplateSyncTask(mc)
		if err != nil {
			return nil, nil, err
		}

		return cron.AwaitedExpr(expr), task, nil
	}
}
