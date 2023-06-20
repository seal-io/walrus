package runtime

import (
	"github.com/prometheus/client_golang/prometheus"
)

var _statsCollector = newStatsCollector()

func NewStatsCollector() prometheus.Collector {
	return _statsCollector
}

func newStatsCollector() *statsCollector {
	ns := "api"

	return &statsCollector{
		requestInflight: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: ns,
				Name:      "request_inflight",
				Help:      "The number of inflight request.",
			},
			[]string{"proto", "path", "method"},
		),
		requestCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: ns,
				Name:      "request_total",
				Help:      "The total number of requests.",
			},
			[]string{"proto", "path", "method", "code"},
		),
		requestDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: ns,
				Name:      "request_duration_seconds",
				Help:      "The response latency distribution in seconds.",
				Buckets: []float64{
					0.005,
					0.025,
					0.05,
					0.1,
					0.2,
					0.4,
					0.6,
					0.8,
					1.0,
					1.25,
					1.5,
					2,
					3,
					4,
					5,
					6,
					8,
					10,
					15,
					20,
					30,
					45,
					60,
				},
			},
			[]string{"proto", "path", "method", "code"},
		),
		requestSizes: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: ns,
				Name:      "request_sizes",
				Help:      "The request size distribution in bytes.",
				Buckets:   prometheus.ExponentialBuckets(128, 2.0, 15), // 128B, 256B, ..., 2M.
			},
			[]string{"proto", "path", "method"},
		),
		responseSizes: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: ns,
				Name:      "response_sizes",
				Help:      "The response size distribution in bytes.",
				Buckets:   prometheus.ExponentialBuckets(128, 2.0, 15), // 128B, 256B, ..., 2M.
			},
			[]string{"proto", "path", "method"},
		),
	}
}

type statsCollector struct {
	requestInflight  *prometheus.GaugeVec
	requestCounter   *prometheus.CounterVec
	requestDurations *prometheus.HistogramVec
	requestSizes     *prometheus.HistogramVec
	responseSizes    *prometheus.HistogramVec
}

func (c *statsCollector) Describe(ch chan<- *prometheus.Desc) {
	c.requestInflight.Describe(ch)
	c.requestCounter.Describe(ch)
	c.requestDurations.Describe(ch)
	c.requestSizes.Describe(ch)
	c.responseSizes.Describe(ch)
}

func (c *statsCollector) Collect(ch chan<- prometheus.Metric) {
	c.requestInflight.Collect(ch)
	c.requestCounter.Collect(ch)
	c.requestDurations.Collect(ch)
	c.requestSizes.Collect(ch)
	c.responseSizes.Collect(ch)
}
