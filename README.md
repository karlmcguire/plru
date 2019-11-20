# plru

Pseudo-LRU implementation using 1-bit per entry and achieving hit ratios within
1-2% of full LRU. This method is commonly refered to as "PLRUm" because each bit
serves as a MRU flag for each cache entry.

Academic literature where PLRUm is mentioned:

* [Performance evaluation of cache replacement policies for the SPEC CPU2000 benchmark suite](https://dl.acm.org/citation.cfm?id=986601)
    * Shows that PLRUm is a very good approximation of LRU
* [Study of Different Cache Line Replacement Algorithms in Embedded Systems](https://people.kth.se/~ingo/MasterThesis/ThesisDamienGille2007.pdf)
    * Shows that PLRUm usually outperforms other PLRU algorithms
    * In some cases, PLRUm *outperforms* LRU

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
