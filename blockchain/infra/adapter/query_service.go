

 package adapter

 import (
	 "math/rand"
	 "time"
 
	 "it-chain/api_gateway"
	 "it-chain/blockchain"
 )
 
 type PeerQueryApi interface {
	 GetAllPeerList() []api_gateway.Peer
	 GetPeerByID(connectionId string) (api_gateway.Peer, error)
 }
 
 type BlockAdapter interface {
	 GetLastBlockFromPeer(peer blockchain.Peer) (blockchain.DefaultBlock, error)
	 GetBlockByHeightFromPeer(height blockchain.BlockHeight, peer blockchain.Peer) (blockchain.DefaultBlock, error)
 }
 
 type QuerySerivce struct {
	 blockAdapter BlockAdapter
	 peerQueryApi PeerQueryApi
 }
 
 func NewQueryService(blockAdapter BlockAdapter, peerQueryApi PeerQueryApi) *QuerySerivce {
	 return &QuerySerivce{
		 blockAdapter: blockAdapter,
		 peerQueryApi: peerQueryApi,
	 }
 }
 
 func (s QuerySerivce) GetRandomPeer() (blockchain.Peer, error) {
 
	 peerList := s.peerQueryApi.GetAllPeerList()
	 if len(peerList) == 0 {
		 return blockchain.Peer{}, nil
	 }
 
	 randSource := rand.NewSource(time.Now().UnixNano())
	 randInstance := rand.New(randSource)
	 randomIndex := randInstance.Intn(len(peerList))
	 randomPeer := toPeerFromConnection(peerList[randomIndex])
 
	 return randomPeer, nil
 }
 
 func (s QuerySerivce) GetLastBlockFromPeer(peer blockchain.Peer) (blockchain.DefaultBlock, error) {
 
	 block, err := s.blockAdapter.GetLastBlockFromPeer(peer)
	 if err != nil {
		 return blockchain.DefaultBlock{}, err
	 }
 
	 return block, nil
 }
 
 func (s QuerySerivce) GetBlockByHeightFromPeer(height blockchain.BlockHeight, peer blockchain.Peer) (blockchain.DefaultBlock, error) {
 
	 block, err := s.blockAdapter.GetBlockByHeightFromPeer(height, peer)
	 if err != nil {
		 return blockchain.DefaultBlock{}, err
	 }
 
	 return block, nil
 }
 