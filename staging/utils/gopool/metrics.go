package gopool

import "github.com/prometheus/client_golang/prometheus"

func NewStatsCollector() prometheus.Collector {
	fqName := func(name string) string {
		return "go_pool_" + name
	}

	return &statsCollector{
		maxWorkers: prometheus.NewDesc(
			fqName("max_workers"),
			"The maximum number of pooling goroutines.",
			nil, nil,
		),
		inUseWorkers: prometheus.NewDesc(
			fqName("in_use_workers"),
			"The number of pooling goroutines in use.",
			nil, nil,
		),
		idleWorkers: prometheus.NewDesc(
			fqName("idle_workers"),
			"The number of idle pooling goroutines.",
			nil, nil,
		),
		maxTasks: prometheus.NewDesc(
			fqName("max_tasks"),
			"The maximum number of queuing tasks.",
			nil, nil,
		),
		submittedTasks: prometheus.NewDesc(
			fqName("submitted_tasks_total"),
			"The total number of tasks submitted.",
			nil, nil,
		),
		waitingTasks: prometheus.NewDesc(
			fqName("waiting_tasks"),
			"The number of tasks waiting for.",
			nil, nil,
		),
		succeededTasks: prometheus.NewDesc(
			fqName("succeeded_tasks_total"),
			"The total number of tasks successful completed.",
			nil, nil,
		),
		failedTasks: prometheus.NewDesc(
			fqName("failed_tasks_total"),
			"The total number of tasks unsuccessful completed.",
			nil, nil,
		),
	}
}

type statsCollector struct {
	maxWorkers   *prometheus.Desc
	inUseWorkers *prometheus.Desc
	idleWorkers  *prometheus.Desc

	maxTasks       *prometheus.Desc
	submittedTasks *prometheus.Desc
	waitingTasks   *prometheus.Desc
	succeededTasks *prometheus.Desc
	failedTasks    *prometheus.Desc
}

func (c *statsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxWorkers
	ch <- c.inUseWorkers
	ch <- c.idleWorkers
	ch <- c.maxTasks
	ch <- c.submittedTasks
	ch <- c.waitingTasks
	ch <- c.succeededTasks
	ch <- c.failedTasks
}

func (c *statsCollector) Collect(ch chan<- prometheus.Metric) {
	var (
		maxWorkers     = gp.MaxWorkers()
		runningWorkers = gp.RunningWorkers()
		idleWorkers    = gp.IdleWorkers()
		inUseWorkers   = runningWorkers - idleWorkers
		maxTasks       = gp.MaxCapacity()
		submittedTasks = gp.SubmittedTasks()
		waitingTasks   = gp.WaitingTasks()
		succeededTasks = gp.SuccessfulTasks()
		failedTasks    = submittedTasks - waitingTasks - succeededTasks
	)

	ch <- prometheus.MustNewConstMetric(c.maxWorkers, prometheus.GaugeValue, float64(maxWorkers))
	ch <- prometheus.MustNewConstMetric(c.inUseWorkers, prometheus.GaugeValue, float64(inUseWorkers))
	ch <- prometheus.MustNewConstMetric(c.idleWorkers, prometheus.GaugeValue, float64(idleWorkers))
	ch <- prometheus.MustNewConstMetric(c.maxTasks, prometheus.GaugeValue, float64(maxTasks))
	ch <- prometheus.MustNewConstMetric(c.submittedTasks, prometheus.CounterValue, float64(submittedTasks))
	ch <- prometheus.MustNewConstMetric(c.waitingTasks, prometheus.GaugeValue, float64(waitingTasks))
	ch <- prometheus.MustNewConstMetric(c.succeededTasks, prometheus.CounterValue, float64(succeededTasks))
	ch <- prometheus.MustNewConstMetric(c.failedTasks, prometheus.CounterValue, float64(failedTasks))
}
