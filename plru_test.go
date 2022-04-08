package plru

import (
	"math/rand"
	"testing"
	"time"

	gen "github.com/pingcap/go-ycsb/pkg/generator"
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
	p := NewPolicy(129)
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
	p.Hit(2)
	p.Del(1)
	if p.Has(1) {
		t.Fatal("Clear not working")
	}
	if !p.Has(2) {
		t.Fatal("Clear other bit")
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
			t.Fatal("lookup should panic")
		}
	}()
	lookup(3)
}

const (
	benchPolicySize = 1e6
)

func benchAccess() (bits [benchPolicySize]uint64) {
	z := gen.NewScrambledZipfian(0, benchPolicySize-1, gen.ZipfianConstant)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < benchPolicySize; i++ {
		bits[i] = uint64(z.Next(r))
	}
	return
}

func BenchmarkHas(b *testing.B) {
	a := benchAccess()
	b.Run("single", func(b *testing.B) {
		b.SetBytes(1)
		p := NewPolicy(benchPolicySize)
		for n := 0; n < b.N; n++ {
			p.Has(a[n%benchPolicySize])
		}
	})
	b.Run("concurrent", func(b *testing.B) {
		b.SetBytes(1)
		p := NewPolicy(benchPolicySize)
		b.RunParallel(func(pb *testing.PB) {
			for n := 0; pb.Next(); n++ {
				p.Has(a[n%benchPolicySize])
			}
		})
	})
}

func BenchmarkHit(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		b.SetBytes(1)
		p := NewPolicy(benchPolicySize)
		for n := 0; n < b.N; n++ {
			p.Hit(1)
		}
	})
	b.Run("concurrent", func(b *testing.B) {
		b.SetBytes(1)
		p := NewPolicy(benchPolicySize)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				p.Hit(1)
			}
		})
	})
}

func BenchmarkClear(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		b.SetBytes(1)
		p := NewPolicy(benchPolicySize)
		for n := 0; n < b.N; n++ {
			p.Del(1)
		}
	})
	b.Run("concurrent", func(b *testing.B) {
		b.SetBytes(1)
		p := NewPolicy(benchPolicySize)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				p.Del(1)
			}
		})
	})
}

func BenchmarkEvict(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		b.SetBytes(1)
		p := NewPolicy(benchPolicySize)
		for n := 0; n < b.N; n++ {
			p.Evict()
		}
	})
	b.Run("concurrent", func(b *testing.B) {
		b.SetBytes(1)
		p := NewPolicy(benchPolicySize)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				p.Evict()
			}
		})
	})
}
