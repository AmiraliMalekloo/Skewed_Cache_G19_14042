package benchmark

import (
	"math/rand"
)

// PatternType specifies the type of memory access pattern.
type PatternType int

const (
	PatternStride PatternType = iota
	PatternLocality
	PatternRandom
)

// AccessGenerator handles generation of trace address sequences.
type AccessGenerator struct {
	BlockSize int
	NumSets   int
}

func NewAccessGenerator(blockSize, numSets int) *AccessGenerator {
	return &AccessGenerator{
		BlockSize: blockSize,
		NumSets:   numSets,
	}
}

// GenerateStrideAccess creates an access pattern with a specific stride to induce thrashing.
func (g *AccessGenerator) GenerateStrideAccess(length int, strideFactor int) []uint64 {
	addrs := make([]uint64, length)
	// TODO: Generate addresses with stride equal to:
	// Stride = strideFactor * g.NumSets * g.BlockSize
	// This pattern causes severe thrashing in traditional set-associative caches,
	// but a Skewed Cache should disperse the conflicts more evenly.
	return addrs
}

// GenerateLocalityAccess creates access patterns exhibiting high spatial and temporal locality.
func (g *AccessGenerator) GenerateLocalityAccess(length int, workingSetSize int) []uint64 {
	addrs := make([]uint64, length)
	// TODO: Generate accesses concentrated within a bounded working set,
	// combining small strides (spatial locality) and repeated addresses (temporal locality).
	return addrs
}

// GenerateRandomAccess creates a uniformly random access trace.
func (g *AccessGenerator) GenerateRandomAccess(length int, addressSpace uint64) []uint64 {
	addrs := make([]uint64, length)
	// TODO: Generate uniformly distributed random addresses within addressSpace bounds.
	for i := 0; i < length; i++ {
		addrs[i] = uint64(rand.Int64n(int64(addressSpace)))
	}
	return addrs
}

// --- Akita Discrete-Event Simulator Integration ---

// TODO: Add Akita integration structures:
//
// type SkewedCacheComponent struct {
//     *sim.ComponentBase
//     Engine sim.Engine
//     Cache  *cache.SkewedCache
//     TopPort sim.Port
//     BottomPort sim.Port
// }
//
// Handle method:
// func (c *SkewedCacheComponent) Handle(e sim.Event) error {
//     // TODO:
//     // 1. Receive incoming request event (e.g., ReadReq or WriteReq) from TopPort.
//     // 2. Call hit, way, set := c.Cache.Lookup(req.Address).
//     // 3. On miss, invoke c.Cache.Replace.
//     // 4. Construct response event (Rsp) and schedule it via c.Engine.Schedule().
//     return nil
// }