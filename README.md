# WTFCache

It's a simple (non-fancy) LRU cache with generic interface
 without annoying limitations to `comparable` keys.

## Example

```
package main

import (
	"fmt"

	"github.com/irr123/wtfcache"
)

func main() {
	c := wtfcache.New[string, string]().MakeWithLock(1e4)
	c.Set("k1", "v1")
	c.Set("k2", "v2")

	v1, ok := c.Get("k1")
	fmt.Printf("%q, %t", v1, ok)

	c.Del("k2")
}
```

[Live example](https://go.dev/play/p/08aS5zNuMnX)

## Here is some benchmarks:

 - https://github.com/vmihailenco/go-cache-benchmark
 - https://github.com/irr123/go-cache-benchmark-1/tree/wtf

```
$ go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/irr123/wtfcache
cpu: 12th Gen Intel(R) Core(TM) i9-12900H
BenchmarkWTF/Set-20             23604550                49.69 ns/op           15 B/op          1 allocs/op
BenchmarkWTF/Get-20             86443056                14.10 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/irr123/wtfcache      2.466s
```
