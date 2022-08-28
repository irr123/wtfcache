package wtfcache

import "sync"

var (
	_ cache[uint, string] = new(LRU[uint, string])
	_ cache[uint, string] = new(LRUWithLock[uint, string])
)

type (
	cache[K, V any] interface {
		Set(K, V)
		Get(K) (V, bool)
		Del(K)
	}

	element[V any] struct {
		next, prev *element[V]

		k any
		v V
	}

	linkedList[V any] struct{ root *element[V] }
	Fab[K, V any]     struct{ pool sync.Pool }

	LRU[K, V any] struct {
		fab *Fab[K, V]

		limit int
		ll    linkedList[V]
		dict  map[any]*element[V]
	}
	LRUWithLock[K, V any] struct {
		mu  sync.Mutex
		lru *LRU[K, V]
	}
)

func New[K, V any]() *Fab[K, V] {
	return &Fab[K, V]{sync.Pool{
		New: func() any { return new(element[V]) },
	}}
}

func (f *Fab[K, V]) MakeWithLock(limit int) *LRUWithLock[K, V] {
	return &LRUWithLock[K, V]{lru: f.Make(limit)}
}

func (c *LRUWithLock[K, V]) Set(k K, v V) {
	c.mu.Lock()
	c.lru.Set(k, v)
	c.mu.Unlock()
}

func (c *LRUWithLock[K, V]) Get(k K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lru.Get(k)
}

func (c *LRUWithLock[K, V]) Del(k K) {
	c.mu.Lock()
	c.lru.Del(k)
	c.mu.Unlock()
}

func (f *Fab[K, V]) Make(limit int) *LRU[K, V] {
	var ll = linkedList[V]{root: new(element[V])}
	ll.root.next = ll.root
	ll.root.prev = ll.root

	return &LRU[K, V]{
		fab:   f,
		limit: limit - 1,
		ll:    ll,
		dict:  make(map[any]*element[V], limit),
	}
}

func (c *LRU[K, V]) Set(k K, v V) {
	if e, ok := c.dict[k]; ok {
		c.ll.Move(e)
		return
	}

	if len(c.dict) > c.limit {
		last := c.ll.Last()
		c.ll.Del(last)
		delete(c.dict, last.k)
		c.fab.Put(last)
	}

	var elem = c.fab.Get()
	elem.k, elem.v = k, v
	c.ll.Push(elem)
	c.dict[k] = elem
}

func (c *LRU[K, V]) Get(k K) (nilVal V, _ bool) {
	e, ok := c.dict[k]
	if !ok {
		return nilVal, false
	}

	c.ll.Move(e)

	return e.v, ok
}

func (c *LRU[K, V]) Del(k K) {
	e, ok := c.dict[k]
	if !ok {
		return
	}

	c.ll.Del(e)
	delete(c.dict, k)
	c.fab.Put(e)
}

func (ll linkedList[V]) Push(e *element[V]) {
	e.prev = ll.root
	e.next = ll.root.next
	e.prev.next = e
	e.next.prev = e
}

func (ll linkedList[V]) Move(e *element[V]) {
	if ll.root.next == e {
		return
	}

	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = ll.root
	e.next = ll.root.next
	e.prev.next = e
	e.next.prev = e
}

func (ll linkedList[V]) Last() *element[V] {
	return ll.root.prev
}

func (ll linkedList[V]) Del(e *element[V]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
}

func (f *Fab[K, V]) Put(x any) {
	f.pool.Put(x)
}

func (f *Fab[K, V]) Get() *element[V] {
	return f.pool.Get().(*element[V])
}
