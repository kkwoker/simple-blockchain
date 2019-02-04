package blockchain

import (
	// "fmt"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

type Blockchain struct {
	Chain []Block
}

func (chain *Blockchain) Add(blk Block) {
	if !blk.ValidHash() {
		panic("adding block with invalid hash")
	}
	chain.Chain = append(chain.Chain, blk)
}

func (chain Blockchain) IsValid() bool {

	validity := true

	// Validate Initial Block
	// The initial block has previous hash all null bytes and is generation zero.
	genesis := chain.Chain[0]
	zeros := make([]byte, 32)
	validity = validity && bytes.Equal(genesis.PrevHash, zeros)
	validity = validity && (genesis.Generation == 0)

	// Validate all other blocks
	for i := 1; i < len(chain.Chain); i++ {
		// Each block has the same difficulty value.
		validity = validity && (genesis.Difficulty == chain.Chain[i].Difficulty)

		// Each block has a generation value that is one more than the previous block.
		validity = validity && (chain.Chain[i].Generation == uint64(i))

		// Each block's previous hash matches the previous block's hash.
		validity = validity && (bytes.Equal(chain.Chain[i].PrevHash, chain.Chain[i-1].Hash))

		// Each block's hash value actually matches its contents.
		prev := hex.EncodeToString(chain.Chain[i].PrevHash)
		gen := strconv.Itoa(int(chain.Chain[i].Generation))
		diff := strconv.Itoa(int(chain.Chain[i].Difficulty))
		proof := strconv.Itoa(int(chain.Chain[i].Proof))
		res := []string{prev, gen, diff, chain.Chain[i].Data, proof}
		str := strings.Join(res, ":")

		// Make hash from values
		h := sha256.New()
		h.Write([]byte(str))
		hash := h.Sum(nil)

		validity = validity && (bytes.Equal(chain.Chain[i].Hash, hash))

		// Each block's hash value ends in difficulty null bytes.
		nulls := make([]byte, chain.Chain[i].Difficulty)
		validity = validity && (bytes.HasSuffix(chain.Chain[i].Hash, nulls))
	}
	return validity
}
