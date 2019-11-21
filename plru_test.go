package plru

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestHas(t *testing.T) {
	p := NewPolicy(128)
	p.Hit(1)
	p.Hit(64)
	if p.Has(0) || !p.Has(1) || !p.Has(64) {
		t.Fatal("Hit or Has not working")
	}
}

func TestHit(t *testing.T) {
	p := NewPolicy(128)
	if p.Has(0) {
		t.Fatal("Has not working")
	}
	p.Hit(0)
	if !p.Has(0) {
		t.Fatal("Hit not working")
	}
	p.Hit(120)
	if !p.Has(120) {
		t.Fatal("Hit not working")
	}
}

func TestDel(t *testing.T) {
	p := NewPolicy(64)
	p.Hit(1)
	p.Del(1)
	if p.blocks[0] != 0 {
		t.Fatal("Clear not working")
	}
}

func TestEvict(t *testing.T) {
	p := NewPolicy(128)
	for i := 0; i < 128; i++ {
		victim := p.Evict()
		if p.Has(victim) {
			t.Fatal("Evict returning used block")
		}
		p.Hit(p.Evict())
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("bitLookup should panic")
		}
	}()
	bitLookup(3)
}

func BenchmarkHas(b *testing.B) {
	b.SetBytes(1)
	p := NewPolicy(64)
	for n := 0; n < b.N; n++ {
		p.Has(1)
	}
}

func BenchmarkHit(b *testing.B) {
	b.SetBytes(1)
	p := NewPolicy(64)
	for n := 0; n < b.N; n++ {
		p.Hit(1)
	}
}

func BenchmarkClear(b *testing.B) {
	b.SetBytes(1)
	p := NewPolicy(64)
	for n := 0; n < b.N; n++ {
		p.Del(1)
	}
}

func BenchmarkEvict(b *testing.B) {
	b.SetBytes(1)
	p := NewPolicy(64)
	for n := 0; n < b.N; n++ {
		p.Evict()
	}
}
