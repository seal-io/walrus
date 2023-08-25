package cron

import (
	"context"
	"errors"
	"fmt"
	"sync"

	settingbus "github.com/seal-io/walrus/pkg/bus/setting"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/cron"
	"github.com/seal-io/walrus/utils/log"
)

type (
	// JobCreator is the creator for creating {cron.Expr, cron.Task} tuple.
	JobCreator func(logger log.Logger, expr string) (Expr, Task, error)

	// JobExpressionName holds the name of the cron.Expr.
	JobExpressionName = string

	// JobCreators holds JobCreator with its expression name.
	JobCreators map[JobExpressionName]JobCreator
)

var (
	js = JobCreators{}
	o  sync.Once
)

// Register executes all job creators and
// schedules the returning task with the returning expression.
func Register(ctx context.Context, mc *model.Client, cs JobCreators) (err error) {
	err = errors.New("not allowed duplicated registering")

	o.Do(func() {
		for n, c := range cs {
			js[n] = c
		}

		err = doRegister(ctx, mc)
	})

	return
}

func doRegister(ctx context.Context, mc *model.Client) error {
	logger := log.WithName("task")

	// Create locker.
	locker := NewLocker(logger, mc)

	// NB(thxCode): don't stop the core cron scheduler.
	err := cron.Start(ctx, locker)
	if err != nil {
		return err
	}

	for n, c := range js {
		if n == "" || c == nil {
			continue
		}

		s := settings.Index(n)
		if s == nil {
			continue
		}

		// Get cron expr of the job from global model client.
		var v string

		v, err = s.Value(ctx, mc)
		if err != nil {
			return fmt.Errorf("error gettting job cron expr: %w", err)
		}

		ce, ct, err := c(logger.WithValues("createdBy", n), v)
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
func Sync(ctx context.Context, m settingbus.BusMessage) error {
	logger := log.WithName("task")

	type job struct {
		Name string
		Expr Expr
		Task Task
	}

	var jobs []job

	for i := 0; i < len(m.Refers); i++ {
		if m.Refers[i] == nil || m.Refers[i].Name == "" {
			continue
		}

		n := m.Refers[i].Name

		c, exist := js[n]
		if !exist {
			continue
		}

		s := settings.Index(n)
		if s == nil {
			continue
		}

		// Get cron expr of the job from transactional model client.
		v, err := s.Value(ctx, m.TransactionalModelClient)
		if err != nil {
			return fmt.Errorf("error gettting job cron expr: %w", err)
		}

		j := job{Name: n}

		j.Expr, j.Task, err = c(logger.WithValues("createdBy", n), v)
		if err != nil {
			return fmt.Errorf("error creating job for %s: %w", n, err)
		}

		jobs = append(jobs, j)
	}

	for i := 0; i < len(jobs); i++ {
		j := jobs[i]

		err := cron.Schedule(j.Name, j.Expr, j.Task)
		if err != nil {
			// NB(thxCode): raising error cannot roll back successfully scheduled job in the same for-loop,
			// so just warn out here.
			logger.Errorf("error scheduling job for: %v", j.Name, err)
		}
		// TODO(thxCode): support rolling back successfully scheduled job.
	}

	return nil
}
