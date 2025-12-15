package cache

import (
	"runtime"
	"sync"
	"time"
)

type DataVal[T any] struct {
	value      *T
	expires_at time.Time
}

type Cache[K comparable, T any] struct {
	data map[K]*DataVal[T]
	mu   sync.Mutex
}

func NewCache[K comparable, T any]() *Cache[K, T] {
	return &Cache[K, T]{
		data: make(map[K]*DataVal[T]),
		mu:   sync.Mutex{},
	}
}

func (c *Cache[K, T]) Get(key K) (*T, bool) {
	data, ok := c.data[key]

	if !ok {
		return nil, false
	}

	if isExpired := data.expires_at.Before(time.Now()); isExpired {
		c.Remove(key)
		return nil, false
	}

	return data.value, true
}

func (c *Cache[K, T]) Set(key K, data T, duration time.Duration) {
	c.data[key] = &DataVal[T]{
		value:      &data,
		expires_at: time.Now().Add(duration * time.Minute),
	}
}

func (c *Cache[K, T]) Remove(key K) {
	delete(c.data, key)
}

func (c *Cache[K, T]) Clear() {
	c.data = make(map[K]*DataVal[T])
	runtime.GC()
}
