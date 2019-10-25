# gocache [![Build Status](https://travis-ci.org/millken/groupcache.svg?branch=master)](https://travis-ci.org/millken/groupcache)

go-ache is an in-memory key:value store/cache similar to redis that is
suitable for applications running on a single machine. Its major advantage is
that, being essentially a thread-safe `map[string]interface{}` with expiration
times, it doesn't need to serialize or transmit its contents over the network.

Any object can be stored, for a given duration or forever, and the cache can be
safely used by multiple goroutines.

### Installation

`go get github.com/millken/gocache`

### Usage
```go
import (
	"fmt"
	"github.com/millken/gocache"
	"time"
)

func main() {
	c := gocache.NewCache(gocache.DefaultConfig)

	// Set the value of the key "foo" to "bar", with the default expiration time
	c.Set("foo", "bar", gocache.DefaultExpiration)

	// Set the value of the key "baz" to 42, with no expiration time
	// (the item won't be removed until it is re-set, or removed using
	// c.Delete("baz")
	c.Set("baz", 42, gocache.NoExpiration)

	// Get the string associated with the key "foo" from the cache
	foo, found := c.Get("foo")
	if found {
		fmt.Println(foo)
	}
}
```

