package fixedchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/onepointsixtwo/rumoursgo/chain"
	"os"
	"time"
)

/*
	Fixed length blocks are specified as follows:

	hash (32 bytes)
	previousHash (32 bytes)
	dataLength (8 bytes)
	data (4008 bytes)
	block number (8 bytes)
	timestamp (8 bytes)

	TOTAL: 4096 bytes (purposely made data 40 bytes short of 4096 so the whole thing fits into a block that is a round number (in binary))
	Should hopefully mean that the blocks line up with either the file read block size or some multiple of it
	Non-data sections: 88 bytes

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
	blockNumber  int64
	timestamp    int64
}

// Chain Methods

func NewFixedChain(file *os.File) *FixedChain {
	return &FixedChain{file}
}

//TODO: get public interface working.
func (chain *FixedChain) AddBlock(data []byte) (chain.Block, error) {
	return nil, nil
}

func (chain *FixedChain) GetSize() uint64 {
	info, err := chain.file.Stat()
	if err != nil {
		return uint64(info.Size() / FixedBlockSize)
	}
	return 0
}

func (chain *FixedChain) GetBlockAtIndex(index uint64) (chain.Block, error) {
	return nil, nil
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

func newFixedSizeBlockFromData(diskData []byte) (*FixedChainBlock, error) {
	if len(diskData) != FixedBlockSize {
		return nil, fmt.Errorf("Expected fixed block length from diskdata of 4096 but was %v", len(diskData))
	}

	hash := diskData[0:32]
	previousHash := diskData[32:64]
	dataLength, err := bytesToLongInteger(diskData[64:72])
	if err != nil {
		return nil, err
	}
	if dataLength > 4008 {
		return nil, fmt.Errorf("Data length too long - should be a max of 4008 but is %v", dataLength)
	}
	data := diskData[72:(72 + dataLength)]
	blockNumber, err2 := bytesToLongInteger(diskData[4080:4088])
	if err2 != nil {
		return nil, err2
	}
	timestamp, err3 := bytesToLongInteger(diskData[4088:4096])
	if err3 != nil {
		return nil, err3
	}

	block := &FixedChainBlock{hash, previousHash, data, blockNumber, timestamp}

	return block, nil
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
	return block.blockNumber
}

func (block *FixedChainBlock) GetCreationTimestamp() int64 {
	return block.timestamp
}

// The amount of bytes you are off being correct on bytes appears to still be one. You need to work out why it's still one short of where it should be - deserialised
// is one byte less than input serialised. Your calculation is incorrect.
func (block *FixedChainBlock) getDiskWriteData() []byte {
	output := make([]byte, 4096)
	writeBytesToSlice(output, block.hash, 0, 32)
	writeBytesToSlice(output, block.previousHash, 32, 32)
	writeBytesToSlice(output, longIntegerToBytes(int64(len(block.data))), 64, 8)
	writeBytesToSlice(output, block.data, 72, 4008)
	writeBytesToSlice(output, longIntegerToBytes(block.blockNumber), 4080, 8)
	writeBytesToSlice(output, longIntegerToBytes(block.timestamp), 4088, 8)
	return output
}

// Helpers

func writeBytesToSlice(slice []byte, data []byte, writeOffset int, maxWriteLength int) {
	for i := 0; i < min(len(data), maxWriteLength); i++ {
		slice[writeOffset+i] = data[i]
	}
}

func min(one int, two int) int {
	if one < two {
		return one
	}
	return two
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

func bytesToLongInteger(bytesRepresentation []byte) (int64, error) {
	buffer := bytes.NewBuffer(bytesRepresentation)
	return binary.ReadVarint(buffer)
}
