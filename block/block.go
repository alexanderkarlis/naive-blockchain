package block

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"time"

	"os/exec"
)

// Blocks function interface
type Blocks interface {
	New(data string)
}

// Block struct
type Block struct {
	ID           string
	Index        int
	Hash         string
	PreviousHash string
	Timestamp    int
	Data         string
}

// Blockchain is a slice of Blocks.
type Blockchain []Block

// AddNewBlockToBlockChain returns a new instance of a single Block in the Blockchain.
func (bc *Blockchain) AddNewBlockToBlockChain(d string) Blockchain {
	id := generateUUID()
	index := 0
	previousHash := (*bc)[len(*bc)-1].Hash
	timeStamp := generateTimeStamp()
	data := d
	hash := hash256(index, previousHash, timeStamp, data)

	newBlock := &Block{
		ID:           id,
		Index:        index,
		Timestamp:    timeStamp,
		Data:         data,
		Hash:         hash,
		PreviousHash: previousHash,
	}
	*bc = append(*bc, *newBlock)

	return *bc
}

// IsValidBlockChain validates the active blockchain by checking previous hash
// against current 'previous hash'.
// Returns valid (bool) and index of problem block in chain. (-1 if blockchain is ok.)
func (bc *Blockchain) IsValidBlockChain() (valid bool, i int) {
	valid = true
	for i := 1; i < len(*bc); i++ {
		if (*bc)[i].PreviousHash != (*bc)[i-1].Hash {
			valid = false
			return valid, i
		}
	}
	return valid, -1
}

// CreateGenesisBlock creates the first block in the blockchain. Has index of 0.
func CreateGenesisBlock() *Block {
	timeStamp := generateTimeStamp()
	data := "GENESIS BLOCK"
	hash := hash256(0, "", timeStamp, data)

	return &Block{
		ID:           generateUUID(),
		Index:        0,
		Timestamp:    timeStamp,
		Data:         data,
		Hash:         hash,
		PreviousHash: "",
	}
}

// CheckPreviousHash should be called as soon as any data on the block changes.
func (b *Block) CheckPreviousHash() *Block {
	newHash := hash256(b.Index, b.PreviousHash, b.Timestamp, b.Data)
	b.Hash = string(newHash)
	b.PreviousHash = b.Hash
	return b
}

func generateTimeStamp() int {
	return int(time.Now().Unix())
}

func generateUUID() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	uuid := strings.TrimRight(string(out), "\n")
	return uuid
}

// 256 Hashing func
func hash256(index int, previousHash string, timestamp int, data string) string {
	hashString := fmt.Sprintf("%d%s%d%s", index, previousHash, timestamp, data)
	sum := sha256.Sum256([]byte(hashString))
	return fmt.Sprintf("%+x", sum[:])
}
