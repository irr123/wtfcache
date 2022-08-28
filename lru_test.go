package wtfcache

import (
	"fmt"
	"testing"
)

func TestCommon(t *testing.T) {
	const (
		limit      = 9
		iterations = 101
	)

	var c = New[int, int]().MakeWithLock(limit)

	for i := 0; i < iterations; i++ {
		c.Set(i, i)
	}

	if len(c.lru.dict) != limit {
		t.Fatalf("1, %v != %v", limit, len(c.lru.dict))
	}

	for i := 0; i < iterations; i++ {
		if i < iterations-limit {
			v, ok := c.Get(i)
			if ok {
				t.Fatalf("3, %v = %v, (%v)", i, v, len(c.lru.dict))
			}

			continue
		}

		v, ok := c.Get(i)
		if !ok {
			t.Fatalf("4, %v = %v, (%v)", i, v, len(c.lru.dict))
		}

		if v != i {
			t.Fatalf("5, %v = %v, (%v)", i, v, len(c.lru.dict))
		}

		c.Del(i)
	}

	if len(c.lru.dict) != 0 {
		t.Fatalf("6, %v", len(c.lru.dict))
	}
}

func TestResetExistingKey(t *testing.T) {
	var c = New[int, int]().MakeWithLock(10)

	c.Set(1, 2)
	c.Set(1, 3)

	v, ok := c.Get(1)
	if !ok {
		t.Fatalf("1, %v", ok)
	}

	if v != 3 {
		t.Fatalf("2, %v != 3", v)
	}
}

func TestDelFromEmpty(t *testing.T) {
	var c = New[int, int]().MakeWithLock(10)

	c.Del(1)
}

func Example() {
	var c = New[int, struct{}]().MakeWithLock(10)
	fmt.Println(c.String())

	// Output:
	//
}

func ExampleLRU_String() {
	var c = New[int, struct{}]().Make(10)

	c.Set(0, struct{}{})
	c.Set(1, struct{}{})
	c.Set(2, struct{}{})
	c.Set(3, struct{}{})
	c.Set(4, struct{}{})
	c.Set(5, struct{}{})
	c.Set(6, struct{}{})
	c.Set(7, struct{}{})
	c.Set(8, struct{}{})
	c.Set(9, struct{}{})
	c.Set(0, struct{}{})

	fmt.Printf("len: %d\n", len(c.dict))
	fmt.Println(c.String())

	// Output:
	// len: 10
	// 1:{}, 2:{}, 3:{}, 4:{}, 5:{}, 6:{}, 7:{}, 8:{}, 9:{}, 0:{}
}

func ExampleLRUWithLock_String() {
	type (
		key   struct{ K string }
		value struct{ V int }
	)

	var c = New[key, value]().MakeWithLock(3)

	c.Set(key{"one"}, value{1})
	c.Set(key{"two"}, value{2})
	c.Set(key{"3"}, value{3})
	c.Set(key{"3"}, value{3})
	c.Set(key{"3"}, value{3})
	c.Set(key{"4"}, value{4})
	c.Set(key{"5"}, value{5})
	c.Set(key{"six"}, value{6})
	c.Set(key{"seven"}, value{7})
	c.Del(key{"seven"})

	fmt.Printf("len: %d\n", len(c.lru.dict))
	fmt.Println(c.String())

	// Output:
	// len: 2
	// {5}:{5}, {six}:{6}
}

func ExampleLRU_Del() {
	var c = New[int, struct{}]().Make(3)

	c.Set(0, struct{}{})
	c.Set(1, struct{}{})
	c.Set(2, struct{}{})
	delete(c.dict, 1)

	fmt.Printf("len: %d\n", len(c.dict))
	fmt.Println(c.String())

	// Output:
	// len: 2
	// ! 0:{}, (missed)1:{}, 2:{}
}

func BenchmarkWTF(b *testing.B) {
	cache := New[int, int]().MakeWithLock(b.N)

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Set(i, i)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, ok := cache.Get(i)
			if ok {
				_ = value
			}
		}
	})
}
