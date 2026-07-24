package main

import (
	"fmt"

	"skewed-cache-sim/pkg/benchmark"
	"skewed-cache-sim/pkg/cache"
	"skewed-cache-sim/pkg/stats"
)

func main() {
	// 1. Simulation parameters
	numSets := 64
	associativity := 4
	blockSize := 64 // 64 Bytes

	// 2. Instantiate Skewed Cache
	skewedCache := cache.NewSkewedCache(numSets, associativity, blockSize)

	// 3. Instantiate workload generator
	gen := benchmark.NewAccessGenerator(blockSize, numSets)

	// TODO: Select workload pattern and generate access trace
	// trace := gen.GenerateStrideAccess(10000, 1)
	// trace := gen.GenerateLocalityAccess(10000, 512)
	trace := gen.GenerateRandomAccess(10000, 1024*1024)

	// 4. Run standalone simulation loop (or drive through Akita engine)
	fmt.Println("Starting simulation...")
	for _, addr := range trace {
		hit, _, _ := skewedCache.Lookup(addr)
		if !hit {
			// On miss, fetch and replace line
			skewedCache.Replace(addr, make([]byte, blockSize))
		}
	}

	// 5. Output statistics
	reporter := stats.NewReporter(skewedCache)
	reporter.PrintConsoleReport()
	
	// TODO: Save CSV report for analysis in the project report
	// reporter.ExportToCSV("results_stride.csv")
}