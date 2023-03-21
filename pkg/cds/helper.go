package cds

import (
	"context"
	"errors"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // db = mysql
	_ "github.com/lib/pq"              // db = postgres
	_ "github.com/mattn/go-sqlite3"    // db = sqlite3
	"github.com/redis/go-redis/v9"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/seal-io/seal/utils/log"
)

func GetDriverAndName(cacheSourceAddress string) (csd, csn string, err error) {
	if cacheSourceAddress == "" {
		err = errors.New("blank cache source address")
		return
	}
	switch {
	case strings.HasPrefix(cacheSourceAddress, "redis://"):
		csd = "redis"
		csn, err = clipRedisAddress(cacheSourceAddress)
	case strings.HasPrefix(cacheSourceAddress, "rediss://"):
		csd = "rediss"
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
	case "redis", "rediss":
		var in any
		if drvDialect == "redis" {
			// redis single
			in, err = redis.ParseURL(drvSource)
		} else {
			// redis cluster
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

func Wait(ctx context.Context, drv Driver) error {
	return wait.PollImmediateUntilWithContext(ctx, 2*time.Second,
		func(ctx context.Context) (bool, error) {
			var err = drv.PingContext(ctx)
			if err != nil {
				log.Warnf("waiting for cache to be ready: %v", err)
			}
			return err == nil, ctx.Err()
		},
	)
}
