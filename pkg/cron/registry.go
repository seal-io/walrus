package cron

import (
	"context"
	"errors"
	"fmt"
	"sync"

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
