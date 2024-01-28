package adapter

import (
	"it-chain/blockchain"
	"it-chain/common"
	"it-chain/common/event"

	"github.com/DE-labtory/sdk/logger"
)

type BlockApiForCommitAndStage interface {
	CommitBlock(block blockchain.DefaultBlock) error
	StageBlock(block blockchain.DefaultBlock)
}

type ConsensusEventHandler struct {
	SyncStateRepository blockchain.SyncStateRepository
	BlockApi            BlockApiForCommitAndStage
}

func NewConsensusEventHandler(syncStateRepository blockchain.SyncStateRepository, blockApi BlockApiForCommitAndStage) *ConsensusEventHandler {

	return &ConsensusEventHandler{
		SyncStateRepository: syncStateRepository,
		BlockApi:            blockApi,
	}
}

/**
receive consensus finished event
if block sync is on progress, change state to 'staged' and add to block pool
if block sync is not on progress, commit block
*/
func (c *ConsensusEventHandler) HandleConsensusFinishedEvent(event event.ConsensusFinished) error {
	receivedBlock := extractBlockFromEvent(event)

	if receivedBlock.Seal == nil {
		return ErrBlockSealNil
	}

	syncState := c.SyncStateRepository.Get()

	if syncState.SyncProgressing {
		c.BlockApi.StageBlock(*receivedBlock)
	} else {
		if err := c.BlockApi.CommitBlock(*receivedBlock); err != nil {
			return err
		}
	}

	return nil
}

func extractBlockFromEvent(event event.ConsensusFinished) *blockchain.DefaultBlock {
	block := &blockchain.DefaultBlock{}

	if err := common.Deserialize(event.Body, block); err != nil {
		logger.Error(nil, "[Blockchain] Deserialize error")
	}

	return block
}
