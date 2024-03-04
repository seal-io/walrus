package memcache

import (
	"context"
	"errors"
	"io"
	"strings"
)

var (
	ErrEntryNotFound = errors.New("entry is not found")
	ErrEntryTooBig   = errors.New("entry is too big")
)

// Entry holds the action of cache entry.
type Entry interface {
	// Key returns entry key.
	Key() string

	// Value returns entry value.
	Value() ([]byte, error)
}

// EntryKeyMatcher holds the predication of the iteration.
type EntryKeyMatcher interface {
	// Match returns true if the given key matched.
	Match(key string) bool
}

// HasPrefix implements the EntryKeyMatcher stereotype,
// which means the key has the given prefix.
type HasPrefix string

func (s HasPrefix) Match(key string) bool {
	return strings.HasPrefix(key, string(s))
}

// HasSuffix implements the EntryKeyMatcher stereotype,
// which means the key has the given suffix.
type HasSuffix string

func (s HasSuffix) Match(key string) bool {
	return strings.HasSuffix(key, string(s))
}

// HasPartial implements the EntryKeyMatcher stereotype,
// which means the key has the given partial.
type HasPartial string

func (s HasPartial) Match(key string) bool {
	return strings.Contains(key, string(s))
}

// EntryAccessor accesses the entry during iterating,
// it can break the iteration with a false returning or none nil error.
type EntryAccessor func(ctx context.Context, entry Entry) (next bool, err error)

// Cache holds the action of caching.
type Cache interface {
	io.Closer

	Name() string

	// Set saves entry with the given key,
	// it returns an ErrEntryTooBig when entry is too big.
	Set(ctx context.Context, key string, entry []byte) error

	// Delete removes the keys.
	Delete(ctx context.Context, keys ...string) error

	// Get reads entry for the key,
	// it returns an ErrEntryNotFound when no entry exists for the given key.
	Get(ctx context.Context, key string) ([]byte, error)

	// List reads entries for the key list.
	List(ctx context.Context, keys ...string) (entries [][]byte, err error)

	// Iterate iterates all entries of the whole cache,
	// breaks with none nil error,
	// do not do time-expensive callback during iteration.
	Iterate(ctx context.Context, m EntryKeyMatcher, a EntryAccessor) error
}

// Underlay gets the underlay client.
type Underlay[T any] interface {
	Underlay() T
}
