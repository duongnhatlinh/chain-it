package mem

import (
	"sync"

	"github.com/gogo/protobuf/sortkeys"
	"it-chain/blockchain"
)

type BlockPool struct {
	blockMap blockchain.BlockMap
	mux      sync.Mutex
}

func NewBlockPool() *BlockPool {
	return &BlockPool{
		blockMap: make(map[uint64]blockchain.DefaultBlock),
		mux:      sync.Mutex{},
	}
}

func (b *BlockPool) Add(block blockchain.DefaultBlock) {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.blockMap[block.GetHeight()] = block
}

func (b *BlockPool) Delete(height uint64) {
	b.mux.Lock()
	defer b.mux.Unlock()

	delete(b.blockMap, height)
}

func (b *BlockPool) GetByHeight(height uint64) blockchain.DefaultBlock {
	b.mux.Lock()
	defer b.mux.Unlock()

	if block, ok := b.blockMap[height]; ok {
		return block
	}

	return blockchain.DefaultBlock{}
}

func (b *BlockPool) GetSortedKeys() []blockchain.BlockHeight {
	keys := make([]blockchain.BlockHeight, 0)
	for h := range b.blockMap {
		keys = append(keys, h)
	}

	sortkeys.Uint64s(keys)

	return keys
}

func (b *BlockPool) Size() int {
	b.mux.Lock()
	defer b.mux.Unlock()

	return len(b.blockMap)
}
