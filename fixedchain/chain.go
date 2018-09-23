package fixedchain

import (
	"crypto/sha256"
	"encoding/binary"
	"os"
	"time"
)

/*
	Fixed length blocks are specified as follows:

	hash (8 bytes)
	previousHash (8 bytes)
	data (4064 bytes)
	block number (8 bytes)
	timestamp (8 bytes)

	TOTAL: 4096 bytes (purposely made data 32 bytes short of 4096 so the whole thing fits into a block that is a round number (in binary))
	Should hopefully mean that the blocks line up with either the file read block size or some multiple of it

*/

const (
	FixedBlockSize = 4096
)

// Types

type FixedChain struct {
	file *os.File
}

type FixedChainBlock struct {
	hash         []byte
	previousHash []byte
	data         []byte
	blockNum     int64
	timestamp    int64
}

// Chain Methods
func NewFixedChain(file *os.File) *FixedChain {
	return &FixedChain{file}
}

//TODO: get public interface working.
func (chain *FixedChain) AddBlock(data []byte) (Block, error) {

}

func (chain *FixedChain) GetSize() uint64 {
	info, err := file.Stat()
	if err != nil {
		return info.Size() / FixedBlockSize
	}
	return 0
}

func (chain *FixedChain) GetBlockAtIndex(index uint64) Block {

}

func (chain *FixedChain) IsValid() bool {
	// Function should be moved out to a general function which checks if any Chain (interface) is valid.
	// Doesn't need to know internals at all.
	return true
}

// Block Methods
func NewFixedSizeBlock(data []byte, previousBlockNumber int64, previousHash []byte) *FixedChainBlock {
	blockNumber := previousBlockNumber + 1
	timestamp := time.Now().UnixNano()
	newBlockHash := hash(data, blockNumber, timestamp, previousHash)
	return &FixedChainBlock{newBlockHash, previousHash, data, blockNumber, timestamp}
}

func (block *FixedChainBlock) GetHash() []byte {
	return block.hash
}

func (block *FixedChainBlock) GetPrevious() []byte {
	return block.previousHash
}

func (block *FixedChainBlock) GetData() []byte {
	return block.data
}

func (block *FixedChainBlock) GetBlockNumber() int64 {
	return block.blockNum
}

func (block *FixedChainBlock) GetCreationTimestamp() int64 {
	return block.timestamp
}

// Using the same hash as simplechain
func hash(data []byte, blockNumber int64, timestamp int64, previousHash []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	hash.Write(longIntegerToBytes(blockNumber))
	hash.Write(longIntegerToBytes(timestamp))
	if len(previousHash) > 0 {
		hash.Write(previousHash)
	}
	return hash.Sum(nil)
}

func longIntegerToBytes(longInteger int64) []byte {
	buffer := make([]byte, binary.MaxVarintLen64)
	lengthWritten := binary.PutVarint(buffer, longInteger)
	return buffer[:lengthWritten]
}
