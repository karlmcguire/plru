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

func TestClear(t *testing.T) {
	p := NewPolicy(64)
	p.Hit(1)
	p.Clear(1)
	if p.blocks[0] != 0 {
		t.Fatal("Clear not working")
	}
}

func TestEvict(t *testing.T) {
	p := NewPolicy(128)
	for i := 0; i < 100; i++ {
		p.Hit(rand.Uint64() % 128)
	}
	for i := 0; i < 32; i++ {
		if p.Has(p.Evict()) {
			t.Fatal("Evict returning invalid block")
		}
	}
}

func BenchmarkHas(b *testing.B) {
	p := NewPolicy(64)
	b.SetBytes(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Has(0)
		}
	})
}

func BenchmarkHit(b *testing.B) {
	p := NewPolicy(64)
	b.SetBytes(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Hit(1)
		}
	})
}

func BenchmarkClear(b *testing.B) {
	p := NewPolicy(64)
	b.SetBytes(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Clear(1)
		}
	})
}

func BenchmarkEvict(b *testing.B) {
	p := NewPolicy(64)
	b.SetBytes(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Evict()
		}
	})
}
