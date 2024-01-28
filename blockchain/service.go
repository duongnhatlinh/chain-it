package blockchain

type QueryService interface {
	GetLastBlockFromPeer(peer Peer) (DefaultBlock, error)
	GetBlockByHeightFromPeer(height BlockHeight, peer Peer) (DefaultBlock, error)
}

type EventService interface {
	Publish(topic string, event interface{}) error
}
