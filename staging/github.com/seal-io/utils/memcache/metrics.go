package memcache

import (
	"github.com/allegro/bigcache/v3"
	"github.com/prometheus/client_golang/prometheus"
)

// NewStatsCollector returns a new prometheus.Collector for global cache.
func NewStatsCollector() prometheus.Collector {
	fqName := func(name string) string {
		return "memcache_" + name
	}

	return &statsCollector{
		entries: prometheus.NewDesc(
			fqName("entries"),
			"The number of entries in the cache.",
			nil, nil,
		),
		capacity: prometheus.NewDesc(
			fqName("capacity"),
			"The store bytes amount in the cache.",
			nil, nil,
		),
		hitsTotal: prometheus.NewDesc(
			fqName("hits_total"),
			"The total number of cache hits(successfully found).",
			nil, nil,
		),
		missesTotal: prometheus.NewDesc(
			fqName("misses_total"),
			"The total number of cache misses(not found).",
			nil, nil,
		),
		deleteHitsTotal: prometheus.NewDesc(
			fqName("delete_hits_total"),
			"The total number of cache delete hits(successfully deleted).",
			nil, nil,
		),
		deleteMissesTotal: prometheus.NewDesc(
			fqName("delete_misses_total"),
			"The total number of cache delete misses(not deleted).",
			nil, nil,
		),
		collisionsTotal: prometheus.NewDesc(
			fqName("collisions_total"),
			"The total number of cache key-collisions.",
			nil, nil,
		),
	}
}

type statsCollector struct {
	entries  *prometheus.Desc
	capacity *prometheus.Desc

	hitsTotal         *prometheus.Desc
	missesTotal       *prometheus.Desc
	deleteHitsTotal   *prometheus.Desc
	deleteMissesTotal *prometheus.Desc
	collisionsTotal   *prometheus.Desc
}

func (c *statsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.entries
	ch <- c.capacity
	ch <- c.hitsTotal
	ch <- c.missesTotal
	ch <- c.deleteHitsTotal
	ch <- c.deleteMissesTotal
	ch <- c.collisionsTotal
}

func (c *statsCollector) Collect(ch chan<- prometheus.Metric) {
	bc := gc.(Underlay[*bigcache.BigCache]).Underlay()

	var (
		entries  = bc.Len()
		capacity = bc.Capacity()
		stats    = bc.Stats()
	)

	ch <- prometheus.MustNewConstMetric(c.entries, prometheus.GaugeValue, float64(entries))
	ch <- prometheus.MustNewConstMetric(c.capacity, prometheus.GaugeValue, float64(capacity))
	ch <- prometheus.MustNewConstMetric(c.hitsTotal, prometheus.CounterValue, float64(stats.Hits))
	ch <- prometheus.MustNewConstMetric(c.missesTotal, prometheus.CounterValue, float64(stats.Misses))
	ch <- prometheus.MustNewConstMetric(c.deleteHitsTotal, prometheus.CounterValue, float64(stats.DelHits))
	ch <- prometheus.MustNewConstMetric(c.deleteMissesTotal, prometheus.CounterValue, float64(stats.DelMisses))
	ch <- prometheus.MustNewConstMetric(c.collisionsTotal, prometheus.CounterValue, float64(stats.Collisions))
}
