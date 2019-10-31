package gocache

import (
	"runtime"
	"time"
)

var instance *Cache

func init() {
	InitConfig(DefaultConfig)
}
func InitConfig(config Config) {

	if config.DefaultExpiration == 0 {
		config.DefaultExpiration = -1
	}
	c := &cache{
		defaultExpiration: config.DefaultExpiration,
		items:             make(map[string]Item),
		group:             Group{},
	}
	instance = &Cache{c}

	if config.CleanupInterval > 0 {
		runJanitor(c, config.CleanupInterval)
		runtime.SetFinalizer(instance, stopJanitor)
	}

}

func Increment(k string, n int64) error {
	return instance.Increment(k, n)
}

func Decrement(k string, n int64) error {
	return instance.Decrement(k, n)
}

func Set(k string, x interface{}, d time.Duration) {
	instance.Set(k, x, d)
}

func Get(k string) (interface{}, bool) {
	return instance.Get(k)
}

func Delete(k string) {
	instance.Delete(k)
}

func HSet(k, f string, x interface{}) {
	instance.HSet(k, f, x)
}

func HGet(k, f string) (interface{}, bool) {
	return instance.HGet(k, f)
}

func HGetAll(k string) (interface{}, bool) {
	return instance.HGetAll(k)
}

func HDel(k, f string) {
	instance.HDel(k, f)
}

func LPush(k string, x interface{}) {
	instance.LPush(k, x)
}

func LPop(k string) (interface{}, bool) {
	return instance.LPop(k)
}

func RPush(k string, x interface{}) {
	instance.RPush(k, x)
}

func RPop(k string) (interface{}, bool) {
	return instance.RPop(k)
}

func OnEvicted(f func(string, interface{})) {
	instance.OnEvicted(f)
}

func TTL(k string, d time.Duration) {
	instance.TTL(k, d)
}

func Memoize(k string, fn func() (interface{}, error), d time.Duration) (interface{}, error) {
	return instance.Memoize(k, fn, d)
}

// Copies all unexpired items in the cache into a new map and returns it.
func Items() map[string]Item {
	return instance.Items()
}

// Returns the number of items in the cache. This may include items that have
// expired, but have not yet been cleaned up.
func ItemCount() int {
	return instance.ItemCount()
}

// Delete all items from the cache.
func Flush() {
	instance.Flush()
}
