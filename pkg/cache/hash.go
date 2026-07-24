package cache

import "math/bits"

// we assume NumSets is a power of two, as required by the
// bit-manipulation hash functions below.
func log2(n int) int {
	return bits.Len(uint(n)) - 1
}
// HashXORFold implements the first hash function based on XOR bit folding.
func HashXORFold(blockAddr uint64, nSets int) int {
	if nSets <= 1 {
		return 0
	}
 
	shiftAmount := uint(log2(nSets))
	folded := blockAddr ^ (blockAddr >> shiftAmount)
 
	return int(folded) & (nSets - 1)
}

const fibonacciConstant uint64 = 11400714819323198485

// HashFibonacci implements the second hash function based on Fibonacci/Golden Ratio Hashing.
func HashFibonacci(blockAddr uint64, nSets int) int {
	if nSets <= 1 {
		return 0
	}
 
	shiftAmount := uint(64 - log2(nSets))
	hashed := blockAddr * fibonacciConstant
	index := hashed >> shiftAmount
 
	return int(index) & (nSets - 1)
}

// MapToSet maps a block address to a Set Index based on the given Way ID.
func (c *SkewedCache) MapToSet(blockAddr uint64, wayID int) int {
	rotateAmount := uint(wayID * 11)
	permuted := bits.RotateLeft64(blockAddr, int(rotateAmount))
 
	if wayID%2 == 0 {
		return HashXORFold(permuted, c.NumSets)
	}
	return HashFibonacci(permuted, c.NumSets)
}