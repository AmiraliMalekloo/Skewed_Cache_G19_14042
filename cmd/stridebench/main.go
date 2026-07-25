// Command stridebench sends a stride-patterned sequence of reads through a
// real writeback cache component (github.com/sarchlab/akita/v5/mem/cache/writeback)
// and prints the cache's own hit/miss statistics counters afterwards.
//
// The stride is fixed to NumSets * BlockSize, so every address the loop
// touches maps to the SAME set under a plain "mod NumSets" addressing
// scheme. On the baseline (non-skewed) cache this should produce almost
// nothing but Conflict Misses and a very low Set Utilization, even though
// the cache as a whole has plenty of free capacity elsewhere.
package main

import (
	"fmt"

	"github.com/sarchlab/akita/v5/mem"
	"github.com/sarchlab/akita/v5/mem/cache/writeback"
	"github.com/sarchlab/akita/v5/mem/idealmemcontroller"
	"github.com/sarchlab/akita/v5/mem/memprotocol"
	"github.com/sarchlab/akita/v5/messaging"
	"github.com/sarchlab/akita/v5/modeling"
	"github.com/sarchlab/akita/v5/noc/directconnection"
	"github.com/sarchlab/akita/v5/timing"
)

func main() {
	engine := timing.NewSerialEngine()

	// --- Backing memory (DRAM) --------------------------------------------
	dramStorage := mem.NewStorage(4 * mem.GB)
	dramSpec := idealmemcontroller.DefaultSpec()
	dramSpec.Width = 1
	dramSpec.Latency = 200
	dramSpec.CacheLineSize = 64

	dram := idealmemcontroller.MakeBuilder().
		WithRegistrar(modeling.NewStandaloneRegistrar(engine)).
		WithResources(idealmemcontroller.Resources{Storage: dramStorage}).
		WithSpec(dramSpec).
		Build("DRAM")
	dram.AssignPort("Top", messaging.NewPort(dram, 256, 256, dram.Name()+".Top"))
	dram.AssignPort("Control", messaging.NewPort(dram, 256, 256, dram.Name()+".Control"))
	addressToPortMapper := &mem.SinglePortMapper{
		Port: dram.GetPortByName("Top").AsRemote(),
	}

	// --- Cache under test ---------------------------------------------------
	cacheSpec := writeback.DefaultSpec()
	cacheSpec.TotalByteSize = 16 * 4 * 64
	cacheSpec.NumReqPerCycle = 1

	cacheComp := writeback.MakeBuilder().
		WithRegistrar(modeling.NewStandaloneRegistrar(engine)).
		WithSpec(cacheSpec).
		WithResources(writeback.Resources{
			AddressToPortMapper: addressToPortMapper,
		}).
		Build("Cache")

	for _, name := range []string{"Top", "Bottom", "Control"} {
		cacheComp.AssignPort(name,
			messaging.NewPort(cacheComp, 256, 256, cacheComp.Name()+"."+name))
	}

	// --- Wiring --------------------------------------------------------------
	agentPort := messaging.NewPort(nil, 256, 256, "Agent.Top")

	conn := directconnection.MakeBuilder().
		WithRegistrar(modeling.NewStandaloneRegistrar(engine)).
		Build("Connection")
	conn.PlugIn(cacheComp.GetPortByName("Top"))
	conn.PlugIn(cacheComp.GetPortByName("Bottom"))
	conn.PlugIn(cacheComp.GetPortByName("Control"))
	conn.PlugIn(dram.GetPortByName("Top"))
	conn.PlugIn(agentPort)

	// --- Derive the real Spec back off the built component -------------------
	blockSize := 1 << cacheSpec.Log2BlockSize
	numSets := int(cacheSpec.TotalByteSize / uint64(cacheSpec.WayAssociativity*blockSize))
	stride := uint64(numSets) * uint64(blockSize)

	fmt.Printf("Cache geometry: NumSets=%d WayAssociativity=%d BlockSize=%d\n",
		numSets, cacheSpec.WayAssociativity, blockSize)
	fmt.Printf("Stride = NumSets * BlockSize = %d bytes\n", stride)

	const traceLength = 200
	const workingSetSize = 5
	var sentIDs []uint64

	// --- Fire the stride-patterned read trace --------------------------------
	for i := 0; i < traceLength; i++ {
		addrIndex := i % workingSetSize
		addr := uint64(addrIndex) * stride

		read := memprotocol.ReadReq{}
		read.ID = timing.GetIDGenerator().Generate()
		read.Src = agentPort.AsRemote()
		read.Dst = cacheComp.GetPortByName("Top").AsRemote()
		read.Address = addr
		read.AccessByteSize = 4
		read.TrafficBytes = 12
		read.TrafficClass = "memprotocol.ReadReq"

		cacheComp.GetPortByName("Top").Deliver(read)
		sentIDs = append(sentIDs, read.ID)
	}

	engine.Run()

	received := 0
	for received < len(sentIDs) {
		rsp := agentPort.RetrieveIncoming()
		if rsp == nil {
			break
		}
		received++
	}
	fmt.Printf("Requests sent: %d, responses received: %d\n", len(sentIDs), received)

	printStats(cacheComp)
}

func printStats(c *writeback.Comp) {
	s := c.State

	var hitRate float64
	if s.TotalAccesses > 0 {
		hitRate = float64(s.TotalHits) / float64(s.TotalAccesses) * 100
	}

	occupiedSets := 0
	for _, count := range s.SetUtilization {
		if count > 0 {
			occupiedSets++
		}
	}
	var setSpreadPct float64
	if len(s.SetUtilization) > 0 {
		setSpreadPct = float64(occupiedSets) / float64(len(s.SetUtilization)) * 100
	}

	fmt.Println("\n=== Baseline Cache Stride-Benchmark Report ===")
	fmt.Printf("Total Accesses     : %d\n", s.TotalAccesses)
	fmt.Printf("Total Hits         : %d (%.2f%% hit rate)\n", s.TotalHits, hitRate)
	fmt.Printf("Compulsory Misses  : %d\n", s.CompulsoryMisses)
	fmt.Printf("Conflict Misses    : %d\n", s.ConflictMisses)
	fmt.Printf("Capacity Misses    : %d\n", s.CapacityMisses)
	fmt.Printf("Sets ever touched  : %d / %d (%.2f%%)\n",
		occupiedSets, len(s.SetUtilization), setSpreadPct)
	fmt.Println("================================================")
}
