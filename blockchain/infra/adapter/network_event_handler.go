

 package adapter

 import (
	 "it-chain/blockchain"
	 "it-chain/common/event"
	 "github.com/DE-labtory/iLogger"
 )
 
 type SynchronizeApi interface {
	 HandleNetworkJoined(peer []blockchain.Peer) error
 }
 
 type NetworkEventHandler struct {
	 SyncApi SynchronizeApi
 }
 
 func NewNetworkEventHandler(syncApi SynchronizeApi) *NetworkEventHandler {
 
	 return &NetworkEventHandler{
		 SyncApi: syncApi,
	 }
 }
 
 func (n *NetworkEventHandler) HandleNetworkJoinedEvent(networkJoindEvent event.NetworkJoined) {
	 iLogger.Infof(nil, "[Blockchain] Network Joined")
	 if err := n.SyncApi.HandleNetworkJoined(createPeerListFromNetworkJoinedEvent(networkJoindEvent)); err != nil {
		 iLogger.Errorf(nil, "[Blockchain] Fail to handle network joined event - Err: [%s]", err.Error())
	 }
 }
 
 func createPeerListFromNetworkJoinedEvent(networkJoindEvent event.NetworkJoined) []blockchain.Peer {
	 peerList := make([]blockchain.Peer, 0)
 
	 for _, c := range networkJoindEvent.Connections {
		 peerList = append(peerList, blockchain.Peer{
			 Id:                c.ConnectionID,
			 ApiGatewayAddress: c.ApiGatewayAddress,
		 })
	 }
 
	 return peerList
 }
 