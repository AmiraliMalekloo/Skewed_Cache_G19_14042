package cache

import "sync"

// BlockTracker is a helper structure to identify Compulsory Misses.
// It keeps track of all block addresses that have been accessed at least once.
type BlockTracker struct {
	mu sync.RWMutex
	
	// TODO [Option 1 - Simple]: Use a map for exact tracking of seen addresses
	seenBlocks map[uint64]bool

	// TODO [Option 2 - Advanced/Memory Efficient]: Implement a Bloom Filter
	// to reduce memory consumption during long simulation runs.
	// Note: Bloom Filters may yield False Positives, so select based on accuracy needs.
}

func NewBlockTracker() *BlockTracker {
	// TODO: Initialize the map or Bloom Filter
	return nil
}

// IsFirstAccess checks if the given block address is being accessed for the very first time.
func (bt *BlockTracker) IsFirstAccess(blockAddr uint64) bool {
	// TODO:
	// 1. Acquire appropriate read/write locks.
	// 2. Query blockAddr in the map or Bloom Filter.
	// 3. If not found, insert it and return true (indicating a Compulsory Miss).
	// 4. If found, return false.
	return false
}