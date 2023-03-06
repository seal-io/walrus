package server

import (
	"context"
	"fmt"

	costskd "github.com/seal-io/seal/pkg/costs/scheduler"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/cron"
	"github.com/seal-io/seal/utils/log"
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

type (
	// jobName is the name of job.
	jobName = string

	// jobCreator is the creator for creating {cron.Expr, cron.Task} tuple,
	// the life of given context.Context ends by this creation,
	// do not use the long-term processing with this context.Context.
	jobCreator func(ctx context.Context, name, expr string) (cron.Expr, cron.Task, error)

	// jobCreators holds all job creators.
	jobCreators struct {
		creators    map[jobName]jobCreator
		modelClient model.ClientSet
	}
)

// Register calls all register functions for the job factory.
func (js jobCreators) Register(ctx context.Context) error {
	for n, c := range js.creators {
		if c == nil {
			continue
		}

		var s = settings.Index(n)
		if s == nil {
			continue
		}
		// get cron expr of the job from global model client.
		var v, err = s.Value(ctx, js.modelClient)
		if err != nil {
			return fmt.Errorf("error gettting job cron expr: %w", err)
		}

		ce, ct, err := c(ctx, n, v)
		if err != nil {
			return fmt.Errorf("error creating %s job: %w", n, err)
		}
		err = cron.Schedule(n, ce, ct)
		if err != nil {
			return fmt.Errorf("error scheduling %s job: %w", n, err)
		}
	}
	return nil
}

// Sync observes the cron expr setting changes and re-register jobs.
func (js jobCreators) Sync(ctx context.Context, m settings.BusMessage) error {
	var logger = log.WithName("jobs")

	type job struct {
		Name string
		Expr cron.Expr
		Task cron.Task
	}

	var jobs []job
	for i := 0; i < len(m.Refer); i++ {
		if m.Refer[i] == nil {
			continue
		}

		var n = m.Refer[i].Name
		var c, exist = js.creators[n]
		if !exist {
			continue
		}

		var s = settings.Index(n)
		if s == nil {
			continue
		}
		// get cron expr of the job from transactional model client.
		var v, err = s.Value(ctx, m.ModelClient)
		if err != nil {
			return fmt.Errorf("error gettting job cron expr: %w", err)
		}

		var j = job{Name: n}
		j.Expr, j.Task, err = c(ctx, n, v)
		if err != nil {
			return fmt.Errorf("error creating %s job: %w", n, err)
		}
		jobs = append(jobs, j)
	}

	for i := 0; i < len(jobs); i++ {
		var j = jobs[i]
		var err = cron.Schedule(j.Name, j.Expr, j.Task)
		if err != nil {
			// NB(thxCode): raising error cannot roll back successfully scheduled job in the same for-loop,
			// so just warn out here.
			logger.Errorf("error scheduling %s job: %v", j.Name, err)
		}
		// TODO(thxCode): support rolling back successfully scheduled job.
	}
	return nil
}
