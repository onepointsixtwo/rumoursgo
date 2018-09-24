package fixedchain

import (
	"bytes"
	"testing"
)

func TestFixedChainBlockSerialisationAndDeserialisation(t *testing.T) {
	hash := testHash()
	prevHash := testHash()
	block := &FixedChainBlock{hash, prevHash, []byte("Some arbitrary data"), 3, 500}

	serialised := block.GetDiskWriteData()

	serialisedLength := len(serialised)
	if serialisedLength != FixedBlockSize {
		t.Errorf("Expected serialised length to be %v but was %v", FixedBlockSize, serialisedLength)
	}

	deserialised, err := NewFixedSizeBlockFromData(serialised)

	if err != nil {
		t.Errorf("Unable to deserialise fixed size block from its own data %v", err)
	}

	if deserialised.blockNumber != 3 {
		t.Errorf("Deserialised block number should have been 3 but was %v", deserialised.blockNumber)
	}

	dataString := string(deserialised.data)
	if dataString != "Some arbitrary data" {
		t.Errorf("Expected deserialised data to be same as input string but was '%v'", dataString)
	}

	if !bytes.Equal(hash, deserialised.hash) {
		t.Errorf("Expected hash %v to be equal to deserialised hash %v", hash, deserialised.hash)
	}

	if !bytes.Equal(prevHash, deserialised.previousHash) {
		t.Errorf("Expected previous hash %v to be equal to deserialised previous hash %v", prevHash, deserialised.previousHash)
	}

	if deserialised.timestamp != 500 {
		t.Errorf("Expected deserialised timestamp to be 500 but was %v", deserialised.timestamp)
	}
}
