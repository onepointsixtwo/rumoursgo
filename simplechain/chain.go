package simplechain

import (
	"bytes"
	"crypto/sha256"
)

// Types

type SimpleBlock struct {
	hash         []byte
	previousHash []byte
	data         []byte
}

type SimpleChain struct {
	blocks []*SimpleBlock
}

// Chain functions

func NewSimpleChain() *SimpleChain {
	return &SimpleChain{make([]*SimpleBlock, 0)}
}

func (chain *SimpleChain) AddDataBlock(data []byte) error {
	chainLength := chain.GetLength()
	if chainLength > 0 {
		lastBlockIndex := chainLength - 1
		lastBlock := chain.blocks[lastBlockIndex]
		chain.addBlock(NewSimpleBlock(data, lastBlock.hash))
	} else {
		chain.addBlock(NewSimpleBlock(data, make([]byte, 0)))
	}

	return nil
}

func (chain *SimpleChain) GetLength() int {
	return len(chain.blocks)
}

func (chain *SimpleChain) addBlock(block *SimpleBlock) {
	chain.blocks = append(chain.blocks, block)
}

func (chain *SimpleChain) IsValid() bool {
	length := chain.GetLength()

	if length < 1 {
		// Can't really have an invalid chain of zero items!
		return true
	}

	var nextBlock *SimpleBlock = nil
	for i := (length - 1); i >= 0; i-- {
		workingBlock := chain.blocks[i]

		// If the next block's previous hash or hash sum is incorrect,
		// our chain is invalid.
		if nextBlock != nil &&
			(!bytes.Equal(nextBlock.previousHash, workingBlock.hash) ||
				!bytes.Equal(nextBlock.hash, hash(nextBlock.data, workingBlock.hash))) {
			return false
		}

		nextBlock = workingBlock
	}

	// Valid unless we find otherwise
	return true
}

// Block functions

func NewSimpleBlock(data []byte, previousHash []byte) *SimpleBlock {
	return &SimpleBlock{hash: hash(data, previousHash), previousHash: previousHash, data: data}
}

// Hash functions
// Could have used the hash interface to inject this into the chain
// but I can't really see how you could change the hashing func? Without
// the whole chain then being invalid?

func hash(data []byte, previousHash []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	if len(previousHash) > 0 {
		hash.Write(previousHash)
	}
	return hash.Sum(nil)
}
