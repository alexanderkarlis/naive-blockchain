package block

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"os/exec"
)

// Block struct
type Block struct {
	ID           string
	Index        int
	Hash         string
	PreviousHash string
	Timestamp    int
	Data         string
}

// New returns a new instance of a single Block in the Blockchain.
func New(d string) *Block {
	id := generateUUID()
	index := 0
	previousHash := ""
	timeStamp := generateTimeStamp()
	data := d
	hash := hash256(index, previousHash, timeStamp, data)

	return &Block{
		ID:           id, // generic ID, TODO: creeate algorithm
		Index:        index,
		Timestamp:    timeStamp,
		Data:         data,
		Hash:         hash,
		PreviousHash: previousHash,
	}
}

// CreateGenesisBlock creates the first block in the blockchain. Has index of 0.
func CreateGenesisBlock() *Block {
	timeStamp := generateTimeStamp()
	data := ""
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

// SetPreviousHash should be called as soon as any data on the block changes.
func (b *Block) SetPreviousHash() *Block {
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
	return string(out)
}

// 256 Hashing func
func hash256(index int, previousHash string, timestamp int, data string) string {
	hashString := fmt.Sprintf("%d%s%d%s", index, previousHash, timestamp, data)
	sum := sha256.Sum256([]byte(hashString))
	return fmt.Sprintf("%+x", sum[:])
}
