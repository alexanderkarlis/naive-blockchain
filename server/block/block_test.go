package block

import (
	"testing"
)

var (
	genesisBlock = CreateGenesisBlock()
)

type bcTest struct {
	bc    *Blockchain
	valid bool
}

func TestIsValidGoodBlockChain(t *testing.T) {
	var goodBlockchain Blockchain
	goodBlockchain = append(goodBlockchain, *genesisBlock)
	goodBlockchain.AddNewBlockToBlockChain("test data")
	if valid, index := goodBlockchain.IsValidBlockChain(); !valid {
		t.Errorf("previousHash doesn't match hash of previous at %d", index)
	}
}

func TestIsValidBadBlockChain(t *testing.T) {
	var badBlockchain Blockchain
	badBlockchain = append(badBlockchain, *genesisBlock)
	badBlock := &Block{
		ID:           "ID",
		Index:        0,
		Hash:         "aihdashdiuhasidh823h83hd2hndn",
		PreviousHash: "wrong-hash",
		Timestamp:    21234234,
		Data:         "bad block in chain",
	}
	badBlockchain = append(badBlockchain, *badBlock)
	if valid, index := badBlockchain.IsValidBlockChain(); valid {
		t.Errorf("previousHash doesn't match hash of previous at index %d", index)
	}
}
