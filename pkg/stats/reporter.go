package stats

import (
	"encoding/csv"
	"fmt"
	"os"

	"skewed-cache-sim/pkg/cache"
)

type Reporter struct {
	Cache *cache.SkewedCache
}

func NewReporter(c *cache.SkewedCache) *Reporter {
	return &Reporter{Cache: c}
}

// PrintConsoleReport formats and prints summary statistics to stdout.
func (r *Reporter) PrintConsoleReport() {
	// TODO: Format and print all metrics:
	// - Total Accesses
	// - Total Hits & Hit Ratio
	// - Compulsory Misses (%)
	// - Conflict Misses (%)
	// - Capacity Misses (%)
	// - Set Utilization (%)
	fmt.Println("=== Skewed Cache Simulation Report ===")
	// ...
}

// ExportToCSV writes simulation statistics to a CSV file.
func (r *Reporter) ExportToCSV(filePath string) error {
	// TODO: 
	// 1. Create target file (os.Create).
	// 2. Write CSV header:
	//    "TotalAccesses,TotalHits,HitRate,CompulsoryMisses,ConflictMisses,CapacityMisses,SetUtilization"
	// 3. Extract metrics from r.Cache and write them as a CSV row.
	// 4. Flush and close the file handle.
	return nil
}