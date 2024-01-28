
package event

import (
	"time"
)

/*
 * consensus
 */

// todo 블록 정보 가지고 있어야 함
// consensus가 끝났다는 event
// true면 블록 저장, false면 블록 저장 안함
type ConsensusFinished struct {
	Seal []byte
	Body []byte
}

/*
 * grpc-gateway
 */

// ivm meta 생성
type ICodeCreated struct {
	ID             string
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        string
}

// ivm meta deleted
type ICodeDeleted struct {
	ICodeID string
}

/*
 * blockChain
 */

// event when block is committed to event store
type BlockCommitted struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []Tx
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   string
	State     string
}

// event when block is staged to event store
type BlockStaged struct {
	BlockId string
	State   string
}

type Tx struct {
	ID        string
	ICodeID   string
	PeerID    string
	TimeStamp time.Time
	Jsonrpc   string
	Function  string
	Args      []string
	Signature []byte
}

/*
 * txpool
 */

// transaction created event
type TxCreated struct {
	TransactionId string
	ICodeID       string
	PeerID        string
	TimeStamp     time.Time
	Jsonrpc       string
	Function      string
	Args          []string
	Signature     []byte
}

// when block committed check transaction and delete
type TxDeleted struct {
	TransactionId string
}

/*
 * p2p
 */

type PeerCreated struct {
	PeerId    string
	IpAddress string
}

type PeerDeleted struct {
	PeerId string
}

// handle leader received event
type LeaderUpdated struct {
	LeaderId string
}

type LeaderDelivered struct {
	LeaderId string
}

type LeaderDeleted struct {
}

//Connection

// connection 생성
type ConnectionCreated struct {
	ConnectionID       string
	GrpcGatewayAddress string
	ApiGatewayAddress  string
}

type ConnectionSaved struct {
	ConnectionID string
	Address      string
}

// connection close
type ConnectionClosed struct {
	ConnectionID string
}

// network
type NetworkJoined struct {
	Connections []ConnectionCreated
}