package wtfcache

import (
	"fmt"
	"strings"
	"sync"
)

var (
	_ cache[fmt.Stringer, string] = new(LRU[fmt.Stringer, string])
	_ cache[fmt.Stringer, string] = new(LRUWithLock[fmt.Stringer, string])

	_ fmt.Stringer = new(LRU[uint, string])
	_ fmt.Stringer = new(LRUWithLock[uint, string])
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

func (c *LRUWithLock[K, V]) String() string {
	return c.lru.String()
}

func (f *Fab[K, V]) Make(limit int) *LRU[K, V] {
	ll := linkedList[V]{root: new(element[V])}
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
	e, ok := c.dict[k]
	if ok {
		c.ll.Move(e)
		e.v = v
		return
	}

	if len(c.dict) > c.limit {
		e = c.ll.Last()
		delete(c.dict, e.k)
		e.k, e.v = k, v
		c.ll.Move(e)
		c.dict[k] = e
		return
	}

	e = c.fab.Get()
	e.k, e.v = k, v
	c.ll.Push(e)
	c.dict[k] = e
}

func (c *LRU[K, V]) Get(k K) (zero V, _ bool) {
	e, ok := c.dict[k]
	if !ok {
		return zero, false
	}

	c.ll.Move(e)

	return e.v, true
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

func (c *LRU[K, V]) String() string {
	var (
		count int
		res   = make([]string, 0, len(c.dict))
	)
	for e := c.ll.Last(); e.next != c.ll.root.next; e = e.prev {
		count++
		if _, ok := c.dict[e.k]; ok {
			res = append(res, fmt.Sprintf("%v:%v", e.k, e.v))
		} else {
			res = append(res, fmt.Sprintf("(missed)%v:%v", e.k, e.v))
		}
	}

	ret := strings.Join(res, ", ")
	if count != len(c.dict) {
		ret = "! " + ret
	}

	return ret
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
