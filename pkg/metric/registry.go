package metric

import (
	"context"
	"errors"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	kmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

// Collectors holds the list of prometheus.Collector.
type Collectors = []prometheus.Collector

var (
	reg = kmetrics.Registry
	o   sync.Once
)

// Register registers all metric collectors.
func Register(ctx context.Context, cs Collectors) (err error) {
	err = errors.New("not allowed duplicated registering")

	o.Do(func() {
		err = reg.Register(collectors.NewBuildInfoCollector())
		if err != nil {
			return
		}

		for i := range cs {
			err = reg.Register(cs[i])
			if err != nil {
				break
			}
		}
	})

	return
}
