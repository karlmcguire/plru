# plru
[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg)](https://godoc.org/github.com/karlmcguire/plru)
[![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-green.svg)](https://goreportcard.com/report/github.com/karlmcguire/plru)
[![Coverage](https://img.shields.io/badge/coverage-100%25-ff69b4.svg)](https://gocover.io/karlmcguire/plru)

Pseudo-LRU implementation using 1-bit per entry and achieving hit ratios within
1-2% of full LRU. This method is commonly refered to as "PLRUm" because each bit
serves as a MRU flag for each cache entry.

Academic literature where PLRUm is mentioned:

* [Performance evaluation of cache replacement policies for the SPEC CPU2000 benchmark suite](https://dl.acm.org/citation.cfm?id=986601)
    * Shows that PLRUm is a very good approximation of LRU
* [Study of Different Cache Line Replacement Algorithms in Embedded Systems](https://people.kth.se/~ingo/MasterThesis/ThesisDamienGille2007.pdf)
    * Shows that PLRUm usually outperforms other PLRU algorithms
    * In some cases, PLRUm *outperforms* LRU

**NOTE**: This is a small experiment repo and everything's still up in the air.
Future plans include trying to implement a full cache out of this where it
handles collisions and increases data locality.

# usage

This library is intended to be small and flexible for use within full cache
implementations. Therefore, it is not safe for concurrent use out of the box and
it does nothing to handle collisions. It makes sense to use this within a mutex
lock close to a hashmap.

```go
// create a new Policy tracking 1024 entries
p := NewPolicy(1024)

// on a Get operation, call policy.Hit() for the cache entry
p.Hit(1)

// when the cache is full, call policy.Evict() to get a LRU bit
victim := p.Evict()

// add some things to the victim location
// ...

// call policy.Hit() on the new entry to flag it as MRU
p.Hit(victim)
```

# about

This PLRUm implementation uses `uint64` blocks and uses round-robin for
selecting which block to evict from (performs better than random selection).

## 1. empty state

Before being populated, each block is 0.

<p align="center">
    <img src="https://karlmcguire.com/images/plru_1.svg">
</p>

## 2. adding items

Adding new items behaves just like a hit operation in which the bit is set to 1
because of the recent access.

After adding items A, B, C, D, E, F, and G, the PLRUm state looks like this:

<p align="center">
    <img src="https://karlmcguire.com/images/plru_2.svg">
</p>

## 3. reaching capacity

The eviction operation requires 0 bits to be present and sample from. To
prevent all bits from being 1 and a potential deadlock situation, all bits
except the newest are set to 0 when capacity is reached.

After adding item H, the PLRUm state looks like this:

<p align="center">
    <img src="https://karlmcguire.com/images/plru_3.svg">
</p>

## 4. hit

A hit just sets the bit to 1.

After hitting item D, the PLRUm state looks like this:

<p align="center">
    <img src="https://karlmcguire.com/images/plru_5.svg">
</p>

## 5. eviction

The eviction process returns the location of the first 0 bit found.

After eviction and adding item I, the PLRUm state looks like this:

<p align="center">
    <img src="https://karlmcguire.com/images/plru_6.svg">
</p>
