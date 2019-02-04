package blockchain

import (
	"work_queue"
)

type miningWorker struct {
	Block Block
	Start uint64
	End   uint64
}

type MiningResult struct {
	Proof uint64 // proof-of-work value, if found.
	Found bool   // true if valid proof-of-work was found.
}

func (m miningWorker) Run() interface{} {
	found := false
	for prf := m.Start; prf <= m.End; prf++ {
		m.Block.SetProof(uint64(prf))
		found = m.Block.ValidHash()
		if found {
			return MiningResult{Proof: uint64(prf), Found: true}
		}
	}
	return MiningResult{Found: false}
}

// Give the worker the hash so it can validate
func newWorker(blk Block, start uint64, end uint64) miningWorker {
	w := new(miningWorker)
	w.Block = blk
	w.Start = start
	w.End = end
	return *w
}

// Mine the range of proof values, by breaking up into chunks and checking
// "workers" chunks concurrently in a work queue. Should return shortly after a result
// is found.
// eg. from 0 to 16, split it into chunks. then give each chunk to the worker to see if any of them are valid
func (blk Block) MineRange(start uint64, end uint64, workers uint64, chunks uint64) MiningResult {
	queue := work_queue.Create(uint(workers), uint(2*chunks))
	chunkSize := (end - start) / chunks

	for i := uint64(0); i < end; i += chunkSize * workers {
		// Build n workers
		for w := i; w < i+chunkSize*workers; w += chunkSize {
			worker := newWorker(blk, w, w+chunkSize)
			queue.Enqueue(worker)
		}
	}

	for mining_result := range queue.Results {
		mr, _ := mining_result.(MiningResult)
		if mr.Found == true {
			queue.Shutdown()
			return mr
		}
	}
	return MiningResult{}
}

// Call .MineRange with some reasonable values that will probably find a result.
// Good enough for testing at least. Updates the block's .Proof and .Hash if successful.
func (blk *Block) Mine(workers uint64) bool {
	reasonableRangeEnd := uint64(4 * 1 << (8 * blk.Difficulty)) // 4 * 2^(bits that must be zero)
	mr := blk.MineRange(0, reasonableRangeEnd, workers, 4321)
	if mr.Found {
		blk.SetProof(mr.Proof)
	}
	return mr.Found
}
