
package txpool

import (
	"time"

	"github.com/rs/xid"
)

type TransactionId = string

type TxData struct {
	Jsonrpc   string
	ICodeID   string
	Function  string
	Args      []string
	Signature []byte
}

//Aggregate root must implement aggregate interface
type Transaction struct {
	ID        TransactionId
	TimeStamp time.Time
	Jsonrpc   string
	ICodeID   string
	Function  string
	Args      []string
	Signature []byte
	PeerID    string
}

func CreateTransaction(publisherId string, txData TxData) (Transaction, error) {

	id := xid.New().String()
	timeStamp := time.Now()

	transaction := Transaction{
		ID:        id,
		PeerID:    publisherId,
		TimeStamp: timeStamp,
		ICodeID:   txData.ICodeID,
		Jsonrpc:   txData.Jsonrpc,
		Signature: txData.Signature,
		Args:      txData.Args,
		Function:  txData.Function,
	}

	return transaction, nil
}

type TransactionRepository interface {
	FindAll() ([]Transaction, error)
	Save(transaction Transaction) error
	Remove(id TransactionId)
	FindById(id TransactionId) (Transaction, error)
}
