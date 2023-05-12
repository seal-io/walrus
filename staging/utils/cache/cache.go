package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/dustin/go-humanize"

	"github.com/seal-io/seal/utils/log"
)

func MustNew(ctx context.Context) *bigcache.BigCache {
	n, err := New(ctx)
	if err != nil {
		panic(fmt.Errorf("error creating cache: %w", err))
	}

	return n
}

func New(ctx context.Context) (*bigcache.BigCache, error) {
	// Each shard initializes with `(MaxEntriesInWindows / Shards) * MaxEntrySize` = 300 * 512 = 150kb
	// each shard limits in `(HardMaxCacheSize * 1024 * 1024) / Shards` = 64 * 1024 * 1024 / 64 = 1mb
	// initializes with 64 * 150kb = 9mb, limits with 64 * 1mb = 64mb.
	cfg := bigcache.Config{
		LifeWindow:         15 * time.Minute,
		CleanWindow:        3 * time.Minute,
		Shards:             64,
		MaxEntriesInWindow: 64 * 5 * 60,
		MaxEntrySize:       512,
		HardMaxCacheSize:   64,
		StatsEnabled:       false,
		Verbose:            false,
	}

	return NewWithConfig(ctx, cfg)
}

func NewWithConfig(ctx context.Context, cfg bigcache.Config) (*bigcache.BigCache, error) {
	if cfg.Logger == nil {
		cfg.Logger = log.WithName("cache")
	}

	if cfg.OnRemoveWithReason == nil {
		cfg.OnRemoveWithReason = func(key string, entry []byte, reason bigcache.RemoveReason) {
			desc := "unknown"

			switch reason {
			case bigcache.Deleted:
				desc = "deleted"
			case bigcache.Expired:
				desc = "expired"
			case bigcache.NoSpace:
				desc = "nospace"
			}

			size := humanize.IBytes(uint64(len(entry)))
			cfg.Logger.Printf("%s: %10s | %s", desc, size, key)
		}
	}

	return bigcache.New(ctx, cfg)
}
