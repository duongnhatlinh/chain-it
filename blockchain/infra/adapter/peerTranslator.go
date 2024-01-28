package adapter

import (
	"it-chain/api_gateway"
	"it-chain/blockchain"
)

func toPeerFromConnection(peer api_gateway.Peer) blockchain.Peer {
	return blockchain.Peer{
		Id:                peer.ID,
		ApiGatewayAddress: peer.ApiGatewayAddress,
	}
}
