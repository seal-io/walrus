package kms

import (
	"context"
	stdpath "path"
)

type (
	KeyValue struct {
		Path      string `json:"path"`
		Key       string `json:"key"`
		ValueHash string `json:"valueHash"`
		ValueSize string `json:"valueSize"`
	}

	Driver interface {
		// Get retrieves a data for the given key(i.e. /dir/key) from the Driver,
		// returns an error if no such key.
		Get(ctx context.Context, key string) ([]byte, error)

		// Put stores the data in the Driver under the given key(i.e. /dir/key),
		// returns error if the operation was not successful,
		// otherwise, Get, must results in the original data.
		Put(ctx context.Context, key string, value []byte) error

		// Delete removes a data from the Driver under the given key(i.e. /dir/key),
		// returns nil if no such key.
		Delete(ctx context.Context, key string) error

		// List lists all data under the path(i.e. /dir).
		List(ctx context.Context, path string) ([]KeyValue, error)
	}
)

func normalize(p string) string {
	p = stdpath.Clean(p)
	p = stdpath.Join("/", p)

	return p
}

func point[T ~string](s T) *T {
	return &s
}
