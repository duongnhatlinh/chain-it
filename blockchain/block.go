package blockchain

import (
	"bytes"
	"encoding/json"
	"reflect"
	"time"
	ygg "github.com/DE-labtory/yggdrasill/common"
)

type Block = ygg.Block

type BlockHeight = uint64

type BlockState = string

const (
	Created   BlockState = "Created"
	Staged    BlockState = "Staged"
	Committed BlockState = "Committed"
)

type DefaultBlock struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []*DefaultTransaction
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   string
	State     BlockState
}

func (block *DefaultBlock) SetSeal(seal []byte) {
	block.Seal = seal
}

func (block *DefaultBlock) SetPrevSeal(prevSeal []byte) {
	block.PrevSeal = prevSeal
}

func (block *DefaultBlock) SetHeight(height uint64) {
	block.Height = height
}

func (block *DefaultBlock) PutTx(transaction Transaction) error {
	convTx, ok := transaction.(*DefaultTransaction)
	if ok {
		if block.TxList == nil {
			block.TxList = make([]*DefaultTransaction, 0)
		}
		block.TxList = append(block.TxList, convTx)

		return nil
	}

	return ErrTransactionType
}

func (block *DefaultBlock) SetTxSeal(txSeal [][]byte) {
	block.TxSeal = txSeal
}

func (block *DefaultBlock) SetCreator(creator string) {
	block.Creator = creator
}

func (block *DefaultBlock) SetTimestamp(currentTime time.Time) {
	block.Timestamp = currentTime
}

func (block *DefaultBlock) SetState(state BlockState) {
	block.State = state
}

func (block *DefaultBlock) GetSeal() []byte {
	return block.Seal
}

func (block *DefaultBlock) GetPrevSeal() []byte {
	return block.PrevSeal
}

func (block *DefaultBlock) GetHeight() uint64 {
	return block.Height
}

func (block *DefaultBlock) GetTxList() []Transaction {
	txList := make([]Transaction, 0)
	for _, tx := range block.TxList {
		txList = append(txList, tx)
	}
	return txList
}

func (block *DefaultBlock) GetTxSeal() [][]byte {
	return block.TxSeal
}

func (block *DefaultBlock) GetCreator() string {
	return block.Creator
}

func (block *DefaultBlock) GetTimestamp() time.Time {
	return block.Timestamp
}

func (block *DefaultBlock) GetState() BlockState {
	return block.State
}

// TODO: Write test case
func (block *DefaultBlock) Serialize() ([]byte, error) {
	data, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// TODO: Write test case
func (block *DefaultBlock) Deserialize(serializedBlock []byte) error {
	if len(serializedBlock) == 0 {
		return ErrDecodingEmptyBlock
	}

	err := json.Unmarshal(serializedBlock, block)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Write test case
func (block *DefaultBlock) IsReadyToPublish() bool {
	return block.Seal != nil
}

func (block *DefaultBlock) IsPrev(serializedPrevBlock []byte) bool {
	prevBlock := &DefaultBlock{}
	prevBlock.Deserialize(serializedPrevBlock)

	return bytes.Compare(prevBlock.GetSeal(), block.GetPrevSeal()) == 0
}

func (block *DefaultBlock) IsEmpty() bool {
	return reflect.DeepEqual(*block, DefaultBlock{})
}

type BlockRepository interface {
	Save(block DefaultBlock) error
	FindLast() (DefaultBlock, error)
	FindByHeight(height BlockHeight) (DefaultBlock, error)
	FindBySeal(seal []byte) (DefaultBlock, error)
	FindAll() ([]DefaultBlock, error)
}

type BlockMap = map[BlockHeight]DefaultBlock

type BlockPool interface {
	Add(block DefaultBlock)
	Delete(height uint64)
	GetByHeight(height uint64) DefaultBlock
	GetSortedKeys() []BlockHeight
	Size() int
}

type SyncState struct {
	SyncProgressing bool
}

func (s *SyncState) Start() {
	s.SyncProgressing = true
}

func (s *SyncState) Done() {
	s.SyncProgressing = false
}

type SyncStateRepository interface {
	Get() SyncState
	Set(state SyncState)
}
