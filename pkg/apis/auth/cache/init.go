package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"

	"github.com/seal-io/seal/utils/cache"
)

var cacher *bigcache.BigCache

func init() {
	// Narrow the expiry form default cache creator.
	cfg := bigcache.Config{
		LifeWindow:         5 * time.Minute,
		CleanWindow:        2 * time.Minute,
		Shards:             64,
		MaxEntriesInWindow: 64 * 5 * 60,
		MaxEntrySize:       512,
		HardMaxCacheSize:   64,
		StatsEnabled:       false,
		Verbose:            false,
	}
	var err error
	cacher, err = cache.NewWithConfig(context.Background(), cfg)
	if err != nil {
		panic(fmt.Errorf("error creating cache: %w", err))
	}
}
