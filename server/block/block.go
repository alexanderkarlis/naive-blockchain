package block

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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
	Difficulty   int
	Nonce        int
}

// Blockchain is a slice of Blocks.
type Blockchain []Block

// BlockchainToBytes turns blockchain into a slice of bytes to broadcast across sockets.
func (bc *Blockchain) BlockchainToBytes() []byte {
	bSlice, err := json.Marshal(bc)
	if err != nil {
		log.Printf("Could not send blockchain to bytes. %v", err)
		return []byte{}
	}
	return bSlice
}

// AddNewBlockToBlockChain returns a new instance of a single Block in the Blockchain.
func (bc *Blockchain) AddNewBlockToBlockChain(d string) Blockchain {
	id := generateUUID()
	index := len(*bc)
	previousHash := (*bc)[len(*bc)-1].Hash
	timeStamp := generateTimeStamp()
	data := d
	difficulty := 2
	hash := findBlockHash(index, previousHash, timeStamp, data, difficulty)
	newBlock := &Block{
		ID:           id,
		Index:        index,
		Timestamp:    timeStamp,
		Data:         data,
		Hash:         hash,
		PreviousHash: previousHash,
		Difficulty:   difficulty,
		Nonce:        0,
	}
	*bc = append(*bc, *newBlock)

	return *bc
}

func findBlockHash(index int, previousHash string, ts int, data string, difficulty int) string {
	var hash string
	nonce := 0
	for {
		hash = hash256(index, previousHash, ts, data, difficulty, nonce)
		if hashMatchesDifficulty(hash, difficulty) {
			return hash
		}
		nonce++
	}
}

func hashMatchesDifficulty(hash string, difficulty int) bool {
	hashBinary, err := hex.DecodeString(hash)
	if err != nil {
		fmt.Println("could not digest hash into binary")
	}
	requiredPrefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(string(hashBinary), requiredPrefix)
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
	hash := "3e7f6189598ac7cdd024f5d07b1fb1fd18e1ce68a3986d44b36b97c2ecded8e0"

	return &Block{
		ID:           generateUUID(),
		Index:        0,
		Timestamp:    timeStamp,
		Data:         data,
		Hash:         hash,
		PreviousHash: "",
	}
}

func generateTimeStamp() int {
	return int(time.Now().Unix())
}

func generateUUID() string {
	uid := uuid.New()
	return uid.String()
}

// 256 Hashing func
func hash256(index int, previousHash string, timestamp int, data string, difficulty int, nonce int) string {
	hashString := fmt.Sprintf("%d%s%d%s%d%d", index, previousHash, timestamp, data, difficulty, nonce)
	sum := sha256.Sum256([]byte(hashString))
	return fmt.Sprintf("%+x", sum[:])
}
