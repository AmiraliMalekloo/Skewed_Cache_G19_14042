package cache

import (
	"time"
)

// Lookup performs a lookup operation across all ways in the skewed cache.
func (c *SkewedCache) Lookup(addr uint64) (hit bool, matchedWay int, matchedSet int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	blockAddr := addr / uint64(c.BlockSize)
	tag := blockAddr

	c.TotalAccesses++

	for wayID := 0; wayID < c.Associativity; wayID++ {
		setIdx := c.MapToSet(blockAddr, wayID)
		line := &c.Sets[setIdx][wayID]

		if line.Valid && line.Tag == tag {
			line.LastAccessTime = time.Now().UnixNano()
			c.TotalHits++
			return true, wayID, setIdx
		}
	}

	c.handleMiss(blockAddr)
	return false, -1, -1
}

// handleMiss categorizes and updates statistics for different types of cache misses.
func (c *SkewedCache) handleMiss(blockAddr uint64) {
	if c.Tracker.IsFirstAccess(blockAddr) {
		c.CompulsoryMisses++
		return
	}

	totalCapacity := uint64(c.NumSets * c.Associativity)
	if c.OccupiedBlocks >= totalCapacity {
		// The whole cache is full: this block had to lose to some other
		// live block no matter which way/set combination it tried.
		c.CapacityMisses++
		return
	}

	// The cache as a whole still has free lines somewhere, but every one
	// of this block's W candidate slots (across its W distinct sets) was
	// already occupied by something else -> a conflict miss.
	c.ConflictMisses++
}

// Replace selects a victim block among the W candidates (across different sets) using LRU and replaces it.
func (c *SkewedCache) Replace(addr uint64, data []byte) (victimWay int, victimSet int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	blockAddr := addr / uint64(c.BlockSize)

	type candidateSlot struct {
		wayID int
		setID int
	}

	candidates := make([]candidateSlot, c.Associativity)
	for wayID := 0; wayID < c.Associativity; wayID++ {
		candidates[wayID] = candidateSlot{
			wayID: wayID,
			setID: c.MapToSet(blockAddr, wayID),
		}
	}

	chosenWay, chosenSet := -1, -1

	// First, look for any free (invalid) slot among the W candidates -
	// no eviction needed in that case.
	for _, cand := range candidates {
		line := &c.Sets[cand.setID][cand.wayID]
		if !line.Valid {
			chosenWay, chosenSet = cand.wayID, cand.setID
			c.OccupiedBlocks++
			break
		}
	}

	// If every candidate slot is occupied, fall back to LRU: pick the
	// candidate with the smallest (oldest) LastAccessTime, even though
	// the candidates live in different sets.
	if chosenWay == -1 {
		oldestTime := int64(1<<63 - 1) // max int64
		for _, cand := range candidates {
			line := &c.Sets[cand.setID][cand.wayID]
			if line.LastAccessTime < oldestTime {
				oldestTime = line.LastAccessTime
				chosenWay, chosenSet = cand.wayID, cand.setID
			}
		}
	}

	c.Sets[chosenSet][chosenWay] = CacheLine{
		Tag:            blockAddr,
		Valid:          true,
		Dirty:          false,
		Data:           data,
		LastAccessTime: time.Now().UnixNano(),
	}

	return chosenWay, chosenSet
}

// GetSetUtilization calculates the current utilization rate of the cache.
func (c *SkewedCache) GetSetUtilization() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	totalCapacity := float64(c.NumSets * c.Associativity)
	if totalCapacity == 0 {
		return 0.0
	}

	return float64(c.OccupiedBlocks) / totalCapacity
}
