package gocache

import (
	"runtime"
	"sync"
	"time"
)

const (
	// For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
	// For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to New() or
	// NewFrom() when the cache was created (e.g. 5 minutes.)
	DefaultExpiration time.Duration = 0
)

type Item struct {
	Object     interface{}
	Expiration int64
}

// Returns true if the item has expired.
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

type Cache struct {
	*cache
	// If this is confusing, see the comment at the bottom of New()
}
type cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	onEvicted         func(string, interface{})
	janitor           *janitor
	group             Group
}

var DefaultConfig = Config{
	DefaultExpiration: DefaultExpiration,
	CleanupInterval:   5 * time.Minute,
}

type Config struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

func NewCache(config Config) *Cache {
	if config.DefaultExpiration == 0 {
		config.DefaultExpiration = -1
	}
	c := &cache{
		defaultExpiration: config.DefaultExpiration,
		items:             make(map[string]Item),
		group:             Group{},
	}
	C := &Cache{c}

	if config.CleanupInterval > 0 {
		runJanitor(c, config.CleanupInterval)
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}

func (c *cache) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}

	c.mu.Unlock()
}

func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
	return item.Object, true
}

func (c *cache) Delete(k string) {
	c.mu.Lock()
	v, evicted := c.delete(k)
	c.mu.Unlock()
	if evicted {
		c.onEvicted(k, v)
	}
}

func (c *cache) delete(k string) (interface{}, bool) {
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}
	delete(c.items, k)
	return nil, false
}

func (c *cache) HSet(k, f string, x interface{}) {
	var obj map[string]interface{}
	c.mu.Lock()
	item, found := c.items[k]
	if !found {
		obj = make(map[string]interface{})
	} else {

		switch item.Object.(type) {
		case map[string]interface{}:
			obj = item.Object.(map[string]interface{})
		default:
			obj = make(map[string]interface{})

		}
	}
	obj[f] = x
	c.items[k] = Item{
		Object:     obj,
		Expiration: 0, //Hset can not
	}

	c.mu.Unlock()
}

func (c *cache) HGet(k, f string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	obj := item.Object.(map[string]interface{})
	val, found := obj[f]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	c.mu.RUnlock()
	return val, true
}

func (c *cache) HGetAll(k string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	obj := item.Object.(map[string]interface{})
	c.mu.RUnlock()
	return obj, true
}

func (c *cache) HDel(k, f string) {
	c.mu.Lock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return
	}

	obj := item.Object.(map[string]interface{})
	_, found = obj[f]
	if !found {
		c.mu.Unlock()
		return
	}
	delete(obj, f)
	c.items[k] = Item{
		Object:     obj,
		Expiration: item.Expiration,
	}
	c.mu.Unlock()
	return
}

// Sets an (optional) function that is called with the key and value when an
// item is evicted from the cache. (Including when it is deleted manually, but
// not when it is overwritten.) Set to nil to disable.
func (c *cache) OnEvicted(f func(string, interface{})) {
	c.mu.Lock()
	c.onEvicted = f
	c.mu.Unlock()
}

func (c *cache) TTL(k string, d time.Duration) {
	var e int64
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	item, found := c.items[k]
	if !found {
		c.mu.Unlock()
		return
	}
	item.Expiration = e
	c.items[k] = item

	c.mu.Unlock()
}

type keyAndValue struct {
	key   string
	value interface{}
}

func (c *cache) DeleteExpired() {
	var evictedItems []keyAndValue
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		// "Inlining" of expired
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := c.delete(k)
			if evicted {
				evictedItems = append(evictedItems, keyAndValue{k, ov})
			}
		}
	}
	c.mu.Unlock()
	for _, v := range evictedItems {
		c.onEvicted(v.key, v.value)
	}
}

// Memoize executes and returns the results of the given function, unless there was a cached value of the same key.
// Only one execution is in-flight for a given key at a time.
// The boolean return value indicates whether v was previously stored.
func (c *cache) Memoize(k string, fn func() (interface{}, error), d time.Duration) (interface{}, error) {
	// Check cache
	value, found := c.Get(k)
	if found {
		return value, nil
	}

	value, err := c.group.Do(k, func() (interface{}, error) {
		data, innerErr := fn()

		if innerErr == nil {
			c.Set(k, data, d)
		}

		return data, innerErr
	})
	return value, err
}
