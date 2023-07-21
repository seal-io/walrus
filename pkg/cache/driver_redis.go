package cache

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/seal-io/seal/utils/log"
)

// clipRedisAddress tries to normalize the given redis address,
// removes options that disallowing customized.
func clipRedisAddress(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}
	q := u.Query()
	q.Del("min_idle_conns")
	q.Del("max_idle_conns")
	q.Del("conn_max_lifetime")
	q.Del("conn_max_idle_time")
	q.Del("pool_fifo")
	q.Del("pool_size")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func newRedisDriver(v any) (drv Driver, err error) {
	switch opts := v.(type) {
	default:
		err = errors.New("unknown options type")
		return
	case *redis.Options:
		drv = &redisDriver{singleOpts: opts}
	case *redis.ClusterOptions:
		drv = &redisDriver{clusterOpts: opts}
	}
	logger := log.WithName("cache").WithName("redis")
	redis.SetLogger(redisLogger(logger.Debugf))

	return
}

// redisDriver implement Driver interface to configure the pool setting before getting client.
type redisDriver struct {
	sync.Mutex

	// NB(thxCode): go-redis doesn't provide a method to migrate redis.Options/redis.ClusterOptions to
	// redis.UniversalOptions, a simple approach is to maintain two kind of options with a decision logic,
	// it's ugly but work.
	singleOpts  *redis.Options
	clusterOpts *redis.ClusterOptions
	drv         string
	cli         redis.UniversalClient
}

func (r *redisDriver) getClient() redis.UniversalClient {
	r.Lock()
	defer r.Unlock()

	if r.cli == nil {
		// Singleton pattern.
		if r.singleOpts != nil {
			r.singleOpts.PoolFIFO = true
			r.singleOpts.PoolSize = r.singleOpts.MaxIdleConns
			r.drv = DialectRedis
			r.cli = redis.NewClient(r.singleOpts)
		} else {
			r.clusterOpts.PoolFIFO = true
			r.clusterOpts.PoolSize = r.clusterOpts.MaxIdleConns
			r.drv = DialectRedisCluster
			r.cli = redis.NewClusterClient(r.clusterOpts)
		}
	}

	return r.cli
}

func (r *redisDriver) SetMaxIdleConns(n int) {
	if r.singleOpts != nil {
		r.singleOpts.MinIdleConns = n
	} else {
		r.clusterOpts.MinIdleConns = n
	}
}

func (r *redisDriver) SetMaxOpenConns(n int) {
	if r.singleOpts != nil {
		r.singleOpts.MaxIdleConns = n
	} else {
		r.clusterOpts.MaxIdleConns = n
	}
}

func (r *redisDriver) SetConnMaxLifetime(d time.Duration) {
	if r.singleOpts != nil {
		r.singleOpts.ConnMaxLifetime = d
	} else {
		r.clusterOpts.ConnMaxLifetime = d
	}
}

func (r *redisDriver) SetConnMaxIdleTime(d time.Duration) {
	if r.singleOpts != nil {
		r.singleOpts.ConnMaxIdleTime = d
	} else {
		r.clusterOpts.ConnMaxIdleTime = d
	}
}

func (r *redisDriver) PingContext(ctx context.Context) error {
	return r.getClient().Ping(ctx).Err()
}

func (r *redisDriver) Underlay(ctx context.Context) (string, any, error) {
	return r.drv, r.getClient(), nil
}

func (r *redisDriver) Stats() DriverStats {
	s := r.getClient().PoolStats()

	return DriverStats{
		MaxOpenConnections: int64(s.TotalConns),
		IdleConnections:    int64(s.IdleConns),
		NewOpenCount:       int64(s.Misses),
		TimeoutCount:       int64(s.Timeouts),
		ClosedCount:        int64(s.StaleConns),
	}
}

// redisLogger implements the redis internal.Logging.
type redisLogger func(string, ...interface{})

func (l redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	l(format, v...)
}
