package fixedchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestAddToChain(t *testing.T) {
	runTestAgainstFile(t, func(file *os.File) {
		// Add block to chain to test
		chain := NewFixedChain(file)

		addedBlock, err := chain.AddBlock([]byte("Some data"))

		if err != nil {
			t.Errorf("Error attempting to add block to chain %v", err)
		}

		if !bytes.Equal(addedBlock.GetData(), []byte("Some data")) {
			t.Errorf("Expected block to contain correct input data but did not")
		}
	})
}

func TestAddToAndReadFromChain(t *testing.T) {
	runTestAgainstFile(t, func(file *os.File) {
		chain := NewFixedChain(file)

		_, err := chain.AddBlock([]byte("Some data"))

		if err != nil {
			t.Errorf("Error attempting to add block to chain %v", err)
		}

		readBlock, err := chain.GetBlockAtIndex(0)

		if err != nil {
			t.Errorf("Failed to read added block %v", err)
			return
		}

		if !bytes.Equal(readBlock.GetData(), []byte("Some data")) {
			t.Errorf("Expected block read back from chain to contain input data but did not")
		}
	})
}

func TestAddMultipleBlocksToChainAndReadFromPosition(t *testing.T) {
	runTestAgainstFile(t, func(file *os.File) {
		chain := NewFixedChain(file)

		chain.AddBlock([]byte("Some data"))
		chain.AddBlock([]byte("Some more data"))
		chain.AddBlock([]byte("Third block"))
		chain.AddBlock([]byte("Fourth block"))
		chain.AddBlock([]byte("Fifth block"))

		fourthBlockRead, _ := chain.GetBlockAtIndex(3)
		if !bytes.Equal(fourthBlockRead.GetData(), []byte("Fourth block")) {
			t.Errorf("Expected the data in the fourth block stored to be 'Fourth block' but was not")
		}

		fifthBlockRead, _ := chain.GetBlockAtIndex(4)
		if !bytes.Equal(fifthBlockRead.GetData(), []byte("Fifth block")) {
			t.Errorf("Expected the data in the fifth block stored to be 'Fifth block' but was not")
		}
	})
}

// Helpers

func runTestAgainstFile(t *testing.T, testRunner func(*os.File)) {
	// Create test file
	file, filePath, fileCreateErr := createFile()
	if fileCreateErr != nil {
		t.Errorf("Error creating test blockchain file %v", fileCreateErr)
	}
	// run test
	testRunner(file)
	// delete file
	deleteFile(filePath)
}

func testHash() []byte {
	hash := sha256.New()
	hash.Write([]byte("This will convert to a sha256 hash"))
	return hash.Sum(nil)
}

func deleteFile(path string) {
	// Cleanup: delete test file
	err := os.Remove(path)
	if err != nil {
		fmt.Errorf("Error cleaning up test file %v", err)
	}
}

func createFile() (*os.File, string, error) {
	fileName, err := exec.Command("uuidgen").Output()
	if err != nil {
		return nil, "", fmt.Errorf("Error generating filename for test blockchain file")
	}

	path := fmt.Sprintf("./%v", fileName)
	file, err2 := os.Create(path)
	return file, path, err2
}
