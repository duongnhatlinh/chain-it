package api

import (
	"math/rand"
	"time"

	"it-chain/blockchain"

	"github.com/DE-labtory/iLogger"
	"github.com/pkg/errors"
)

type SyncApi struct {
	publisherId         string
	blockRepository     blockchain.BlockRepository
	syncStateRepository blockchain.SyncStateRepository
	eventService        blockchain.EventService
	queryService        blockchain.QueryService
	blockPool           blockchain.BlockPool
}

func NewSyncApi(publisherId string, blockRepository blockchain.BlockRepository, syncStateRepository blockchain.SyncStateRepository, eventService blockchain.EventService, queryService blockchain.QueryService, blockPool blockchain.BlockPool) (SyncApi, error) {
	return SyncApi{
		publisherId:         publisherId,
		blockRepository:     blockRepository,
		syncStateRepository: syncStateRepository,
		eventService:        eventService,
		queryService:        queryService,
		blockPool:           blockPool,
	}, nil
}

func (sApi SyncApi) Synchronize(peer blockchain.Peer) error {
	syncState := sApi.syncStateRepository.Get()

	syncState.Start()
	sApi.syncStateRepository.Set(syncState)

	defer func() {
		syncState.Done()
		sApi.syncStateRepository.Set(syncState)
	}()

	iLogger.Infof(nil, "[Blockchain] Start to synchronize - Peer: [%s]", peer.Id)

	if sApi.isSynced(peer) {
		iLogger.Infof(nil, "[Blockchain] Already synchronized - Peer: [%s]", peer.Id)
		return nil
	}

	// if sync has not done, on sync
	err := sApi.syncWithPeer(peer)
	if err != nil {
		iLogger.Errorf(nil, "[Blockchain] Fail to synchronize - Err: [%s]", err)
		return err
	}

	if err = sApi.CommitStagedBlocks(); err != nil {
		return err
	}

	iLogger.Infof(nil, "[Blockchain] Synchronized Successfully - Peer: [%s]", peer.Id)

	return nil
}

func (sApi SyncApi) isSynced(peer blockchain.Peer) bool {

	// If nil peer is given(when i'm the first node of p2p network) : Synced
	if peer.ApiGatewayAddress == "" {
		return true
	}
	// Get last block of my blockChain
	lastBlock, err := sApi.blockRepository.FindLast()
	if err != nil {
		return false
	}
	// Get last block of other peer's blockChain
	standardBlock, err := sApi.queryService.GetLastBlockFromPeer(peer)
	if err != nil {
		return false
	}
	// Compare last block vs standard block
	if lastBlock.GetHeight() < standardBlock.GetHeight() {
		return false
	}

	return true
}

func (sApi SyncApi) syncWithPeer(peer blockchain.Peer) error {
	standardBlock, err := sApi.queryService.GetLastBlockFromPeer(peer)

	if err != nil {
		return err
	}

	standardHeight := standardBlock.GetHeight()
	lastBlock, err := sApi.blockRepository.FindLast()

	if err != nil {
		return err
	}

	lastHeight := lastBlock.GetHeight()

	return sApi.construct(peer, standardHeight, lastHeight)

}

func (sApi SyncApi) construct(peer blockchain.Peer, standardHeight blockchain.BlockHeight, lastHeight blockchain.BlockHeight) error {

	for lastHeight < standardHeight {

		targetHeight := setTargetHeight(lastHeight)
		retrievedBlock, err := sApi.queryService.GetBlockByHeightFromPeer(targetHeight, peer)
		if err != nil {
			return err
		}

		err = sApi.commitBlock(retrievedBlock)
		if err != nil {
			return err
		}

		raiseHeight(&lastHeight)
	}

	return nil
}

func setTargetHeight(lastHeight blockchain.BlockHeight) blockchain.BlockHeight {
	return lastHeight + 1
}
func raiseHeight(height *blockchain.BlockHeight) {
	*height++
}

func (sApi SyncApi) commitBlock(block blockchain.DefaultBlock) error {

	// save(commit)
	err := sApi.blockRepository.Save(block)
	if err != nil {
		iLogger.Errorf(nil, "[Blockchain] Block is not committed - Err: [%s]", err)
		return ErrSaveBlock
	}

	// publish
	commitEvent, err := createBlockCommittedEvent(block)
	if err != nil {
		return ErrCreateEvent
	}

	iLogger.Infof(nil, "[Blockchain] Block has committed - Seal: [%x],  Height: [%d]", block.Seal, block.Height)

	return sApi.eventService.Publish("block.committed", commitEvent)
}

func (sApi *SyncApi) CommitStagedBlocks() error {
	lastBlock, err := sApi.blockRepository.FindLast()
	if err != nil {
		return err
	}

	targetHeight := setTargetHeight(lastBlock.GetHeight())

	for _, h := range sApi.blockPool.GetSortedKeys() {
		height := blockchain.BlockHeight(h)
		switch {
		case height > targetHeight:
			return nil

		case height < targetHeight:
			sApi.blockPool.Delete(height)

		case height == targetHeight:
			block := sApi.blockPool.GetByHeight(height)
			if err := sApi.commitBlock(block); err != nil {
				return err
			}

			sApi.blockPool.Delete(height)
			raiseHeight(&targetHeight)
		}
	}

	return nil
}

func (sApi *SyncApi) HandleNetworkJoined(peerList []blockchain.Peer) error {
	peer, err := getRandomPeer(peerList)
	if err != nil {
		return err
	}

	return sApi.Synchronize(peer)
}

func getRandomPeer(peerList []blockchain.Peer) (blockchain.Peer, error) {
	if len(peerList) == 0 {
		return blockchain.Peer{}, errors.New("No peer")
	}

	randSource := rand.NewSource(time.Now().UnixNano())
	randInstance := rand.New(randSource)
	randomIndex := randInstance.Intn(len(peerList))
	randomPeer := peerList[randomIndex]

	return randomPeer, nil
}
