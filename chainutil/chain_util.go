package chainutil

import (
	"bytes"
	"github.com/onepointsixtwo/rumoursgo/chain"
)

func ChainIsValid(ch chain.Chain) bool {
	length := ch.GetSize()

	if length < 1 {
		// Can't really have an invalid chain of zero items!
		return true
	}

	var nextBlock chain.Block = nil
	for i := (length - 1); i >= 0; i-- {
		workingBlock, err := ch.GetBlockAtIndex(i)
		if err != nil {
			// If we have an error reading a block we can't verify this chain is valid
			// so return false
			return false
		}

		// If the next block's previous hash or hash sum is incorrect,
		// our chain is invalid.
		if nextBlock != nil &&
			(!bytes.Equal(nextBlock.GetPreviousHash(), workingBlock.GetHash()) ||
				!workingBlock.VerifyBlockHash()) {
			return false
		}

		nextBlock = workingBlock
	}

	// Valid unless we find otherwise
	return true
}
