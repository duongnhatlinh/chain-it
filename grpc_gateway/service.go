
package grpc_gateway

const RequestPeerList = "RequestPeerList"
const ResponsePeerList = "ResponsePeerList"

type GrpcService interface {
	Dial(address string) (Connection, error)
	CloseConnection(connID string) error
	SendMessages(message []byte, protocol string, connIDs ...string) error
	GetAllConnections() ([]Connection, error)
	CloseAllConnections() error
	IsConnectionExist(connectionID string) bool
	GetHostID() string
}
