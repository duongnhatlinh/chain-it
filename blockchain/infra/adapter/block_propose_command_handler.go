
 package adapter

 import (
	 "it-chain/blockchain"
	 "it-chain/common/command"
	 "github.com/DE-labtory/iLogger"
 )
 
 type BlockProposeApi interface {
	 CreateProposedBlock(txList []*blockchain.DefaultTransaction) (blockchain.DefaultBlock, error)
	 ConsentBlock(consensusType string, block blockchain.DefaultBlock) error
 }
 
 type BlockProposeCommandHandler struct {
	 blockApi   BlockProposeApi
	 engineMode string
 }
 
 func NewBlockProposeCommandHandler(blockApi BlockProposeApi, engineMode string) *BlockProposeCommandHandler {
	 return &BlockProposeCommandHandler{
		 blockApi:   blockApi,
		 engineMode: engineMode,
	 }
 }
 
 func (h *BlockProposeCommandHandler) HandleProposeBlockCommand(command command.ProposeBlock) error {
 
	 iLogger.Debug(nil, "[Blockchain] Received proposed block command from txpool component")
 
	 if err := validateCommand(command); err != nil {
		 return err
	 }
 
	 txList := command.TxList
	 defaultTxList := getBackTxList(txList)
 
	 proposedBlock, err := h.blockApi.CreateProposedBlock(defaultTxList)
	 if err != nil {
		 return err
	 }
 
	 if err := h.blockApi.ConsentBlock(h.engineMode, proposedBlock); err != nil {
		 return err
	 }
 
	 return nil
 }
 
 func validateCommand(command command.ProposeBlock) error {
	 txList := command.TxList
 
	 if txList == nil || len(txList) == 0 {
		 return ErrCommandTransactions
	 }
	 return nil
 }
 
 func getBackTxList(txList []command.Tx) []*blockchain.DefaultTransaction {
	 defaultTxList := make([]*blockchain.DefaultTransaction, 0)
 
	 for _, tx := range txList {
		 defaultTx := getBackTx(tx)
		 defaultTxList = append(defaultTxList, defaultTx)
	 }
	 return defaultTxList
 }
 
 func getBackTx(tx command.Tx) *blockchain.DefaultTransaction {
	 return &blockchain.DefaultTransaction{
		 ID:        tx.ID,
		 ICodeID:   tx.ICodeID,
		 PeerID:    tx.PeerID,
		 Timestamp: tx.TimeStamp,
		 Jsonrpc:   tx.Jsonrpc,
		 Function:  tx.Function,
		 Args:      tx.Args,
		 Signature: tx.Signature,
	 }
 }
 