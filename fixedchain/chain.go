package fixedchain

import (
	"fmt"
	"github.com/onepointsixtwo/rumoursgo/chain"
	"github.com/onepointsixtwo/rumoursgo/chainutil"
	"os"
)

// Types

type FixedChain struct {
	file *os.File
}

// Chain Methods

func NewFixedChain(file *os.File) *FixedChain {
	return &FixedChain{file}
}

func (chain *FixedChain) AddBlock(data []byte) (chain.Block, error) {
	dataLen := len(data)
	if dataLen > MaxDataSize {
		return nil, fmt.Errorf("Max data size should be %v but was %v", MaxDataSize, dataLen)
	}

	latestBlock, err := chain.getLatestBlock()
	if err != nil {
		return nil, err
	}

	var previousBlockNumber int64 = -1
	previousHash := make([]byte, 0)
	if latestBlock != nil {
		previousBlockNumber = latestBlock.blockNumber
		previousHash = latestBlock.hash
	}

	newBlock := NewFixedSizeBlock(data, previousBlockNumber, previousHash)
	return chain.appendBlock(newBlock)
}

func (chain *FixedChain) GetSize() uint64 {
	info, err := chain.file.Stat()
	if err != nil {
		return uint64(info.Size() / FixedBlockSize)
	}
	return 0
}

func (chain *FixedChain) GetBlockAtIndex(index uint64) (chain.Block, error) {
	return chain.getBlockAtIndex(index)
}

func (chain *FixedChain) IsValid() bool {
	return chainutil.ChainIsValid(chain)
}

func (chain *FixedChain) getLatestBlock() (*FixedChainBlock, error) {
	size := chain.GetSize()
	if size == 0 {
		return nil, nil
	}
	return chain.getBlockAtIndex(size - 1)
}

func (chain *FixedChain) appendBlock(block *FixedChainBlock) (chain.Block, error) {
	info, statErr := chain.file.Stat()
	if statErr != nil {
		return nil, fmt.Errorf("Error attempting to get stats for blockchain file %v", statErr)
	}

	fileLength := info.Size()
	_, err := chain.file.WriteAt(block.GetDiskWriteData(), fileLength)
	if err != nil {
		return nil, fmt.Errorf("Error writing block to disk %v", err)
	}

	return block, nil
}

func (chain *FixedChain) getBlockAtIndex(index uint64) (*FixedChainBlock, error) {
	byteOffset := int64(index * FixedBlockSize)
	bytes := make([]byte, FixedBlockSize)

	bytesRead, err := chain.file.ReadAt(bytes, byteOffset)
	if err != nil {
		return nil, fmt.Errorf("Error attempting to read block from file %v", err)
	}
	if bytesRead != FixedBlockSize {
		return nil, fmt.Errorf("Incorrect number of bytes read attempting to read block from file - got %v", bytesRead)
	}

	return NewFixedSizeBlockFromData(bytes)
}
