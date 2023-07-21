package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/seal-io/seal/utils/log"
)

const (
	DialectRedis        = "redis"
	DialectRedisCluster = "rediss"
)

func GetDriverAndName(cacheSourceAddress string) (csd, csn string, err error) {
	if cacheSourceAddress == "" {
		err = errors.New("blank cache source address")
		return
	}

	switch {
	case strings.HasPrefix(cacheSourceAddress, "redis://"):
		csd = DialectRedis
		csn, err = clipRedisAddress(cacheSourceAddress)
	case strings.HasPrefix(cacheSourceAddress, "rediss://"):
		csd = DialectRedisCluster
		csn, err = clipRedisAddress(cacheSourceAddress)
	}

	if csd == "" {
		err = errors.New("cannot recognize driver from cache source address")
	}

	return
}

func LoadDriver(cacheSourceAddress string) (drvDialect string, drv Driver, err error) {
	drvDialect, drvSource, err := GetDriverAndName(cacheSourceAddress)
	if err != nil {
		return
	}

	switch drvDialect {
	case DialectRedis, DialectRedisCluster:
		var in any
		if drvDialect == DialectRedis {
			// Redis single.
			in, err = redis.ParseURL(drvSource)
		} else {
			// Redis cluster.
			in, err = redis.ParseClusterURL(drvSource)
		}

		if err != nil {
			return
		}

		drv, err = newRedisDriver(in)
		if err != nil {
			return
		}
	}

	return
}

func Wait(ctx context.Context, drv Driver) (err error) {
	var lastErr error

	err = wait.PollUntilContextCancel(ctx, 2*time.Second, true,
		func(ctx context.Context) (bool, error) {
			lastErr = IsConnected(ctx, drv)
			if lastErr != nil {
				log.Warnf("waiting for database to be ready: %v", lastErr)
			}

			return lastErr == nil, ctx.Err()
		},
	)
	if err != nil && lastErr != nil {
		err = lastErr // Use last error to overwrite context error while existed.
	}

	return
}

func IsConnected(ctx context.Context, drv Driver) error {
	return drv.PingContext(ctx)
}
