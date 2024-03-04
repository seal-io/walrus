package memcache

import (
	"context"
	"sync"
	"time"
)

var (
	once sync.Once
	gc   = MustNew(context.Background())
)

// Configure configures the global cache with given entry max life once,
// multiple calls will be ignored.
func Configure(entryMaxLife time.Duration) {
	once.Do(func() {
		_ = gc.Close()
		gc = MustNewWithConfig(context.Background(), Config{
			EntryMaxLife: entryMaxLife,
		})
	})
}

// Set sets the entry associated with the key to the global cache.
func Set(ctx context.Context, key string, entry []byte) error {
	return gc.Set(ctx, key, entry)
}

// Delete removes the keys from the global cache.
func Delete(ctx context.Context, keys ...string) error {
	return gc.Delete(ctx, keys...)
}

// Get returns the entry associated with the key from the global cache.
func Get(ctx context.Context, key string) ([]byte, error) {
	return gc.Get(ctx, key)
}

// List returns the entries associated with the keys from the global cache.
func List(ctx context.Context, keys ...string) ([][]byte, error) {
	return gc.List(ctx, keys...)
}

// Iterate iterates over the entries in the global cache.
// breaks with none nil error,
// do not do time-expensive callback during iteration.
func Iterate(ctx context.Context, m EntryKeyMatcher, a EntryAccessor) error {
	return gc.Iterate(ctx, m, a)
}
