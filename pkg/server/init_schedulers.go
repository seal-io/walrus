package server

import (
	"context"

	costskd "github.com/seal-io/seal/pkg/costs/scheduler"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/cron"
)

func (r *Server) runCronJobs(ctx context.Context, opts initOptions) error {
	defer func() {
		_ = cron.Stop()
	}()
	var err = cron.Start(ctx)
	if err != nil {
		return err
	}

	var candidateRegisters = []func(context.Context, model.ClientSet) error{
		r.registerCostCollectTask,
		r.registerCostToolsCheckTask,
	}

	for _, register := range candidateRegisters {
		err = register(ctx, opts.ModelClient)
		if err != nil {
			return err
		}
	}

	err = settings.AddSubscriber("cron-expression", r.SyncCronExprFromSetting)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (r *Server) registerCostCollectTask(ctx context.Context, client model.ClientSet) error {
	var name = settings.CostCollectCronExpr.Name()
	var value, err = settings.CostCollectCronExpr.Value(ctx, client)
	if err != nil {
		return err
	}

	task, err := costskd.NewCostSyncTask(client)
	if err != nil {
		return err
	}
	return cron.Schedule(name, cron.ImmediateExpr(value), task)
}

func (r *Server) registerCostToolsCheckTask(ctx context.Context, client model.ClientSet) error {
	var name = settings.CostToolsCheckCronExpr.Name()
	var value, err = settings.CostToolsCheckCronExpr.Value(ctx, client)
	if err != nil {
		return err
	}

	task, err := costskd.NewToolsCheckTask(client)
	if err != nil {
		return err
	}
	return cron.Schedule(name, cron.ImmediateExpr(value), task)
}

// SyncCronExprFromSetting observes the cron expr change and update the cronJobs.
func (r *Server) SyncCronExprFromSetting(ctx context.Context, m settings.BusMessage) error {
	for i := 0; i < len(m.Refer); i++ {
		var err error
		switch m.Refer[i].Name {
		default:
			continue
		case settings.CostCollectCronExpr.Name():
			err = r.registerCostCollectTask(ctx, m.ModelClient)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
