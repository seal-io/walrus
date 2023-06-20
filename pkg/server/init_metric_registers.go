package server

import (
	"context"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/metric"
	"github.com/seal-io/seal/pkg/rds"
	"github.com/seal-io/seal/utils/cron"
	"github.com/seal-io/seal/utils/gopool"
)

func (r *Server) initMetrics(ctx context.Context, opts initOptions) error {
	cs := metric.Collectors{
		rds.NewStatsCollectorWith(opts.RdsDriver),
		gopool.NewStatsCollector(),
		cron.NewStatsCollector(),
		runtime.NewStatsCollector(),
	}

	return metric.Register(ctx, cs)
}
