package cache

import (
	"fmt"
	"time"
)

// Lookup performs a lookup operation across all ways in the skewed cache.
func (c *SkewedCache) Lookup(addr uint64) (hit bool, matchedWay int, matchedSet int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	blockAddr := addr / uint64(c.BlockSize)
	tag := blockAddr

	c.TotalAccesses++

	// TODO: Implement the lookup loop across all ways.
	// Key difference in Skewed Cache:
	// For each wayID from 0 to Associativity-1:
	//   1. Calculate setIdx = c.MapToSet(blockAddr, wayID)
	//   2. Inspect the line c.Sets[setIdx][wayID]
	//   3. If Valid == true and Tag == tag:
	//        - Hit detected!
	//        - Update last access timestamp (LastAccessTime).
	//        - c.TotalHits++
	//        - return true, wayID, setIdx

	// If block is not found after inspecting all ways -> Miss occurs.
	c.handleMiss(blockAddr)
	return false, -1, -1
}

// handleMiss categorizes and updates statistics for different types of cache misses.
func (c *SkewedCache) handleMiss(blockAddr uint64) {
	// TODO: Implement Miss Classification logic:
	//
	// 1. Check for Compulsory Miss:
	//    if c.Tracker.IsFirstAccess(blockAddr) {
	//        c.CompulsoryMisses++
	//        return
	//    }
	//
	// 2. If not Compulsory, classify as Capacity Miss vs Conflict Miss:
	//    - Is the entire cache capacity full? (c.OccupiedBlocks == c.NumSets * c.Associativity)
	//      - If YES -> c.CapacityMisses++
	//      - If NO  -> c.ConflictMisses++ 
	//        (Because space exists in the cache, but all W candidate slots for this block are occupied)
}

// Replace selects a victim block among the W candidates (across different sets) using LRU and replaces it.
func (c *SkewedCache) Replace(addr uint64, data []byte) (victimWay int, victimSet int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	blockAddr := addr / uint64(c.BlockSize)

	// TODO: Implement Victim Selection logic for Skewed Cache:
	//
	// 1. Build a candidate array of size W. For each wayID from 0 to Associativity-1:
	//    setIdx = c.MapToSet(blockAddr, wayID)
	//    Candidate i: c.Sets[setIdx][wayID]
	//
	// 2. Check if an invalid/empty block (Valid == false) exists among the W candidates:
	//    - If found, select it directly (no valid block needs eviction).
	//    - c.OccupiedBlocks++
	//
	// 3. If all candidates are valid, apply the LRU policy:
	//    - Identify the candidate with the oldest LastAccessTime among these W distinct sets.
	//    - Select that entry as the victim.
	//
	// 4. Replace with the new block:
	//    - c.Sets[chosenSet][chosenWay] = CacheLine{ Tag: blockAddr, Valid: true, LastAccessTime: time.Now().UnixNano(), Data: data }

	return -1, -1
}

// GetSetUtilization calculates the current utilization rate of the cache.
func (c *SkewedCache) GetSetUtilization() float64 {
	// TODO: Calculate the ratio of occupied blocks relative to total cache capacity,
	// or the percentage of sets containing at least one valid line.
	return 0.0
}