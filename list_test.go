package gocache

import "testing"

func TestCache_LPush_LPop(t *testing.T) {

	const k = "k"
	tc := NewCache(DefaultConfig)

	for i := 0; i <= 5; i++ {
		tc.LPush(k, i)
	}
	for i := 5; i >= 0; i-- {
		x, found := tc.LPop(k)

		if !found {
			t.Errorf("LPush[%s] was not found", k)
		}
		if x == nil {
			t.Error("x is nil")
		} else if b2 := x.(int); b2 != i {
			t.Errorf("'%d' does not equal to '%d'", b2, i)
		}

	}
	x, found := tc.LPop(k)

	if found {
		t.Errorf("LPop[%s] was found", k)
	}
	if x != nil {
		t.Error("x is not nil")
	}

}

func BenchmarkCache_LPush_LPop(b *testing.B) {
	b.StopTimer()
	const k = "k"
	tc := NewCache(DefaultConfig)
	b.StartTimer()
	for i := 0; i <= b.N; i++ {
		tc.LPush(k, i)
	}
	for i := b.N; i >= 0; i-- {
		tc.LPop(k)
	}

}

func TestCache_RPush_RPop(t *testing.T) {

	const k = "k"
	tc := NewCache(DefaultConfig)

	for i := 0; i <= 5; i++ {
		tc.RPush(k, i)
	}
	for i := 5; i >= 0; i-- {
		x, found := tc.RPop(k)

		if !found {
			t.Errorf("RPush[%s] was not found", k)
		}
		if x == nil {
			t.Error("x is nil")
		} else if b2 := x.(int); b2 != i {
			t.Errorf("'%d' does not equal to '%d'", b2, i)
		}

	}
	x, found := tc.RPop(k)

	if found {
		t.Errorf("RPop[%s] was found", k)
	}
	if x != nil {
		t.Error("x is not nil")
	}

}

func BenchmarkCache_RPush_RPop(b *testing.B) {
	b.StopTimer()
	const k = "k"
	tc := NewCache(DefaultConfig)
	b.StartTimer()
	for i := 0; i <= b.N; i++ {
		tc.RPush(k, i)
	}
	for i := b.N; i >= 0; i-- {
		tc.RPop(k)
	}

}
