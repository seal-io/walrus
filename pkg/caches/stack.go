package caches

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/utils/cache"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/vars"
)

var Stack = vars.SetOnce[[]cache.Cache]{}

// Set saves entry with the given key.
func Set(ctx context.Context, key string, entry []byte) (err error) {
	logger := log.WithName("cache")
	layers := Stack.Get()

	for i := range layers {
		if layers[i] == nil {
			continue
		}

		err = layers[i].Set(ctx, key, entry)
		if err != nil && !errors.Is(err, cache.ErrEntryTooBig) {
			logger.Warnf("cannot set entry %q into layer \"%s(%d)\": %v",
				key, layers[i].Name(), i, err)
		}
	}

	return nil
}

// Delete removes the keys.
func Delete(ctx context.Context, keys ...string) (err error) {
	logger := log.WithName("cache")
	layers := Stack.Get()

	for i := range layers {
		if layers[i] == nil {
			continue
		}

		err = layers[i].Delete(ctx, keys...)
		if err != nil && !errors.Is(err, cache.ErrEntryNotFound) {
			logger.Warnf("cannot delete from layer \"%s(%d)\": %v",
				layers[i].Name(), i, err)
		}
	}

	return nil
}

// Get reads entry for the given key,
// it returns an ErrEntryNotFound when no entry exists for the given key.
func Get(ctx context.Context, key string) (entry []byte, err error) {
	logger := log.WithName("cache")
	layers := Stack.Get()
	err = cache.ErrEntryNotFound

	for i := range layers {
		if layers[i] == nil {
			continue
		}

		if entry != nil {
			// Back source.
			for j := i - 1; j >= 0; j-- {
				if layers[j] == nil {
					continue
				}

				err = layers[j].Set(ctx, key, entry)
				if err != nil && !errors.Is(err, cache.ErrEntryTooBig) {
					logger.Warnf("cannot back-set entry %q into layer \"%s(%d)\": %v",
						key, layers[j].Name(), j, err)
				}
			}

			break
		}

		entry, err = layers[i].Get(ctx, key)
		if err != nil {
			if !errors.Is(err, cache.ErrEntryNotFound) {
				return
			}
		} else {
			logger.V(5).Infof("get entry %q from layer \"%s(%d)\"",
				key, layers[i].Name(), i)
		}
		// NB(thxCode): cannot find from this layer,
		// try to find from next layer.
	}

	return
}

// List reads entries for the given key list,
// it returns an entry list with nil value when no entry exists for the given key.
func List(ctx context.Context, keys ...string) (entries [][]byte, err error) {
	logger := log.WithName("cache")
	layers := Stack.Get()

	keysIndex := map[string]int{}
	for i := range keys {
		keysIndex[keys[i]] = i
	}

	entries = make([][]byte, len(keys))

	for i := range layers {
		if layers[i] == nil {
			continue
		}

		if len(keys) == 0 {
			break
		}

		var r [][]byte

		r, err = layers[i].List(ctx, keys...)
		if err != nil {
			return
		}

		var residueKeys, foundKeys []string

		for j := range r {
			if r[j] == nil {
				// NB(thxCode): cannot find from this layer,
				// try to find from next layer.
				residueKeys = append(residueKeys, keys[j])
				continue
			}

			entries[keysIndex[keys[j]]] = r[j]
			foundKeys = append(foundKeys, keys[j])
		}

		keys = residueKeys

		logger.V(5).Infof("list entries %q from layer \"%s(%d)\"",
			foundKeys, layers[i].Name(), i)
	}

	return
}

// Iterate iterates all entries of the whole cache,
// breaks with none nil error,
// do not do time-expensive callback during iteration.
func Iterate(ctx context.Context, m cache.EntryKeyMatcher, a cache.EntryAccessor) (err error) {
	if a == nil {
		return
	}

	layers := Stack.Get()

	av := sets.New[string]()
	ax := func(ctx context.Context, entry cache.Entry) (bool, error) {
		k := entry.Key()
		if av.Has(k) {
			// NB(thxCode): already visited.
			return true, nil
		}

		av.Insert(k)

		return a(ctx, entry)
	}

	for i := range layers {
		if layers[i] == nil {
			continue
		}

		err = layers[i].Iterate(ctx, m, ax)
		if err != nil {
			return
		}
	}

	return
}

// Underlay returns the underlay client with the given generic, e.g.
// process cache implements with *bigcache.BigCache,
// remote cache implements with redis.UniversalClient.
// DO NOT close the underlay client.
func Underlay[T any]() (T, error) {
	layers := Stack.Get()

	for i := range layers {
		if layers[i] == nil {
			continue
		}

		t, ok := layers[i].(cache.Underlay[T])
		if ok {
			return t.Underlay(), nil
		}
	}

	panic(errors.New("not found cache underlay"))
}
