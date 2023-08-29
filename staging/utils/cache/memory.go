package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/dustin/go-humanize"

	"github.com/seal-io/walrus/utils/log"
)

// MemoryConfig holds the configuration of the in-memory cache,
// entry indexes by key and stores in one bucket,
// the total cache size is BucketCapacity * Buckets.
type MemoryConfig struct {
	// Namespace indicates the operating workspace.
	Namespace string
	// EntryMaxAge indicates the maximum lifetime of each entry,
	// default is 15 mins.
	EntryMaxAge time.Duration
	// LazyEntryEviction indicates to evict an expired entry at next peeking,
	// by default, a background looping tries to evict expired entries per 3 mins.
	LazyEntryEviction bool
	// Buckets indicates the bucket number of cache,
	// value must be a power of two,
	// default is 64.
	Buckets int
	// BucketCapacity indicates the maximum MB of each bucket,
	// default is 1 MB.
	BucketCapacity int
	// LazyBucketScale indicates to scale when the current bucket is not enough to put a new entry,
	// by default, create the bucket with the given capacity to avoid any array copying.
	// It's worth noticing that the bucket capacity can not exceed even configured LazyBucketScale to true.
	LazyBucketScale bool
}

func (c *MemoryConfig) Default() {
	c.Namespace = strings.TrimSpace(c.Namespace)
	if c.EntryMaxAge == 0 {
		c.EntryMaxAge = 15 * time.Minute
	}

	if c.Buckets == 0 {
		c.Buckets = 64
	}

	if c.BucketCapacity == 0 {
		c.BucketCapacity = 1
	}
}

func (c *MemoryConfig) Validate() error {
	if c.EntryMaxAge < 0 {
		return errors.New("invalid entry max age: negative")
	}

	if c.Buckets < 0 {
		return errors.New("invalid buckets: negative")
	}

	if c.BucketCapacity < 0 {
		return errors.New("invalid bucket capacity: negative")
	}

	return nil
}

// NewMemory returns an in-memory Cache implementation.
func NewMemory(ctx context.Context) (Cache, error) {
	return NewMemoryWithConfig(ctx, MemoryConfig{})
}

// MustNewMemory likes NewMemory, but panic if error found.
func MustNewMemory(ctx context.Context) Cache {
	n, err := NewMemory(ctx)
	if err != nil {
		panic(fmt.Errorf("error creating in-memory cache: %w", err))
	}

	return n
}

// NewMemoryWithConfig returns an in-memory Cache implementation with given configuration.
func NewMemoryWithConfig(ctx context.Context, cfg MemoryConfig) (Cache, error) {
	// Default, validate.
	cfg.Default()

	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	// Generate bigcache configuration with MemoryConfig.
	//
	// For example:
	//
	// bigcache.Config{
	//		LifeWindow:         15 * time.Minute,
	//		CleanWindow:        3 * time.Minute,
	//		Shards:             64,
	//		MaxEntriesInWindow: 64 * 300,  // works with MaxEntrySize to determinate the cache initialization.
	//		MaxEntrySize:       512,
	//		HardMaxCacheSize:   64,
	//		StatsEnabled:       false,
	//		Verbose:            false,
	//	}
	//
	// Each shard initializes with `(MaxEntriesInWindows / Shards) * MaxEntrySize` = 300 * 512 = 150kb.
	// Each shard limits in `(HardMaxCacheSize * 1024 * 1024) / Shards` = 64 * 1024 * 1024 / 64 = 1mb.
	// Initializes with 64 * 150kb = 9mb, limits with 64 * 1mb = 64mb.
	//
	capacity := cfg.BucketCapacity * cfg.Buckets
	logger := log.WithName("cache").WithName("memory")

	underlayCfg := bigcache.Config{
		Shards:             cfg.Buckets,
		LifeWindow:         cfg.EntryMaxAge,
		CleanWindow:        0,
		MaxEntriesInWindow: cfg.Buckets << 4,
		MaxEntrySize:       cfg.BucketCapacity << (20 - 4),
		HardMaxCacheSize:   capacity,
		StatsEnabled:       false,
		Verbose:            false,
		Logger:             logger,
		OnRemoveWithReason: func(key string, entry []byte, reason bigcache.RemoveReason) {
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
			if len(key) > 10 {
				logger.Debugf("%s: %10s | %10s...", desc, size, key[:10])
			} else {
				logger.Debugf("%s: %10s | %10s", desc, size, key)
			}
		},
	}
	if !cfg.LazyEntryEviction {
		// Set up a background looping to clean.
		underlayCfg.CleanWindow = 3 * time.Minute
	}

	if cfg.LazyBucketScale {
		// Initialize the cache queue in 1/4 capacity.
		underlayCfg.MaxEntrySize >>= 2
	}

	// Init.
	underlay, err := bigcache.New(ctx, underlayCfg)
	if err != nil {
		return nil, err
	}
	mc := memoryCache{
		underlay: underlay,
	}

	if cfg.Namespace != "" {
		ns := cfg.Namespace + "#"
		mc.namespace = &ns
	}

	return mc, nil
}

// MustNewMemoryWithConfig likes NewMemoryWithConfig, but panic if error found.
func MustNewMemoryWithConfig(ctx context.Context, cfg MemoryConfig) Cache {
	n, err := NewMemoryWithConfig(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("error creating in-memory cache: %w", err))
	}

	return n
}

// memoryCache adapts Cache interface to implement an in-memory cache with bigcache.BigCache.
type memoryCache struct {
	namespace *string
	underlay  *bigcache.BigCache
}

func (c memoryCache) wrapKey(s *string) *string {
	if c.namespace == nil {
		return s
	}
	r := *c.namespace + *s

	return &r
}

func (c memoryCache) Close() error {
	return c.underlay.Close()
}

func (c memoryCache) Name() string {
	return "memory"
}

func (c memoryCache) Set(ctx context.Context, key string, entry []byte) (err error) {
	wk := c.wrapKey(&key)

	err = c.underlay.Set(*wk, entry)
	if err != nil && err.Error() == "entry is bigger than max shard size" {
		err = ErrEntryTooBig
	}

	return
}

func (c memoryCache) Delete(ctx context.Context, keys ...string) (err error) {
	for i := range keys {
		wk := c.wrapKey(&keys[i])

		err = c.underlay.Delete(*wk)
		if err != nil {
			if !errors.Is(err, bigcache.ErrEntryNotFound) {
				return
			}
			err = nil
		}
	}

	return
}

func (c memoryCache) Get(ctx context.Context, key string) (entry []byte, err error) {
	wk := c.wrapKey(&key)

	entry, err = c.underlay.Get(*wk)
	if err != nil && errors.Is(err, bigcache.ErrEntryNotFound) {
		err = ErrEntryNotFound
	}

	return
}

func (c memoryCache) List(ctx context.Context, keys ...string) (entries [][]byte, err error) {
	entries = make([][]byte, len(keys))
	for i := range keys {
		entries[i], err = c.Get(ctx, keys[i])
		if err != nil {
			return
		}
	}

	return
}

func (c memoryCache) Iterate(ctx context.Context, m EntryKeyMatcher, a EntryAccessor) error {
	if a == nil {
		return nil
	}

	it := c.underlay.Iterator()
	for it.SetNext() {
		e, err := it.Value()
		if err != nil {
			return err
		}

		k := e.Key()
		if c.namespace != nil && !strings.HasPrefix(k, *c.namespace) {
			continue
		}

		if m != nil && !m.Match(k) {
			continue
		}

		n, err := a(ctx, bigcacheEntry{i: e})
		if err != nil {
			return err
		}

		if !n {
			break
		}
	}

	return nil
}

func (c memoryCache) Underlay() *bigcache.BigCache {
	return c.underlay
}

type bigcacheEntry struct {
	i bigcache.EntryInfo
}

func (e bigcacheEntry) Key() string {
	return e.i.Key()
}

func (e bigcacheEntry) Value() ([]byte, error) {
	return e.i.Value(), nil
}
