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
	// Process executes the task main logic.
	Process(ctx context.Context, args ...interface{}) error
}

// TaskFunc implements the Task interface to provide a convenient wrapper.
type TaskFunc func(ctx context.Context, args ...interface{}) error

func (fn TaskFunc) Process(ctx context.Context, args ...interface{}) error {
	if fn == nil {
		return nil
	}
	return fn(ctx, args...)
}

var cronParser = cronlib.NewParser(
	cronlib.Second | cronlib.Minute | cronlib.Hour | cronlib.Dom | cronlib.Month | cronlib.Dow | cronlib.Descriptor)

// parseCronExpr parses the given string as cronlib.Schedule,
// returns nil in none strict mode if passing blank string.
func parseCronExpr(ce string, strict bool) (cronlib.Schedule, error) {
	if !strict && ce == "" {
		return nil, nil
	}
	var cron, err = cronParser.Parse(ce)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression: %w", err)
	}
	return cron, nil
}

// ValidateCronExpr returns error if the given Expr is invalid.
func ValidateCronExpr(ce string) error {
	var _, err = parseCronExpr(ce, true)
	if err != nil {
		return err
	}
	return nil
}

// Expr holds the definition of cron expression.
type Expr interface {
	fmt.Stringer

	runImmediately() bool
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
	Schedule(name string, cron Expr, task Task, taskArgs ...interface{}) error
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

func (in timeoutTask) Process(ctx context.Context, args ...interface{}) error {
	var logger = log.WithName("cronjobs")
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("panic observing of %s task: %v", in.name, r)
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, in.timeout)
	defer cancel()

	var err error
	if len(args) == 0 {
		err = in.task.Process(ctx)
	} else {
		var t, ok = args[0].([]interface{})
		if !ok {
			err = in.task.Process(ctx)
		} else {
			err = in.task.Process(ctx, t...)
		}
	}
	if err != nil {
		logger.Errorf("error executing %s task: %v", in.name, err)
	} else {
		logger.Debugf("executed %s task", in.name)
	}

	// NB(thxCode): always return nil as there is no way to restart the job at present.
	return nil
}

type emptyVariadicList struct{}

func (in *scheduler) Schedule(name string, cron Expr, task Task, taskArgs ...interface{}) (err error) {
	var ce = cron.String()
	ceParsed, err := parseCronExpr(ce, false)
	if err != nil {
		return
	}

	err = in.s.RemoveByTag(name)
	if err != nil && !errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		return
	}
	if ceParsed == nil {
		return
	}

	const atLeast = 5 * time.Minute
	var now = time.Now()
	var next = ceParsed.Next(now).Sub(now)
	if next > atLeast {
		next >>= 1
	}
	if next < atLeast {
		next = atLeast
	}
	var tt = timeoutTask{
		timeout: next,
		name:    name,
		task:    task,
	}
	var variadicArgs interface{}
	if len(taskArgs) == 0 {
		variadicArgs = emptyVariadicList{}
	} else {
		variadicArgs = taskArgs
	}

	var s = in.s.CronWithSeconds(ce).Tag(name)
	if cron.runImmediately() {
		s.StartImmediately()
	}
	_, err = s.Do(tt.Process, in.c, variadicArgs)
	return
}

func (in *scheduler) Start(ctx context.Context) error {
	var s = gocron.NewScheduler(time.Now().Location())
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

var globalScheduler = New()

func init() {
	gocron.SetPanicHandler(func(name string, r interface{}) {
		log.WithName("cronjobs").Errorf("panic observing of %s task: %v", name, r)
	})
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
func Schedule(name string, cron Expr, task Task, taskArgs ...interface{}) error {
	return globalScheduler.Schedule(name, cron, task, taskArgs...)
}

// MustSchedule likes Schedule, but panic if error found.
func MustSchedule(name string, cron Expr, task Task, taskArgs ...interface{}) {
	var err = Schedule(name, cron, task, taskArgs...)
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
