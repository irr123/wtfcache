package wtfcache

import (
	"container/list"
	"sync"
)

type (
	element[V any] struct {
		k any
		v V
	}
	LRU[K, V any] struct {
		limit int
		mu    sync.Mutex
		ll    *list.List
		dict  map[any]*list.Element
	}
)

func New[K, V any](limit int) *LRU[K, V] {
	const size = 3

	return &LRU[K, V]{
		limit: limit - 1,
		ll:    list.New().Init(),
		dict:  make(map[any]*list.Element, size),
	}
}

func (c *LRU[K, V]) Set(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.dict[k]
	if ok {
		c.ll.MoveToFront(e)
		return
	}

	if c.ll.Len() > c.limit {
		last := c.ll.Back()
		c.ll.Remove(last)
		delete(c.dict, last.Value.(element[V]).k)
	}

	c.dict[k] = c.ll.PushFront(element[V]{k, v})
}

func (c *LRU[K, V]) Get(k K) (nilValue V, _ bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.dict[k]
	if !ok {
		return nilValue, false
	}

	c.ll.MoveToFront(e)

	return e.Value.(element[V]).v, ok
}

func (c *LRU[K, V]) Del(k K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.dict[k]
	if !ok {
		return
	}

	c.ll.Remove(e)
	delete(c.dict, k)
}
