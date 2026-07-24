package cache

import "sync"

// CacheLine represents a single block/line in the cache memory.
type CacheLine struct {
	Tag      uint64
	Valid    bool
	Dirty    bool
	Data     []byte
	
	// TODO: Define LRU replacement policy variables/bits here.
	// Note: In a Skewed Cache, a block can reside in different sets depending on the way.
	// You can use a global timestamp or counter to track the last access time.
	LastAccessTime int64
}

// SkewedCache is the main data structure for the skewed cache simulator.
type SkewedCache struct {
	mu           sync.Mutex
	NumSets      int // Number of sets (S)
	Associativity int // Number of ways (W)
	BlockSize    int // Block size in bytes

	// 2D array of cache lines: Sets x Ways
	// Note: In a Skewed Cache, each way is addressed using a distinct hash function.
	Sets [][]CacheLine

	// First-access tracker used to identify compulsory misses
	Tracker *BlockTracker

	// --- Statistical Counters ---
	TotalAccesses    uint64
	TotalHits        uint64
	CompulsoryMisses uint64
	ConflictMisses   uint64
	CapacityMisses   uint64
	
	// Total count of occupied (Valid) blocks across the entire cache for setUtilization
	OccupiedBlocks   uint64
}

// NewSkewedCache constructs and initializes a new SkewedCache instance.
func NewSkewedCache(numSets, associativity, blockSize int) *SkewedCache {
	// TODO: Initialize the 2D Sets matrix, allocate memory for each Way,
	// and instantiate a new BlockTracker instance.
	return nil
}