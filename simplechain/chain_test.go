package simplechain

import (
	"bytes"
	"testing"
)

// Tests

func TestAddingDataBlocks(t *testing.T) {
	chain := NewSimpleChain()

	firstData := stringToBytes("ONE")
	chain.AddDataBlock(firstData)

	if chain.GetLength() != 1 {
		t.Errorf("Expected chain length after adding one data block to be 1 but was %v", chain.GetLength())
	}

	addedBlock := chain.blocks[0]
	if !bytes.Equal(addedBlock.data, firstData) {
		t.Errorf("Expected added bytes to be equal to bytes of first block")
	}
	if !bytes.Equal(hash(addedBlock.data, make([]byte, 0)), addedBlock.hash) {
		t.Errorf("Expected first block's hash to be equal to hash of first block data")
	}

	secondData := stringToBytes("TWO")
	chain.AddDataBlock(secondData)

	if chain.GetLength() != 2 {
		t.Errorf("Expected chain length after adding two data blocks to be 2 but was %v", chain.GetLength())
	}

	secondAddedBlock := chain.blocks[1]
	if !bytes.Equal(secondAddedBlock.data, secondData) {
		t.Errorf("Expected second block's data to be equal to second data input")
	}
	if !bytes.Equal(hash(secondData, addedBlock.hash), secondAddedBlock.hash) {
		t.Errorf("Unexpected calculation for second block's hash - should have been amalgam of second block value and first block hash but was not")
	}
}

func TestChainLength(t *testing.T) {
	chain := NewSimpleChain()

	for i := 0; i < 30; i++ {
		chain.AddDataBlock(stringToBytes("I"))
	}

	if len(chain.blocks) != chain.GetLength() {
		t.Errorf("Somehow got chain length function wrong!")
	}
}

func TestNonCorruptedChainIsValid(t *testing.T) {
	chain := NewSimpleChain()
	chain.AddDataBlock(stringToBytes("Testing"))
	chain.AddDataBlock(stringToBytes("That"))
	chain.AddDataBlock(stringToBytes("My"))
	chain.AddDataBlock(stringToBytes("Chain"))
	chain.AddDataBlock(stringToBytes("Is"))
	chain.AddDataBlock(stringToBytes("Valid"))

	if !chain.IsValid() {
		t.Errorf("Untamptered chain should be perfectly valid")
	}
}

func TestChainIsNotValidWhenCorrupted(t *testing.T) {
	// In this test I manually add an invalid block. Obviously the public interface of the
	// chain struct's methods doesn't actually allow this, but eventually there will be serialise / deserialise methods for
	// writing the chain to a file via a Reader / Writer and that's when we need to test the validity of our chain
	// -- when it could have been tampered with by another party.

	chain := NewSimpleChain()
	chain.AddDataBlock(stringToBytes("Testing"))
	chain.AddDataBlock(stringToBytes("That"))
	chain.AddDataBlock(stringToBytes("This"))
	chain.AddDataBlock(stringToBytes("Chain"))
	chain.AddDataBlock(stringToBytes("Is"))
	chain.AddDataBlock(stringToBytes("Invalid"))

	block := NewSimpleBlock(stringToBytes("Some data we're trying to insert"), stringToBytes("Some invalid hash"))
	chain.addBlock(block)

	chain.AddDataBlock(stringToBytes("After"))
	chain.AddDataBlock(stringToBytes("Rubbish"))
	chain.AddDataBlock(stringToBytes("Block"))
	chain.AddDataBlock(stringToBytes("Added"))

	if chain.IsValid() {
		t.Errorf("Chain should be invalid after entering block with invalid previous hash")
	}
}

// Helpers

func stringToBytes(value string) []byte {
	return []byte(value)
}

func bytesToString(value []byte) string {
	return string(value)
}
