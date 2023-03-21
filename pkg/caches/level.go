package caches

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/utils/cache"
	"github.com/seal-io/seal/utils/vars"
)

var Stack = vars.SetOnce[[]cache.Cache]{}

// Set saves entry with the given key,
// it returns an ErrEntryTooBigger when entry is too bigger.
func Set(ctx context.Context, key string, entry []byte) (err error) {
	var cs = Stack.Get()
	for i := range cs {
		if cs[i] == nil {
			continue
		}
		err = cs[i].Set(ctx, key, entry)
		if err != nil && !errors.Is(err, cache.ErrEntryTooBigger) {
			return
			// e.g. cannot store in l1,
			// but can store in l2.
		}
	}
	return
}

// Delete removes the keys.
func Delete(ctx context.Context, keys ...string) (err error) {
	var cs = Stack.Get()
	for i := range cs {
		if cs[i] == nil {
			continue
		}
		err = cs[i].Delete(ctx, keys...)
		if err != nil {
			return
		}
	}
	return
}

// Get reads entry for the key,
// it returns an ErrEntryNotFound when no entry exists for the given key.
func Get(ctx context.Context, key string) (entry []byte, err error) {
	var cs = Stack.Get()
	for i := range cs {
		if cs[i] == nil {
			continue
		}
		if entry != nil {
			break
		}
		// get
		entry, err = cs[i].Get(ctx, key)
		if err != nil && !errors.Is(err, cache.ErrEntryNotFound) {
			return
			// e.g. cannot find from l1,
			// but can find from l2.
		}
	}
	return
}

// List reads entries for the key list.
func List(ctx context.Context, keys ...string) (entries [][]byte, err error) {
	entries = make([][]byte, len(keys))
	var index = map[string]int{}
	for i := range keys {
		index[keys[i]] = i
	}
	var cs = Stack.Get()
	for i := range cs {
		if cs[i] == nil {
			continue
		}
		if len(keys) == 0 {
			break
		}
		// list
		var r [][]byte
		r, err = cs[i].List(ctx, keys...)
		if err != nil {
			return
		}
		// merge
		var n []string
		for j := range r {
			if r[j] == nil {
				// e.g. cannot find from l1,
				// try to find from l2.
				n = append(n, keys[j])
				continue
			}
			entries[index[keys[j]]] = r[j]
		}
		keys = n
	}
	return
}

// Iterate iterates all entries of the whole cache,
// breaks with false returning or none nil error,
// do not do time-expensive callback during iteration.
func Iterate(ctx context.Context, m cache.EntryKeyMatcher, a cache.EntryAccessor) (err error) {
	if a == nil {
		return
	}
	var av = sets.New[string]()
	var ax = func(ctx context.Context, entry cache.Entry) (bool, error) {
		var k = entry.Key()
		if av.Has(k) {
			// e.g. visit from l1,
			// do not visit again.
			return true, nil
		}
		av.Insert(k)
		return a(ctx, entry)
	}
	var cs = Stack.Get()
	for i := range cs {
		if cs[i] == nil {
			continue
		}
		err = cs[i].Iterate(ctx, m, ax)
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
	var cs = Stack.Get()
	for i := range cs {
		if cs[i] == nil {
			continue
		}
		t, ok := cs[i].(cache.Underlay[T])
		if ok {
			return t.Underlay(), nil
		}
	}
	panic(errors.New("not found cache underlay"))
}
