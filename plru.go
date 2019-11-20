package plru

import (
	"math"
)

const (
	blockSize = 64
	blockMask = blockSize - 1
	blockFull = math.MaxUint64
)

type Policy struct {
	blocks  []uint64
	counter uint64
}

func NewPolicy(size uint64) *Policy {
	return &Policy{
		blocks: make([]uint64, (size+blockMask)/blockSize),
	}
}

func (p *Policy) Has(bit uint64) bool {
	return (p.blocks[bit/blockSize] & (1 << (bit & blockMask))) > 0
}

func (p *Policy) Hit(bit uint64) {
	block := &p.blocks[bit/blockSize]
	*block |= 1 << (bit & blockMask)
	if *block == blockFull {
		*block = 0 | 1<<(bit&blockMask)
	}
}

func (p *Policy) Clear(bit uint64) {
	p.blocks[bit/blockSize] &= 0 << (bit & blockMask)
}

func (p *Policy) Evict() uint64 {
	index := p.counter
	block := &p.blocks[index]
	if p.counter++; p.counter >= uint64(len(p.blocks)) {
		p.counter = 0
	}
	return (index * blockSize) + bitLookup(^*block&(*block+1))
}

func bitLookup(num uint64) uint64 {
	switch num {
	case 0:
		return 0
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
	case 4294967296:
		return 32
	case 8589934592:
		return 33
	case 17179869184:
		return 34
	case 34359738368:
		return 35
	case 68719476736:
		return 36
	case 137438953472:
		return 37
	case 274877906944:
		return 38
	case 549755813888:
		return 39
	case 1099511627776:
		return 40
	case 2199023255552:
		return 41
	case 4398046511104:
		return 42
	case 8796093022208:
		return 43
	case 17592186044416:
		return 44
	case 35184372088832:
		return 45
	case 70368744177664:
		return 46
	case 140737488355328:
		return 47
	case 281474976710656:
		return 48
	case 562949953421312:
		return 49
	case 1125899906842624:
		return 50
	case 2251799813685248:
		return 51
	case 4503599627370496:
		return 52
	case 9007199254740992:
		return 53
	case 18014398509481984:
		return 54
	case 36028797018963968:
		return 55
	case 72057594037927936:
		return 56
	case 144115188075855872:
		return 57
	case 288230376151711744:
		return 58
	case 576460752303423488:
		return 59
	case 1152921504606846976:
		return 60
	case 2305843009213693952:
		return 61
	case 4611686018427387904:
		return 62
	case 9223372036854775808:
		return 63
	case 18446744073709551615:
		return 64
	}
	panic("invalid bit lookup")
	return 0
}
