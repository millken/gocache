package gocache

import (
	"testing"
	"time"
)

func TestGlobalNoConfig(t *testing.T) {

	a, found := Get("a")
	if found || a != nil {
		t.Error("Getting A found value that shouldn't exist:", a)
	}
}

func TestGlobal(t *testing.T) {
	InitConfig(DefaultConfig)

	a, found := Get("a")
	if found || a != nil {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	b, found := Get("b")
	if found || b != nil {
		t.Error("Getting B found value that shouldn't exist:", b)
	}

	c, found := Get("c")
	if found || c != nil {
		t.Error("Getting C found value that shouldn't exist:", c)
	}

	Set("a", 1, DefaultExpiration)
	Set("b", "b", DefaultExpiration)
	Set("c", 3.5, DefaultExpiration)

	x, found := Get("a")
	if !found {
		t.Error("a was not found while getting a2")
	}
	if x == nil {
		t.Error("x for a is nil")
	} else if a2 := x.(int); a2+2 != 3 {
		t.Error("a2 (which should be 1) plus 2 does not equal 3; value:", a2)
	}

	x, found = Get("b")
	if !found {
		t.Error("b was not found while getting b2")
	}
	if x == nil {
		t.Error("x for b is nil")
	} else if b2 := x.(string); b2+"B" != "bB" {
		t.Error("b2 (which should be b) plus B does not equal bB; value:", b2)
	}

	x, found = Get("c")
	if !found {
		t.Error("c was not found while getting c2")
	}
	if x == nil {
		t.Error("x for c is nil")
	} else if c2 := x.(float64); c2+1.2 != 4.7 {
		t.Error("c2 (which should be 3.5) plus 1.2 does not equal 4.7; value:", c2)
	}
}

func TestGlobal_HSet_HGet(t *testing.T) {

	const k = "k"
	const f = "f"
	const v = "v"

	InitConfig(DefaultConfig)

	HSet(k, f, v)

	x, found := HGet(k, f)

	if !found {
		t.Errorf("HGet[%s][%s] was not found", k, f)
	}
	if x == nil {
		t.Error("x is nil")
	} else if b2 := x.(string); b2+"B" != "vB" {
		t.Errorf("'%s' does not equal to '%s'", b2, v)
	}
}

func TestGlobal_HSet_HGetAll(t *testing.T) {

	const k = "k"
	const f = "f"
	const v = "v"

	InitConfig(DefaultConfig)

	HSet(k, f, v)

	x, found := HGetAll(k)

	if !found {
		t.Errorf("HGet[%s][%s] was not found", k, f)
	}
	if x == nil {
		t.Error("x is nil")
	} else if b2 := x.(map[string]interface{}); b2[f].(string) != v {
		t.Errorf("'%s' does not equal to '%s'", b2[f].(string), v)
	}
}

func TestGlobalDelete(t *testing.T) {
	InitConfig(DefaultConfig)
	Set("foo", "bar", DefaultExpiration)
	Delete("foo")
	x, found := Get("foo")
	if found {
		t.Error("foo was found, but it should have been deleted")
	}
	if x != nil {
		t.Error("x is not nil:", x)
	}
}

func TestGlobalHDel(t *testing.T) {
	InitConfig(DefaultConfig)
	HSet("key", "foo", "bar")
	HDel("key", "foo")
	x, found := HGet("key", "foo")
	if found {
		t.Error("foo was found, but it should have been deleted")
	}
	if x != nil {
		t.Error("x is not nil:", x)
	}
}

func BenchmarkGlobalGetExpiring(b *testing.B) {
	benchmarkGlobalGet(b, 5*time.Minute)
}

func BenchmarkGlobalGetNotExpiring(b *testing.B) {
	benchmarkGlobalGet(b, NoExpiration)
}

func benchmarkGlobalGet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	cf := Config{
		DefaultExpiration: exp,
	}
	InitConfig(cf)
	Set("foo", "bar", DefaultExpiration)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Get("foo")
	}
}

func BenchmarkGlobalHGet(b *testing.B) {
	b.StopTimer()
	InitConfig(DefaultConfig)
	HSet("foo", "bar", "x")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		HGet("foo", "bar")
	}
}

func TestGlobal_LPush_LPop(t *testing.T) {

	const k = "k1"
	InitConfig(DefaultConfig)

	for i := 0; i <= 5; i++ {
		LPush(k, i)
	}
	for i := 5; i >= 0; i-- {
		x, found := LPop(k)

		if !found {
			t.Errorf("LPush[%s] was not found", k)
		}
		if x == nil {
			t.Error("x is nil")
		} else if b2 := x.(int); b2 != i {
			t.Errorf("'%d' does not equal to '%d'", b2, i)
		}

	}
	x, found := LPop(k)

	if found {
		t.Errorf("LPop[%s] was found", k)
	}
	if x != nil {
		t.Error("x is not nil")
	}

}

func BenchmarkGlobal_LPush_LPop(b *testing.B) {
	b.StopTimer()
	const k = "k2"
	InitConfig(DefaultConfig)
	b.StartTimer()
	for i := 0; i <= b.N; i++ {
		LPush(k, i)
	}
	for i := b.N; i >= 0; i-- {
		LPop(k)
	}

}

func TestGlobal_RPush_RPop(t *testing.T) {

	const k = "k3"
	InitConfig(DefaultConfig)

	for i := 0; i <= 5; i++ {
		RPush(k, i)
	}
	for i := 5; i >= 0; i-- {
		x, found := RPop(k)

		if !found {
			t.Errorf("RPush[%s] was not found", k)
		}
		if x == nil {
			t.Error("x is nil")
		} else if b2 := x.(int); b2 != i {
			t.Errorf("'%d' does not equal to '%d'", b2, i)
		}

	}
	x, found := RPop(k)

	if found {
		t.Errorf("RPop[%s] was found", k)
	}
	if x != nil {
		t.Error("x is not nil")
	}

}

func BenchmarkGlobal_RPush_RPop(b *testing.B) {
	b.StopTimer()
	const k = "k4"
	InitConfig(DefaultConfig)
	b.StartTimer()
	for i := 0; i <= b.N; i++ {
		RPush(k, i)
	}
	for i := b.N; i >= 0; i-- {
		RPop(k)
	}

}