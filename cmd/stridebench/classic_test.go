package main

import (
	"fmt"

	"Skewed_Cache_G19_14042/pkg/benchmark"
	"Skewed_Cache_G19_14042/pkg/cache"
)

// CacheUnderTest is the shared surface both cache implementations expose.
// Because SkewedCache and SetAssociativeCache have identical method
// signatures, this interface lets us write ONE test loop and run it
// against either cache without duplicating code.
type CacheUnderTest interface {
	Lookup(addr uint64) (hit bool, way int, set int)
	Replace(addr uint64, data []byte) (way int, set int)
	GetSetUtilization() float64
}

// runTrace feeds a full address trace through a cache: on every access it
// looks the address up, and on a miss it installs the block via Replace -
// exactly what a real memory hierarchy does.
func runTrace(c CacheUnderTest, blockSize int, trace []uint64) {
	dummyData := make([]byte, blockSize)
	for _, addr := range trace {
		hit, _, _ := c.Lookup(addr)
		if !hit {
			c.Replace(addr, dummyData)
		}
	}
}

func main() {
	const numSets = 64
	const associativity = 4
	const blockSize = 16
	const traceLength = 5000

	gen := benchmark.NewAccessGenerator(blockSize, numSets)

	scenarios := map[string][]uint64{
		"Stride (thrashing pattern)": gen.GenerateStrideAccess(traceLength, 1),
		"Locality (hot window)":      gen.GenerateLocalityAccess(traceLength, associativity*2),
		"Random":                     gen.GenerateRandomAccess(traceLength, uint64(numSets*associativity*blockSize*8)),
	}

	// IMPORTANT: build a fresh pair of caches per scenario. Reusing the
	// same cache instance across scenarios would let leftover state
	// (and inflated OccupiedBlocks/Tracker history) from one pattern
	// contaminate the results of the next, making the comparison unfair.
	for name, trace := range scenarios {
		fmt.Printf("\n################ Scenario: %s ################\n", name)

		skewed := cache.NewSkewedCache(numSets, associativity, blockSize)
		classic := cache.NewSetAssociativeCache(numSets, associativity, blockSize)

		runTrace(skewed, blockSize, trace)
		runTrace(classic, blockSize, trace)

		printSummary("Skewed Cache          ", skewed.TotalAccesses, skewed.TotalHits,
			skewed.CompulsoryMisses, skewed.ConflictMisses, skewed.CapacityMisses, skewed.GetSetUtilization())
		printSummary("Classic Set-Assoc Cache", classic.TotalAccesses, classic.TotalHits,
			classic.CompulsoryMisses, classic.ConflictMisses, classic.CapacityMisses, classic.GetSetUtilization())
	}
}

func printSummary(label string, accesses, hits, compulsory, conflict, capacity uint64, util float64) {
	var hitRate float64
	if accesses > 0 {
		hitRate = float64(hits) / float64(accesses) * 100
	}
	fmt.Printf("%s | HitRate=%.2f%% | Compulsory=%d Conflict=%d Capacity=%d | Utilization=%.2f%%\n",
		label, hitRate, compulsory, conflict, capacity, util*100)
}
