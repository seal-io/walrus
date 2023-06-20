package cron

import "github.com/prometheus/client_golang/prometheus"

var _statsCollector = newStatsCollector()

func NewStatsCollector() prometheus.Collector {
	return _statsCollector
}

func newStatsCollector() *statsCollector {
	ns := "cron"

	return &statsCollector{
		schedulingJobs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: ns,
				Name:      "scheduling_jobs",
				Help:      "The number of scheduling jobs.",
			},
			[]string{"job"},
		),
		scheduledTasks: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: ns,
				Name:      "scheduled_tasks_total",
				Help:      "The total number of scheduling tasks.",
			},
			[]string{"job"},
		),
		processingTasks: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: ns,
				Name:      "processing_tasks",
				Help:      "The number of tasks processing.",
			},
			[]string{"job"},
		),
		succeededTasks: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: ns,
				Name:      "succeeded_tasks_total",
				Help:      "The total number of tasks successful completed.",
			},
			[]string{"job"},
		),
		failedTasks: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: ns,
				Name:      "failed_tasks_total",
				Help:      "The total number of tasks unsuccessful completed.",
			},
			[]string{"job"},
		),
		taskDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: ns,
				Name:      "task_duration_seconds",
				Help:      "The task consumption distribution in seconds.",
				Buckets: []float64{
					15,
					30,
					60,
					90,
					120,
					180,
					240,
					300,
					600,
					900,
				},
			},
			[]string{"job"},
		),
	}
}

type statsCollector struct {
	schedulingJobs *prometheus.GaugeVec

	scheduledTasks  *prometheus.CounterVec
	processingTasks *prometheus.GaugeVec
	succeededTasks  *prometheus.CounterVec
	failedTasks     *prometheus.CounterVec
	taskDurations   *prometheus.HistogramVec
}

func (c *statsCollector) Describe(ch chan<- *prometheus.Desc) {
	c.schedulingJobs.Describe(ch)
	c.scheduledTasks.Describe(ch)
	c.processingTasks.Describe(ch)
	c.succeededTasks.Describe(ch)
	c.failedTasks.Describe(ch)
	c.taskDurations.Describe(ch)
}

func (c *statsCollector) Collect(ch chan<- prometheus.Metric) {
	c.schedulingJobs.Collect(ch)
	c.scheduledTasks.Collect(ch)
	c.processingTasks.Collect(ch)
	c.succeededTasks.Collect(ch)
	c.failedTasks.Collect(ch)
	c.taskDurations.Collect(ch)
}
