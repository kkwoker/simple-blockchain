package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockCalcHash(t *testing.T) {
	b0 := Initial(2)

	// SetProof calls CalcHash and sets the Hash
	b0.SetProof(242278)
	str0 := hex.EncodeToString(b0.Hash)
	assert.Equal(t, "29528aaf90e167b2dc248587718caab237a81fd25619a5b18be4986f75f30000", str0)

	b1 := b0.Next("message")
	b1.SetProof(75729)
	str1 := hex.EncodeToString(b1.Hash)
	assert.Equal(t, "02b09bde9ff60582ef21baa4bef87a95dfcd67efaf258e6df60463da0a940000", str1)
}

func TestBlockValidHashFalse(t *testing.T) {
	b0 := Initial(2)
	assert.Equal(t, false, b0.ValidHash())
}

func TestBlockValidHashTrue(t *testing.T) {
	b0 := Initial(2)
	b0.SetProof(242278)
	assert.Equal(t, true, b0.ValidHash())
}

func TestBlockMine1(t *testing.T) {
	b0 := Initial(2)
	b0.Mine(1)
	fmt.Println("Proof: ", b0.Proof, hex.EncodeToString(b0.Hash))
	fmt.Println(b0.Hash)
	assert.Equal(t, uint64(242278), b0.Proof)

	b1 := b0.Next("this is an interesting message")
	b1.Mine(1)
	fmt.Println(b1.Proof, hex.EncodeToString(b1.Hash))
	assert.Equal(t, uint64(41401), b1.Proof)

	b2 := b1.Next("this is not interesting")
	b2.Mine(1)
	fmt.Println(b2.Proof, hex.EncodeToString(b2.Hash))
	assert.Equal(t, uint64(195955), b2.Proof)
}

func TestBlockMine2(t *testing.T) {
	b0 := Initial(3)
	b0.Mine(5)
	fmt.Println(b0.Proof, hex.EncodeToString(b0.Hash))
	assert.Equal(t, uint64(8816998), b0.Proof)

	b1 := b0.Next("this is an interesting message")
	b1.Mine(5)
	fmt.Println(b1.Proof, hex.EncodeToString(b1.Hash))
	assert.Equal(t, uint64(16634616), b1.Proof)

	b2 := b1.Next("this is not interesting")
	b2.Mine(5)
	fmt.Println(b2.Proof, hex.EncodeToString(b2.Hash))
	assert.Equal(t, uint64(8148543), b2.Proof)
}

func TestBlockChain(t *testing.T) {
	b0 := Initial(2)
	b0.Mine(1)
	b1 := b0.Next("this is an interesting message")
	b1.Mine(1)
	b2 := b1.Next("this is not interesting")
	b2.Mine(1)

	bc := new(Blockchain)
	bc.Add(b0)
	bc.Add(b1)
	bc.Add(b2)

	fmt.Println(bc)
	assert.Equal(t, true, bc.IsValid())
}
