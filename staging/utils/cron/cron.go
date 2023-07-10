package cron

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	cronlib "github.com/robfig/cron/v3"

	"github.com/seal-io/seal/utils/log"
)

// Task defines the interface to hold the job executing main logic.
type Task interface {
	// Name return name for task.
	Name() string
	// Process executes the task main logic.
	Process(ctx context.Context, args ...any) error
}

// TaskFunc implements the Task interface to provide a convenient wrapper.
type TaskFunc func(ctx context.Context, args ...any) error

// Expr holds the definition of cron expression.
type Expr interface {
	fmt.Stringer

	runImmediately() bool
}

var (
	globalScheduler = New()
	cronParser      = cronlib.NewParser(
		cronlib.Second | cronlib.Minute | cronlib.Hour | cronlib.Dom | cronlib.Month | cronlib.Dow | cronlib.Descriptor)
)

func init() {
	gocron.SetPanicHandler(func(jobName string, recoverData any) {
		log.WithName("task").Errorf("panic in job: %s, recover data: %v", jobName, recoverData)
	})
}

func (fn TaskFunc) Process(ctx context.Context, args ...any) error {
	if fn == nil {
		return nil
	}

	return fn(ctx, args...)
}

// parseCronExpr parses the given string as cronlib.Schedule,
// returns nil in none strict mode if passing blank string.
func parseCronExpr(ce string, strict bool) (cronlib.Schedule, error) {
	if !strict && ce == "" {
		return nil, nil
	}

	cron, err := cronParser.Parse(ce)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression: %w", err)
	}

	return cron, nil
}

// ValidateCronExpr returns error if the given Expr is invalid.
func ValidateCronExpr(ce string) error {
	_, err := parseCronExpr(ce, true)
	if err != nil {
		return err
	}

	return nil
}

type scheduleCronExpr struct {
	raw         string
	immediately bool
}

func (in scheduleCronExpr) String() string {
	return in.raw
}

func (in scheduleCronExpr) runImmediately() bool {
	return in.immediately
}

// AwaitedExpr returns an Expr and runs in the next round.
func AwaitedExpr(raw string) Expr {
	return scheduleCronExpr{raw: raw}
}

// ImmediateExpr returns an Expr and runs immediately.
func ImmediateExpr(raw string) Expr {
	return scheduleCronExpr{raw: raw, immediately: true}
}

// Scheduler defines the interface to maintain the simple on-time scheduling logic.
type Scheduler interface {
	// Schedule registers the job with the given name,
	// and schedules it at the given Expr.
	// Remove from scheduler if the given Expr is blank.
	// If the given name job has found,
	// Schedule updates it with the new Expr.
	Schedule(name string, cron Expr, task Task, taskArgs ...any) error
	// Start starts scheduling.
	Start(ctx context.Context) error
	// Stop stops scheduling.
	Stop() error
}

type scheduler struct {
	c context.Context
	s *gocron.Scheduler
}

type timeoutTask struct {
	timeout time.Duration
	name    string
	task    Task
}

func (in timeoutTask) Process(ctx context.Context, args ...any) error {
	logger := log.WithName("task")

	ctx, cancel := context.WithTimeout(ctx, in.timeout)
	defer cancel()

	// Record scheduled task.
	_statsCollector.scheduledTasks.
		WithLabelValues(in.name).
		Inc()

	// Record processing task.
	_statsCollector.processingTasks.
		WithLabelValues(in.name).
		Inc()

	defer func() {
		_statsCollector.processingTasks.
			WithLabelValues(in.name).
			Dec()
	}()

	start := time.Now()

	var err error
	if len(args) == 0 {
		err = in.task.Process(ctx)
	} else {
		t, ok := args[0].([]any)
		if !ok {
			err = in.task.Process(ctx)
		} else {
			err = in.task.Process(ctx, t...)
		}
	}

	// Record task consumption.
	_statsCollector.taskDurations.
		WithLabelValues(in.name).
		Observe(time.Since(start).Seconds())

	if err != nil {
		// Record failed task.
		_statsCollector.failedTasks.
			WithLabelValues(in.name).
			Inc()
		logger.Errorf("error executing task: %s: %v", in.task.Name(), err)
	} else {
		// Record succeeded task.
		_statsCollector.succeededTasks.
			WithLabelValues(in.name).
			Inc()
		logger.Debugf("executed task: %s", in.task.Name())
	}

	// NB(thxCode): always return nil as there is no way to restart the job at present.
	return nil
}

type emptyVariadicList struct{}

func (in *scheduler) Schedule(name string, cron Expr, task Task, taskArgs ...any) (err error) {
	ce := cron.String()

	ceParsed, err := parseCronExpr(ce, false)
	if err != nil {
		return
	}

	err = in.s.RemoveByTag(name)
	if err != nil && !errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		return
	}

	// Record scheduled job.
	_statsCollector.schedulingJobs.
		WithLabelValues(name).
		Set(0)

	defer func() {
		if err != nil {
			return
		}

		_statsCollector.schedulingJobs.
			WithLabelValues(name).
			Set(1)
	}()

	if ceParsed == nil {
		return
	}

	const atLeast = 5 * time.Minute

	var (
		now  = time.Now()
		next = ceParsed.Next(now).Sub(now)
	)

	if next > atLeast {
		next >>= 1
	}

	if next < atLeast {
		next = atLeast
	}
	tt := timeoutTask{
		timeout: next,
		name:    name,
		task:    task,
	}

	var variadicArgs any
	if len(taskArgs) == 0 {
		variadicArgs = emptyVariadicList{}
	} else {
		variadicArgs = taskArgs
	}

	s := in.s.CronWithSeconds(ce).Tag(name)
	if cron.runImmediately() {
		s.StartImmediately()
	}
	_, err = s.Do(tt.Process, in.c, variadicArgs)

	return err
}

func (in *scheduler) Start(ctx context.Context) error {
	s := gocron.NewScheduler(time.Now().Location())
	s.WaitForScheduleAll()
	s.TagsUnique()
	s.StartAsync()
	in.c = ctx
	in.s = s

	return nil
}

func (in *scheduler) Stop() error {
	if in.s != nil {
		in.s.Stop()
	}

	return nil
}

// New returns a new Scheduler.
func New() Scheduler {
	return &scheduler{}
}

// Schedule registers the task with the given name as a job,
// and schedules it at the given Expr.
// Remove from scheduler if the given Expr is blank.
// If the given name task has found,
// Schedule updates it with the new Expr.
func Schedule(name string, cron Expr, task Task, taskArgs ...any) error {
	return globalScheduler.Schedule(name, cron, task, taskArgs...)
}

// MustSchedule likes Schedule, but panic if error found.
func MustSchedule(name string, cron Expr, task Task, taskArgs ...any) {
	err := Schedule(name, cron, task, taskArgs...)
	if err != nil {
		panic(err)
	}
}

// Start starts scheduling.
func Start(ctx context.Context) error {
	return globalScheduler.Start(ctx)
}

// Stop stops scheduling.
func Stop() error {
	return globalScheduler.Stop()
}
