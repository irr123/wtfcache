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

 - https://github.com/irr123/go-cache-benchmark/tree/wtfcache
 - https://github.com/irr123/go-cache-benchmark-1/tree/wtf
