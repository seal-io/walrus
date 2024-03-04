package memcache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/dustin/go-humanize"
	"k8s.io/klog/v2"
)

// Config holds the configuration of the in-memory cache,
// entry indexes by key and stores in one bucket,
// the total cache size is BucketCapacity * Buckets.
type Config struct {
	// Namespace indicates the operating workspace.
	Namespace string
	// EntryMaxLife indicates the maximum lifetime of each entry,
	// default is 15 mins.
	EntryMaxLife time.Duration
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

func (c *Config) Default() {
	c.Namespace = strings.TrimSpace(c.Namespace)
	if c.EntryMaxLife == 0 {
		c.EntryMaxLife = 15 * time.Minute
	}

	if c.Buckets == 0 {
		c.Buckets = 64
	}

	if c.BucketCapacity == 0 {
		c.BucketCapacity = 1
	}
}

func (c *Config) Validate() error {
	if c.EntryMaxLife < 0 {
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

// New returns an in-memory Cache implementation.
func New(ctx context.Context) (Cache, error) {
	return NewWithConfig(ctx, Config{})
}

// MustNew likes New, but panic if error found.
func MustNew(ctx context.Context) Cache {
	n, err := New(ctx)
	if err != nil {
		panic(fmt.Errorf("create in-memory cache: %w", err))
	}

	return n
}

// NewWithConfig returns an in-memory Cache implementation with given configuration.
func NewWithConfig(ctx context.Context, cfg Config) (Cache, error) {
	// Default, validate.
	cfg.Default()

	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	// Generate bigcache configuration with Config.
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
	logger := klog.Background().WithName("cache").WithName("memory").V(5)

	underlayCfg := bigcache.Config{
		Shards:             cfg.Buckets,
		LifeWindow:         cfg.EntryMaxLife,
		CleanWindow:        0,
		MaxEntriesInWindow: cfg.Buckets << 4,
		MaxEntrySize:       cfg.BucketCapacity << (20 - 4),
		HardMaxCacheSize:   capacity,
		StatsEnabled:       false,
		Verbose:            false,
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
			if logger.Enabled() {
				size := humanize.IBytes(uint64(len(entry)))
				if len(key) > 10 {
					key = key[:10] + "..."
				}
				logger.Info("removed",
					"reason", desc, "size", size, "key", key[:10])
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
	mc := bigCache{
		underlay: underlay,
	}

	if cfg.Namespace != "" {
		ns := cfg.Namespace + "#"
		mc.namespace = &ns
	}

	return mc, nil
}

// MustNewWithConfig likes NewWithConfig, but panic if error found.
func MustNewWithConfig(ctx context.Context, cfg Config) Cache {
	n, err := NewWithConfig(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("create in-memory cache: %w", err))
	}

	return n
}

// bigCache adapts Cache interface to implement an in-memory cache with bigcache.BigCache.
type bigCache struct {
	namespace *string
	underlay  *bigcache.BigCache
}

func (c bigCache) wrapKey(s *string) *string {
	if c.namespace == nil {
		return s
	}
	r := *c.namespace + *s

	return &r
}

func (c bigCache) Close() error {
	return c.underlay.Close()
}

func (c bigCache) Name() string {
	return "memory"
}

func (c bigCache) Set(ctx context.Context, key string, entry []byte) (err error) {
	wk := c.wrapKey(&key)

	err = c.underlay.Set(*wk, entry)
	if err != nil && err.Error() == "entry is bigger than max shard size" {
		err = ErrEntryTooBig
	}

	return
}

func (c bigCache) Delete(ctx context.Context, keys ...string) (err error) {
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

func (c bigCache) Get(ctx context.Context, key string) (entry []byte, err error) {
	wk := c.wrapKey(&key)

	entry, err = c.underlay.Get(*wk)
	if err != nil && errors.Is(err, bigcache.ErrEntryNotFound) {
		err = ErrEntryNotFound
	}

	return
}

func (c bigCache) List(ctx context.Context, keys ...string) (entries [][]byte, err error) {
	entries = make([][]byte, len(keys))
	for i := range keys {
		entries[i], err = c.Get(ctx, keys[i])
		if err != nil {
			return
		}
	}

	return
}

func (c bigCache) Iterate(ctx context.Context, m EntryKeyMatcher, a EntryAccessor) error {
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

func (c bigCache) Underlay() *bigcache.BigCache {
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
