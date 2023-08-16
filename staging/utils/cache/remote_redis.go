package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/seal-io/walrus/utils/strs"
)

// RemoteRedisConfig holds the configuration of the remote cache.
type RemoteRedisConfig struct {
	// Namespace indicates the operating workspace.
	Namespace string
	// EntryMaxAge indicates the maximum lifetime of each entry,
	// default is 15 mins.
	EntryMaxAge time.Duration
	// Client indicates the underlay client.
	Client redis.UniversalClient
}

func (c *RemoteRedisConfig) Default() {
	c.Namespace = strings.TrimSpace(c.Namespace)
	if c.EntryMaxAge == 0 {
		c.EntryMaxAge = 15 * time.Minute
	}
}

func (c *RemoteRedisConfig) Validate() error {
	if c.EntryMaxAge < 0 {
		return errors.New("invalid entry max age: negative")
	}

	if c.Client == nil {
		return errors.New("invalid client: nil")
	}

	return nil
}

// NewRemoteRedis returns a remote Cache implementation.
func NewRemoteRedis(ctx context.Context) (Cache, error) {
	return NewRemoteRedisWithConfig(ctx, RemoteRedisConfig{})
}

// MustNewRemoteRedis likes NewRemoteRedis, but panic if error found.
func MustNewRemoteRedis(ctx context.Context) Cache {
	n, err := NewRemoteRedis(ctx)
	if err != nil {
		panic(fmt.Errorf("error creating remote cache: %w", err))
	}

	return n
}

// NewRemoteRedisWithConfig returns a remote Cache implementation with given configuration.
func NewRemoteRedisWithConfig(ctx context.Context, cfg RemoteRedisConfig) (Cache, error) {
	// Default, validate.
	cfg.Default()

	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	// Generate redis configuration with RemoteRedisConfig.
	underlay := cfg.Client
	rc := remoteRedisCache{
		expiration: cfg.EntryMaxAge,
		underlay:   underlay,
	}

	if cfg.Namespace != "" {
		ns := cfg.Namespace + "#"
		rc.namespace = &ns
	}

	return rc, nil
}

// MustNewRemoteRedisWithConfig likes NewRemoteRedisWithConfig, but panic if error found.
func MustNewRemoteRedisWithConfig(ctx context.Context, cfg RemoteRedisConfig) Cache {
	n, err := NewRemoteRedisWithConfig(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("error creating remote cache: %w", err))
	}

	return n
}

// remoteRedisCache adapts Cache interface to implement a remote cache with redis.Client.
type remoteRedisCache struct {
	namespace  *string
	expiration time.Duration
	underlay   redis.UniversalClient
}

func (c remoteRedisCache) wrapKey(s *string) *string {
	if c.namespace == nil {
		return s
	}
	r := *c.namespace + *s

	return &r
}

func (c remoteRedisCache) Close() error {
	return c.underlay.Close()
}

func (c remoteRedisCache) Name() string {
	return "redis"
}

func (c remoteRedisCache) Set(ctx context.Context, key string, entry []byte) (err error) {
	wk := c.wrapKey(&key)
	err = c.underlay.Set(ctx, *wk, entry, c.expiration).Err()

	return
}

func (c remoteRedisCache) Delete(ctx context.Context, keys ...string) (err error) {
	for i := range keys {
		wk := c.wrapKey(&keys[i])
		keys[i] = *wk
	}

	err = c.underlay.Del(ctx, keys...).Err()
	if err != nil && errors.Is(err, redis.Nil) {
		err = nil
	}

	return
}

func (c remoteRedisCache) Get(ctx context.Context, key string) (entry []byte, err error) {
	wk := c.wrapKey(&key)

	entry, err = c.underlay.Get(ctx, *wk).Bytes()
	if err != nil && errors.Is(err, redis.Nil) {
		err = ErrEntryNotFound
	}

	return
}

func (c remoteRedisCache) List(ctx context.Context, keys ...string) (entries [][]byte, err error) {
	for i := range keys {
		wk := c.wrapKey(&keys[i])
		keys[i] = *wk
	}

	values, err := c.underlay.MGet(ctx, keys...).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		err = ErrEntryNotFound
		return
	}
	entries = make([][]byte, len(keys))

	for i := range values {
		if values[i] == nil {
			continue
		}

		v, ok := values[i].(string)
		if !ok {
			continue
		}
		entries[i] = strs.ToBytes(&v)
	}

	return
}

func (c remoteRedisCache) Iterate(ctx context.Context, m EntryKeyMatcher, a EntryAccessor) error {
	if a == nil {
		return nil
	}

	var prefix string
	if c.namespace != nil {
		prefix = *c.namespace
	}

	if m != nil {
		hasPrefix, ok := m.(HasPrefix)
		if ok {
			prefix += string(hasPrefix)
			m = nil
		}
	}

	if prefix != "" {
		prefix += "*"
	}

	it := c.underlay.Scan(ctx, 0, prefix, 0).Iterator()
	for it.Next(ctx) {
		err := it.Err()
		if err != nil {
			return err
		}

		k := it.Val()
		if m != nil && !m.Match(k) {
			continue
		}

		n, err := a(ctx, redisEntry{c: ctx, u: c.underlay, k: &k})
		if err != nil {
			return err
		}

		if !n {
			break
		}
	}

	return nil
}

func (c remoteRedisCache) Underlay() redis.UniversalClient {
	return c.underlay
}

type redisEntry struct {
	c context.Context
	u redis.UniversalClient
	k *string
}

func (e redisEntry) Key() string {
	return *e.k
}

func (e redisEntry) Value() ([]byte, error) {
	return e.u.Get(e.c, *e.k).Bytes()
}
