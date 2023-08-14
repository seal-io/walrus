package server

import (
	"context"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/cache"
	"github.com/seal-io/seal/pkg/database"
	"github.com/seal-io/seal/pkg/metric"
	"github.com/seal-io/seal/utils/cron"
	"github.com/seal-io/seal/utils/gopool"
)

// registerMetricCollectors registers the metric collectors into the global metric registry.
func (r *Server) registerMetricCollectors(ctx context.Context, opts initOptions) error {
	cs := metric.Collectors{
		database.NewStatsCollectorWith(opts.DatabaseDriver),
		gopool.NewStatsCollector(),
		cron.NewStatsCollector(),
		runtime.NewStatsCollector(),
	}

	if opts.CacheDriver != nil {
		cs = append(cs, cache.NewStatsCollectorWith(opts.CacheDriver))
	}

	return metric.Register(ctx, cs)
}
