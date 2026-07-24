package cache

import "sync"

// CacheLine represents a single block/line in the cache memory.
type CacheLine struct {
	Tag   uint64
	Valid bool
	Dirty bool
	Data  []byte

	LastAccessTime int64
}

// SkewedCache is the main data structure for the skewed cache simulator.
type SkewedCache struct {
	mu            sync.Mutex
	NumSets       int // Number of sets (S)
	Associativity int // Number of ways (W)
	BlockSize     int // Block size in bytes

	// 2D array of cache lines: Sets x Ways
	// In a Skewed Cache, each way is addressed using a distinct hash function.
	Sets [][]CacheLine

	// First-access tracker used to identify compulsory misses
	Tracker *BlockTracker

	// Statistical Counters :
	TotalAccesses    uint64
	TotalHits        uint64
	CompulsoryMisses uint64
	ConflictMisses   uint64
	CapacityMisses   uint64

	// Total count of occupied (Valid) blocks across the entire cache for setUtilization
	OccupiedBlocks uint64
}

// NewSkewedCache constructs and initializes a new SkewedCache instance.
func NewSkewedCache(numSets, associativity, blockSize int) *SkewedCache {
	sets := make([][]CacheLine, numSets)
	for s := 0; s < numSets; s++ {
		sets[s] = make([]CacheLine, associativity)
		for w := 0; w < associativity; w++ {
			sets[s][w] = CacheLine{
				Data: make([]byte, blockSize),
			}
		}
	}

	return &SkewedCache{
		NumSets:       numSets,
		Associativity: associativity,
		BlockSize:     blockSize,
		Sets:          sets,
		Tracker:       NewBlockTracker(),
	}
}
