package gocache

import (
	"testing"
	"time"
)

type TestStruct struct {
	Num      int
	Children []*TestStruct
}

func TestCache(t *testing.T) {
	tc := NewCache(DefaultConfig)

	a, found := tc.Get("a")
	if found || a != nil {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	tc.Set("a", 1, DefaultExpiration)

	x, found := tc.Get("a")
	if !found {
		t.Error("a was not found while getting a2")
	}
	if x == nil {
		t.Error("x for a is nil")
	} else if a2 := x.(int); a2+2 != 3 {
		t.Error("a2 (which should be 1) plus 2 does not equal 3; value:", a2)
	}
}

func TestCache_Memoize(t *testing.T) {
	tc := NewCache(DefaultConfig)

	a, err := tc.Memoize("a", func() (interface{}, error) {
		return 1, nil
	}, 1)

	if err != nil || a.(int) != 1 {
		t.Error("memoize error :", a)
	}
	time.Sleep(2 * time.Second)
	x, found := tc.Get("a")
	if found || x != nil {
		t.Error("a was found while getting a")
	}
}

func TestCacheTimes(t *testing.T) {
	var found bool

	cf := Config{
		DefaultExpiration: 50 * time.Millisecond,
	}
	tc := NewCache(cf)
	tc.Set("a", 1, DefaultExpiration)
	tc.Set("b", 2, NoExpiration)
	tc.Set("c", 3, 20*time.Millisecond)
	tc.Set("d", 4, 70*time.Millisecond)

	<-time.After(25 * time.Millisecond)
	_, found = tc.Get("c")
	if found {
		t.Error("Found c when it should have been automatically deleted")
	}

	<-time.After(30 * time.Millisecond)
	_, found = tc.Get("a")
	if found {
		t.Error("Found a when it should have been automatically deleted")
	}

	_, found = tc.Get("b")
	if !found {
		t.Error("Did not find b even though it was set to never expire")
	}

	_, found = tc.Get("d")
	if !found {
		t.Error("Did not find d even though it was set to expire later than the default")
	}

	<-time.After(20 * time.Millisecond)
	_, found = tc.Get("d")
	if found {
		t.Error("Found d when it should have been automatically deleted (later than the default)")
	}
}

func TestCache_HSet_HGet(t *testing.T) {

	const k = "k"
	const f = "f"
	const v = "v"

	tc := NewCache(DefaultConfig)

	tc.HSet(k, f, v)

	x, found := tc.HGet(k, f)

	if !found {
		t.Errorf("HGet[%s][%s] was not found", k, f)
	}
	if x == nil {
		t.Error("x is nil")
	} else if b2 := x.(string); b2+"B" != "vB" {
		t.Errorf("'%s' does not equal to '%s'", b2, v)
	}
}

func TestCache_HSet_HGetAll(t *testing.T) {

	const k = "k"
	const f = "f"
	const v = "v"

	tc := NewCache(DefaultConfig)

	tc.HSet(k, f, v)

	x, found := tc.HGetAll(k)

	if !found {
		t.Errorf("HGet[%s][%s] was not found", k, f)
	}
	if x == nil {
		t.Error("x is nil")
	} else if b2 := x.(map[string]interface{}); b2[f].(string) != v {
		t.Errorf("'%s' does not equal to '%s'", b2[f].(string), v)
	}
}

func TestDelete(t *testing.T) {
	tc := NewCache(DefaultConfig)
	tc.Set("foo", "bar", DefaultExpiration)
	tc.Delete("foo")
	x, found := tc.Get("foo")
	if found {
		t.Error("foo was found, but it should have been deleted")
	}
	if x != nil {
		t.Error("x is not nil:", x)
	}
}

func TestHDel(t *testing.T) {
	tc := NewCache(DefaultConfig)
	tc.HSet("key", "foo", "bar")
	tc.HDel("key", "foo")
	x, found := tc.HGet("key", "foo")
	if found {
		t.Error("foo was found, but it should have been deleted")
	}
	if x != nil {
		t.Error("x is not nil:", x)
	}
}

func BenchmarkCacheGetExpiring(b *testing.B) {
	benchmarkCacheGet(b, 5*time.Minute)
}

func BenchmarkCacheGetNotExpiring(b *testing.B) {
	benchmarkCacheGet(b, NoExpiration)
}

func benchmarkCacheGet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	cf := Config{
		DefaultExpiration: exp,
	}
	tc := NewCache(cf)
	tc.Set("foo", "bar", DefaultExpiration)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Get("foo")
	}
}

func BenchmarkCacheHGet(b *testing.B) {
	b.StopTimer()
	tc := NewCache(DefaultConfig)
	tc.HSet("foo", "bar", "x")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.HGet("foo", "bar")
	}
}
