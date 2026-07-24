package cache

import "sync"

// BlockTracker is a helper structure to identify Compulsory Misses.
// It keeps track of all block addresses that have been accessed at least once.
type BlockTracker struct {
	mu sync.RWMutex
	
	seenBlocks map[uint64]bool
}

func NewBlockTracker() *BlockTracker {
	return &BlockTracker{
		seenBlocks: make(map[uint64]bool),
	}
}

// IsFirstAccess checks if the given block address is being accessed for the very first time.
func (bt *BlockTracker) IsFirstAccess(blockAddr uint64) bool {
	bt.mu.Lock()
	defer bt.mu.Unlock()
 
	if bt.seenBlocks[blockAddr] {
		return false
	}
 
	bt.seenBlocks[blockAddr] = true
	return true
}