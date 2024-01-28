/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

 package txpoolfx

 import (
	 "context"
	 "time"
 
	 "it-chain/common"
	 "it-chain/common/batch"
	 "it-chain/common/rabbitmq/pubsub"
	 "it-chain/common/rabbitmq/rpc"
	 "it-chain/conf"
	 "it-chain/txpool"
	 "it-chain/txpool/api"
	 "it-chain/txpool/infra/adapter"
	 "it-chain/txpool/infra/mem"
	 "github.com/DE-labtory/iLogger"
	 "go.uber.org/fx"
 )
 
 var Module = fx.Options(
	 fx.Provide(
		 mem.NewTransactionRepository,
		 NewLeaderRepository,
		 NewBlockProposalService,
		 NewTransferService,
		 NewTxpoolApi,
		 NewGrpcMessageHandler,
		 NewLeaderEventHandler,
		 adapter.NewTxCommandHandler,
	 ),
	 fx.Invoke(
		 RunBatcher,
		 RegisterRpcHandlers,
		 RegisterPubsubHandlers,
	 ),
 )
 
 func NewLeaderRepository(config *conf.Configuration) *mem.LeaderRepository {
	 NodeId := common.GetNodeID(config.Engine.KeyPath, "ECDSA256")
	 repo := mem.NewLeaderRepository()
	 if config.Engine.BootstrapNodeAddress == "" {
		 repo.Set(txpool.Leader{NodeId})
	 }
 
	 return repo
 }
 
 func NewBlockProposalService(repository *mem.TransactionRepository, eventService common.EventService) *txpool.BlockProposalService {
	 return txpool.NewBlockProposalService(repository, eventService)
 }
 
 func NewTransferService(transactionRepository *mem.TransactionRepository, leaderRepository *mem.LeaderRepository, eventService common.EventService) *txpool.TransferService {
	 return txpool.NewTransferService(transactionRepository, leaderRepository, eventService)
 }
 
 func NewTxpoolApi(config *conf.Configuration, transactionRepository *mem.TransactionRepository, leaderRepository *mem.LeaderRepository, transferService *txpool.TransferService, blockProposalService *txpool.BlockProposalService) *api.TransactionApi {
	 NodeId := common.GetNodeID(config.Engine.KeyPath, "ECDSA256")
	 return api.NewTransactionApi(NodeId, transactionRepository, leaderRepository, transferService, blockProposalService)
 }
 
 func NewLeaderEventHandler(leaderRepository *mem.LeaderRepository) *adapter.LeaderEventHandler {
 
	 return adapter.NewLeaderEventHandler(leaderRepository)
 }
 
 func NewGrpcMessageHandler(txPoolApi *api.TransactionApi) *adapter.GrpcMessageHandler {
	 return adapter.NewGrpcMessageHandler(txPoolApi)
 }
 
 func RunBatcher(lifecycle fx.Lifecycle, txPoolApi *api.TransactionApi, config *conf.Configuration) {
 
	 var proposeBlockQuit chan struct{}
	 var sendTransactionQuit chan struct{}
	 lifecycle.Append(fx.Hook{
		 OnStart: func(context context.Context) error {
			 proposeBlockQuit = batch.GetTimeOutBatcherInstance().Run(func() error {
				 return txPoolApi.ProposeBlock(config.Engine.Mode)
			 }, (time.Duration(config.Txpool.TimeoutMs) * time.Millisecond))
 
			 sendTransactionQuit = batch.GetTimeOutBatcherInstance().Run(func() error {
				 return txPoolApi.SendLeaderTransaction(config.Engine.Mode)
			 }, (time.Duration(config.Txpool.TimeoutMs) * time.Millisecond))
			 return nil
		 },
		 OnStop: func(context context.Context) error {
			 proposeBlockQuit <- struct{}{}
			 sendTransactionQuit <- struct{}{}
			 return nil
		 },
	 })
 }
 
 func RegisterRpcHandlers(server *rpc.Server, handler *adapter.TxCommandHandler) {
	 iLogger.Infof(nil, "[Main] Txpool is starting")
	 if err := server.Register("transaction.create", handler.HandleTxCreateCommand); err != nil {
		 panic(err)
	 }
 }
 
 func RegisterPubsubHandlers(subscriber *pubsub.TopicSubscriber, leaderEventHandler *adapter.LeaderEventHandler, grpcMessageHandler *adapter.GrpcMessageHandler) {
 
	 if err := subscriber.SubscribeTopic("leader.updated", leaderEventHandler); err != nil {
		 panic(err)
	 }
 
	 if err := subscriber.SubscribeTopic("message.receive", grpcMessageHandler); err != nil {
		 panic(err)
	 }
 
 }
 