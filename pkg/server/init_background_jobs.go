package server

import (
	"context"

	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/cron"
	"github.com/seal-io/seal/pkg/dao/model"
	connskd "github.com/seal-io/seal/pkg/scheduler/connector"
	svcskd "github.com/seal-io/seal/pkg/scheduler/service"
	svcresskd "github.com/seal-io/seal/pkg/scheduler/serviceresource"
	telemetryskd "github.com/seal-io/seal/pkg/scheduler/telemetry"
	tokenskd "github.com/seal-io/seal/pkg/scheduler/token"
	"github.com/seal-io/seal/pkg/settings"
)

func (r *Server) initBackgroundJobs(ctx context.Context, opts initOptions) error {
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
			opts.TlsCertified,
		),
		settings.ServiceDriftDetectCronExpr.Name(): buildServiceDriftDetectJobCreator(
			opts.ModelClient,
			opts.K8sConfig,
			opts.TlsCertified,
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
		task, err := svcresskd.NewStatusSyncTask(mc)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildResourceLabelApplyJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := svcresskd.NewLabelApplyTask(mc)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildResourceComponentsDiscoverJobCreator(mc model.ClientSet) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := svcresskd.NewComponentsDiscoverTask(mc)
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

func buildServiceRelationshipCheckJobCreator(mc model.ClientSet, kc *rest.Config, tlsCertified bool) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := svcskd.NewServiceRelationshipCheckTask(mc, kc, tlsCertified)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}
}

func buildServiceDriftDetectJobCreator(mc model.ClientSet, kc *rest.Config, tlsCertified bool) cron.JobCreator {
	return func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error) {
		task, err := svcskd.NewServiceDriftDetectTask(mc, kc, tlsCertified)
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
