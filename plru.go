// A Pseudo-LRU implementation using 1-bit per entry and hit ratio performance
// nearly identical to full LRU.
package plru

import (
	"math"
	"sync/atomic"
)

const (
	bSize = 32
	bMask = bSize - 1
	bFull = math.MaxUint32
)

type Policy struct {
	blocks []uint32
	cursor uint64
	size   uint64
}

// NewPolicy returns an empty Policy where size is the number of entries to
// track. The size param is rounded up to the next multiple of 32.
func NewPolicy(size uint64) *Policy {
	size = ((size + 32) >> 5) << 5
	return &Policy{
		blocks: make([]uint32, (size+bMask)/bSize),
		size:   ((size + bMask) / bSize),
	}
}

// Has returns true if the bit is set (1) and false if not (0).
func (p *Policy) Has(bit uint64) bool {
	return (atomic.LoadUint32(&p.blocks[bit/bSize]) & (1 << (bit & bMask))) > 0
}

// Hit sets the bit to 1 and clears the other bits in the block if capacity is
// reached.
func (p *Policy) Hit(bit uint64) {
	block := &p.blocks[bit/bSize]
hit:
	o := atomic.LoadUint32(block)
	n := o | 1<<(bit&bMask)
	if n == bFull {
		n = 0 | 1<<(bit&bMask)
	}
	if !atomic.CompareAndSwapUint32(block, o, n) {
		goto hit
	}
}

// Del sets the bit to 0.
func (p *Policy) Del(bit uint64) {
	block := &p.blocks[bit/bSize]
del:
	o := atomic.LoadUint32(block)
	n := o & 0 << (bit & bMask)
	if !atomic.CompareAndSwapUint32(block, o, n) {
		goto del
	}
}

// Evict returns a LRU bit that you can later pass to Hit.
func (p *Policy) Evict() uint64 {
	i := (atomic.AddUint64(&p.cursor, 1) - 1) % p.size
	block := atomic.LoadUint32(&p.blocks[i])
	return (i * bSize) + lookup(^block&(block+1))
}

func lookup(b uint32) uint64 {
	switch b {
	case 1:
		return 0
	case 2:
		return 1
	case 4:
		return 2
	case 8:
		return 3
	case 16:
		return 4
	case 32:
		return 5
	case 64:
		return 6
	case 128:
		return 7
	case 256:
		return 8
	case 512:
		return 9
	case 1024:
		return 10
	case 2048:
		return 11
	case 4096:
		return 12
	case 8192:
		return 13
	case 16384:
		return 14
	case 32768:
		return 15
	case 65536:
		return 16
	case 131072:
		return 17
	case 262144:
		return 18
	case 524288:
		return 19
	case 1048576:
		return 20
	case 2097152:
		return 21
	case 4194304:
		return 22
	case 8388608:
		return 23
	case 16777216:
		return 24
	case 33554432:
		return 25
	case 67108864:
		return 26
	case 134217728:
		return 27
	case 268435456:
		return 28
	case 536870912:
		return 29
	case 1073741824:
		return 30
	case 2147483648:
		return 31
	}
	panic("invalid bit lookup")
}
