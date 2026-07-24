package cache

// HashXORFold implements the first hash function based on XOR bit folding.
func HashXORFold(blockAddr uint64, nSets int) int {
	// TODO: Implement the exact XOR Folding formula:
	// 1. Split blockAddr into L-bit segments (where 2^L = nSets).
	// 2. XOR the segments together to generate a set index in the range [0, nSets-1].
	// 3. Mask the result with (nSets - 1) to guarantee bounds safety.
	return 0
}

// HashFibonacci implements the second hash function based on Fibonacci/Golden Ratio Hashing.
func HashFibonacci(blockAddr uint64, nSets int) int {
	// TODO: Implement the exact Fibonacci Hashing formula:
	// 1. Use the 64-bit Fibonacci constant: 11400714819323198485U (2^64 / GoldenRatio).
	// 2. Multiply blockAddr by this constant.
	// 3. Right-shift the result by (64 - log2(nSets)) to extract the high-order bits.
	return 0
}

// MapToSet maps a block address to a Set Index based on the given Way ID.
func (c *SkewedCache) MapToSet(blockAddr uint64, wayID int) int {
	// TODO: Hash function selection logic based on wayID:
	// In a true Skewed Cache, each way utilizes a different hash function or a skewed signal.
	//
	// Example:
	// switch wayID % 2 {
	// case 0:
	//     return HashXORFold(blockAddr, c.NumSets)
	// case 1:
	//     return HashFibonacci(blockAddr, c.NumSets)
	// default:
	//     ...
	// }
	//
	// You can also rotate or permute blockAddr bits for higher associativity counts.
	return 0
}