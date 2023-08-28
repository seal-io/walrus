package cron

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	cronlib "github.com/robfig/cron/v3"

	"github.com/seal-io/walrus/utils/log"
)

// Task defines the interface to hold the job executing main logic.
type Task interface {
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

// ParseCronExpr parses the given string as cronlib.Schedule,
// returns nil in none strict mode if passing blank string.
func ParseCronExpr(ce string, strict bool) (cronlib.Schedule, error) {
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
	_, err := ParseCronExpr(ce, true)
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
	Start(ctx context.Context, lc Locker) error
	// Stop stops scheduling.
	Stop() error
	// State returns an indexer that index the running job status by its name.
	State() map[string]JobStatus
}

type scheduler struct {
	c context.Context
	s *gocron.Scheduler
}

type timeoutTask struct {
	timeout time.Duration
	jobName string
	task    Task
}

func (in timeoutTask) Process(ctx context.Context, args ...any) error {
	logger := log.WithName("task").WithValues("createBy", in.jobName)

	ctx, cancel := context.WithTimeout(ctx, in.timeout)
	defer cancel()

	// Record scheduled task.
	_statsCollector.scheduledTasks.
		WithLabelValues(in.jobName).
		Inc()

	// Record processing task.
	_statsCollector.processingTasks.
		WithLabelValues(in.jobName).
		Inc()

	defer func() {
		_statsCollector.processingTasks.
			WithLabelValues(in.jobName).
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
		WithLabelValues(in.jobName).
		Observe(time.Since(start).Seconds())

	if err != nil {
		// Record failed task.
		_statsCollector.failedTasks.
			WithLabelValues(in.jobName).
			Inc()
		logger.Errorf("error executing task: %v", err)
	} else {
		// Record succeeded task.
		_statsCollector.succeededTasks.
			WithLabelValues(in.jobName).
			Inc()
		logger.Debugf("executed task")
	}

	// NB(thxCode): always return nil as there is no way to restart the job at present.
	return nil
}

type emptyVariadicList struct{}

func (in *scheduler) Schedule(jobName string, cron Expr, task Task, taskArgs ...any) (err error) {
	ce := cron.String()

	ceParsed, err := ParseCronExpr(ce, false)
	if err != nil {
		return
	}

	err = in.s.RemoveByTag(jobName)
	if err != nil && !errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		return
	}

	// Record scheduled job.
	_statsCollector.schedulingJobs.
		WithLabelValues(jobName).
		Set(0)

	defer func() {
		if err != nil {
			return
		}

		_statsCollector.schedulingJobs.
			WithLabelValues(jobName).
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
		jobName: jobName,
		task:    task,
	}

	var variadicArgs any
	if len(taskArgs) == 0 {
		variadicArgs = emptyVariadicList{}
	} else {
		variadicArgs = taskArgs
	}

	s := in.s.CronWithSeconds(ce).
		Tag(jobName).
		Name(jobName)
	if cron.runImmediately() {
		s.StartImmediately()
	}
	_, err = s.Do(tt.Process, in.c, variadicArgs)

	return err
}

func (in *scheduler) Start(ctx context.Context, lc Locker) error {
	s := gocron.NewScheduler(time.Now().Location())
	s.WaitForScheduleAll()
	s.TagsUnique()
	s.StartAsync()
	s.WithDistributedLocker(lc)
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

func (in *scheduler) State() map[string]JobStatus {
	var (
		sj = in.s.Jobs()
		js = make(map[string]JobStatus, len(sj))
	)

	for i := range sj {
		if !sj[i].IsRunning() {
			continue
		}

		js[sj[i].Tags()[0]] = JobStatus{
			LastRun: sj[i].LastRun(),
			NextRun: sj[i].NextRun(),
		}
	}

	return js
}

// JobStatus holds the status of a running job.
type JobStatus struct {
	// LastRun observes the time job last run.
	LastRun time.Time
	// NextRun observes the next time job should be schedule.
	NextRun time.Time
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
func Start(ctx context.Context, lc Locker) error {
	return globalScheduler.Start(ctx, lc)
}

// Stop stops scheduling.
func Stop() error {
	return globalScheduler.Stop()
}

// State returns an indexer that index the running job status by its name.
func State() map[string]JobStatus {
	return globalScheduler.State()
}
