package gocache

import (
	"runtime"
	"sync"
	"time"
)

type singleton struct {
}

var instance *Cache
var once sync.Once
var initConfig = false

func InitConfig(config Config) {

	once.Do(func() {
		initConfig = true
		if config.DefaultExpiration == 0 {
			config.DefaultExpiration = -1
		}
		c := &cache{
			defaultExpiration: config.DefaultExpiration,
			items:             make(map[string]Item),
		}
		instance = &Cache{c}

		if config.CleanupInterval > 0 {
			runJanitor(c, config.CleanupInterval)
			runtime.SetFinalizer(instance, stopJanitor)
		}
	})

}

func getInstance() *Cache {
	if !initConfig {
		InitConfig(DefaultConfig)
	}
	return instance
}

func Set(k string, x interface{}, d time.Duration) {
	getInstance().Set(k, x, d)
}

func Get(k string) (interface{}, bool) {
	return getInstance().Get(k)
}

func Delete(k string) {
	getInstance().Delete(k)
}

func HSet(k, f string, x interface{}) {
	getInstance().HSet(k, f, x)
}

func HGet(k, f string) (interface{}, bool) {
	return getInstance().HGet(k, f)
}

func HGetAll(k string) (interface{}, bool) {
	return getInstance().HGetAll(k)
}

func HDel(k, f string) {
	getInstance().HDel(k, f)
}

func LPush(k string, x interface{}) {
	getInstance().LPush(k, x)
}

func LPop(k string) (interface{}, bool) {
	return getInstance().LPop(k)
}

func RPush(k string, x interface{}) {
	getInstance().RPush(k, x)
}

func RPop(k string) (interface{}, bool) {
	return getInstance().RPop(k)
}

func OnEvicted(f func(string, interface{})) {
	getInstance().OnEvicted(f)
}

func TTL(k string, d time.Duration) {
	getInstance().TTL(k, d)
}
