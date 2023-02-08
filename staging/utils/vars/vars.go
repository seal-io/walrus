package vars

import "sync"

type SetOnce[T any] struct {
	o sync.Once
	v T
}

func (i *SetOnce[T]) Set(v T) {
	i.o.Do(func() {
		i.v = v
	})
}

func (i *SetOnce[T]) Get() T {
	return i.v
}

type SetMany[T any] struct {
	m sync.RWMutex
	v T
}

func (i *SetMany[T]) Set(v T) {
	i.m.Lock()
	defer i.m.Unlock()
	i.v = v
}

func (i *SetMany[T]) Get() T {
	i.m.RLock()
	defer i.m.RUnlock()
	return i.v
}
