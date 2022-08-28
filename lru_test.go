package wtfcache

import (
	"fmt"
	"sort"
	"testing"
)

func Test(t *testing.T) {
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

func Example() {
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

	fmt.Printf("dict: %d\n", len(c.lru.dict))
	values := make([]string, 0, len(c.lru.dict))
	for k, v := range c.lru.dict {
		values = append(values, fmt.Sprintf(`"%v: %v"`, k, v.v))
	}

	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })

	fmt.Printf("%v\nlist: %d\n", values, len(c.lru.dict))
	elem := c.lru.ll.Last()
	for elem != nil {
		fmt.Printf(`"%v: %v"`+"\n", elem.k, elem.v)
		c.lru.ll.Del(elem)
		elem = c.lru.ll.Last()
	}

	// Output:
	// dict: 2
	// ["{5}: {5}" "{six}: {6}"]
	// list: 2
	// "{5}: {5}"
	// "{six}: {6}"
	// "<nil>: {0}"
}
