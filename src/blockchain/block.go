package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

type Block struct {
	Generation uint64
	Difficulty uint8
	Data       string
	PrevHash   []byte
	Hash       []byte
	Proof      uint64
}

// Create new initial (generation 0) block.
func Initial(difficulty uint8) Block {
	blk := Block{
		Generation: 0,
		Difficulty: difficulty,
		Data:       "",
		PrevHash:   make([]byte, 32),
	}
	return blk
}

// Create new block to follow this block, with provided data.
func (prev_block Block) Next(data string) Block {
	blk := Block{
		Generation: prev_block.Generation + 1,
		Difficulty: prev_block.Difficulty,
		Data:       data,
		PrevHash:   prev_block.Hash,
	}
	return blk
}

// Calculate the block's hash.
func (blk Block) CalcHash() []byte {
	// blk.PrevHash : blk.Generation : blk.Difficulty : blk.Data : blk.Proof
	prev := hex.EncodeToString(blk.PrevHash)
	gen := strconv.Itoa(int(blk.Generation))
	diff := strconv.Itoa(int(blk.Difficulty))
	proof := strconv.Itoa(int(blk.Proof))

	res := []string{prev, gen, diff, blk.Data, proof}
	str := strings.Join(res, ":")

	h := sha256.New()
	h.Write([]byte(str))

	return h.Sum(nil)

}

// Is this block's hash valid?
func (blk Block) ValidHash() bool {
	// returns true if (and only if) the block's hash ends in .Difficulty zero bytes
	// if it had the proof value given as an argument.
	zeros := make([]byte, blk.Difficulty)
	bool := bytes.HasSuffix(blk.Hash, zeros)
	return bool
}

// Set the proof-of-work and calculate the block's "true" hash.
func (blk *Block) SetProof(proof uint64) {
	blk.Proof = proof
	blk.Hash = blk.CalcHash()
}
