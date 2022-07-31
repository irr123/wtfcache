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

	var c = New[int, int](limit)

	for i := 0; i < iterations; i++ {
		c.Set(i, i)
	}

	if c.ll.Len() != limit {
		t.Fatalf("1, %v/%v", c.ll.Len(), len(c.dict))
	}
	if len(c.dict) != limit {
		t.Fatalf("2, %v/%v", c.ll.Len(), len(c.dict))
	}

	for i := 0; i < iterations; i++ {
		if i < iterations-limit {
			v, ok := c.Get(i)
			if ok {
				t.Fatalf("3, %v = %v, (%v/%v)", i, v, c.ll.Len(), len(c.dict))
			}

			continue
		}

		v, ok := c.Get(i)
		if !ok {
			t.Fatalf("4, %v = %v, (%v/%v)", i, v, c.ll.Len(), len(c.dict))
		}

		if v != i {
			t.Fatalf("5, %v = %v, (%v/%v)", i, v, c.ll.Len(), len(c.dict))
		}

		c.Del(i)
	}

	if c.ll.Len() != 0 {
		t.Fatalf("6, %v/%v", c.ll.Len(), len(c.dict))
	}
}

func Example() {
	type (
		key   struct{ K string }
		value struct{ V int }
	)

	var c = New[key, value](3)

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

	fmt.Printf("dict: %d\n", len(c.dict))
	values := make([]string, 0, len(c.dict))
	for k, v := range c.dict {
		values = append(values, fmt.Sprintf(`"%v: %v"`, k, v.Value))
	}

	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })

	fmt.Printf("%v\nlist: %d\n", values, c.ll.Len())
	elem := c.ll.Back()
	for elem != nil {
		fmt.Printf("%v\n", elem.Value)
		c.ll.Remove(elem)
		elem = c.ll.Back()
	}

	// Output:
	// dict: 2
	// ["{5}: {{5} {5}}" "{six}: {{six} {6}}"]
	// list: 2
	// {{5} {5}}
	// {{six} {6}}
}
