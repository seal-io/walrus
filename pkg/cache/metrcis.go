package cache

import "github.com/prometheus/client_golang/prometheus"

func NewStatsCollectorWith(drv Driver) prometheus.Collector {
	fqName := func(name string) string {
		return "go_cache_remote_" + name
	}

	return &statsCollector{
		drv: drv,
		maxOpenConnections: prometheus.NewDesc(
			fqName("max_open_connections"),
			"The maximum number of open connections to the remote cache.",
			nil, nil,
		),
		idleConnections: prometheus.NewDesc(
			fqName("idle_connections"),
			"The number of idle connections.",
			nil, nil,
		),
		newOpenCount: prometheus.NewDesc(
			fqName("new_open_total"),
			"The total number of connections to newly create.",
			nil, nil,
		),
		timeoutCount: prometheus.NewDesc(
			fqName("timeout_total"),
			"The total number of getting a connection timeout times.",
			nil, nil,
		),
		closedCount: prometheus.NewDesc(
			fqName("closed_total"),
			"The total number of connections closed.",
			nil, nil,
		),
	}
}

type statsCollector struct {
	drv Driver

	maxOpenConnections *prometheus.Desc
	idleConnections    *prometheus.Desc

	newOpenCount *prometheus.Desc
	timeoutCount *prometheus.Desc
	closedCount  *prometheus.Desc
}

func (c *statsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpenConnections
	ch <- c.idleConnections
	ch <- c.newOpenCount
	ch <- c.timeoutCount
	ch <- c.closedCount
}

func (c *statsCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.drv.Stats()
	ch <- prometheus.MustNewConstMetric(c.maxOpenConnections, prometheus.GaugeValue, float64(stats.MaxOpenConnections))
	ch <- prometheus.MustNewConstMetric(c.idleConnections, prometheus.GaugeValue, float64(stats.IdleConnections))
	ch <- prometheus.MustNewConstMetric(c.newOpenCount, prometheus.CounterValue, float64(stats.NewOpenCount))
	ch <- prometheus.MustNewConstMetric(c.timeoutCount, prometheus.CounterValue, float64(stats.TimeoutCount))
	ch <- prometheus.MustNewConstMetric(c.closedCount, prometheus.CounterValue, float64(stats.ClosedCount))
}
