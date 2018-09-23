package chain

type Chain interface {
	AddBlock(data []byte) (Block, error)
	GetSize() uint64
	GetBlockAtIndex(index uint64) Block
	IsValid() bool
}

type Block interface {
	GetHash() []byte
	GetPreviousHash() []byte
	GetData() []byte
	GetBlockNumber() int64
	GetCreationTimestamp() int64
}
