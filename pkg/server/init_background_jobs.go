package server

import (
	"context"

	"github.com/seal-io/walrus/pkg/cron"
	catalogskd "github.com/seal-io/walrus/pkg/scheduler/catalog"
	connskd "github.com/seal-io/walrus/pkg/scheduler/connector"
	svcskd "github.com/seal-io/walrus/pkg/scheduler/resource"
	svcresskd "github.com/seal-io/walrus/pkg/scheduler/resourcecomponent"
	telemetryskd "github.com/seal-io/walrus/pkg/scheduler/telemetry"
	tokenskd "github.com/seal-io/walrus/pkg/scheduler/token"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/log"
)

// startBackgroundJobs starts the background jobs by Cron Expression to do something periodically.
// StartBackgroundJobs requires the global settings to be initialized.
func (r *Server) startBackgroundJobs(ctx context.Context, opts initOptions) error {
	bjs := []jobCreatorBuilder{
		buildCatalogTemplateSyncJobCreator,
		buildConnectorCostCollectJobCreator,
		buildConnectorStatusSyncJobCreator,
		buildResourceComponentsDiscoverJobCreator,
		buildResourceLabelApplyJobCreator,
		buildResourceStatusSyncJobCreator,
		buildResourceRelationshipCheckJobCreator,
		buildTelemetryPeriodicReportJobCreator,
		buildTokenDeploymentExpireCleanJobCreator,
	}

	js := cron.JobCreators{}

	for i := range bjs {
		expr, j := bjs[i](opts)
		js[expr.Name()] = j
	}

	return cron.Register(ctx, opts.ModelClient, js)
}

// jobCreatorBuilder is the stereotype of a function that creates a cron.JobCreator,
// must return the expression setting and the job creator.
type jobCreatorBuilder func(initOptions) (exprSetting settings.Value, jobCreator cron.JobCreator)

func buildCatalogTemplateSyncJobCreator(opts initOptions) (es settings.Value, jc cron.JobCreator) {
	es = settings.CatalogTemplateSyncCronExpr
	jc = func(logger log.Logger, expr string) (cron.Expr, cron.Task, error) {
		task, err := catalogskd.NewCatalogTemplateSyncTask(logger, opts.ModelClient)
		if err != nil {
			return nil, nil, err
		}

		return cron.AwaitedExpr(expr), task, nil
	}

	return
}

func buildConnectorCostCollectJobCreator(opts initOptions) (es settings.Value, jc cron.JobCreator) {
	es = settings.ConnectorCostCollectCronExpr
	jc = func(logger log.Logger, expr string) (cron.Expr, cron.Task, error) {
		task, err := connskd.NewCollectTask(logger, opts.ModelClient)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}

	return
}

func buildConnectorStatusSyncJobCreator(opts initOptions) (es settings.Value, jc cron.JobCreator) {
	es = settings.ConnectorStatusSyncCronExpr
	jc = func(logger log.Logger, expr string) (cron.Expr, cron.Task, error) {
		task, err := connskd.NewStatusSyncTask(logger, opts.ModelClient)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}

	return
}

func buildResourceComponentsDiscoverJobCreator(opts initOptions) (es settings.Value, jc cron.JobCreator) {
	es = settings.ResourceComponentsDiscoverCronExpr
	jc = func(logger log.Logger, expr string) (cron.Expr, cron.Task, error) {
		task, err := svcresskd.NewComponentsDiscoverTask(logger, opts.ModelClient)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}

	return
}

func buildResourceLabelApplyJobCreator(opts initOptions) (es settings.Value, jc cron.JobCreator) {
	es = settings.ResourceComponentLabelApplyCronExpr
	jc = func(logger log.Logger, expr string) (cron.Expr, cron.Task, error) {
		task, err := svcresskd.NewLabelApplyTask(logger, opts.ModelClient)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}

	return
}

func buildResourceStatusSyncJobCreator(opts initOptions) (es settings.Value, jc cron.JobCreator) {
	es = settings.ResourceComponentStatusSyncCronExpr
	jc = func(logger log.Logger, expr string) (cron.Expr, cron.Task, error) {
		task, err := svcresskd.NewStatusSyncTask(logger, opts.ModelClient)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}

	return
}

func buildResourceRelationshipCheckJobCreator(opts initOptions) (es settings.Value, jc cron.JobCreator) {
	es = settings.ResourceRelationshipCheckCronExpr
	jc = func(logger log.Logger, expr string) (cron.Expr, cron.Task, error) {
		task, err := svcskd.NewResourceRelationshipCheckTask(logger,
			opts.ModelClient, opts.K8sConfig)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}

	return
}

func buildTelemetryPeriodicReportJobCreator(opts initOptions) (es settings.Value, jc cron.JobCreator) {
	es = settings.TelemetryPeriodicReportCronExpr
	jc = func(logger log.Logger, expr string) (cron.Expr, cron.Task, error) {
		task, err := telemetryskd.NewPeriodicReportTask(logger, opts.ModelClient)
		if err != nil {
			return nil, nil, err
		}

		return cron.AwaitedExpr(expr), task, nil
	}

	return
}

func buildTokenDeploymentExpireCleanJobCreator(opts initOptions) (es settings.Value, jc cron.JobCreator) {
	es = settings.TokenDeploymentExpiredCleanCronExpr
	jc = func(logger log.Logger, expr string) (cron.Expr, cron.Task, error) {
		task, err := tokenskd.NewDeploymentExpiredCleanTask(logger, opts.ModelClient)
		if err != nil {
			return nil, nil, err
		}

		return cron.ImmediateExpr(expr), task, nil
	}

	return
}
