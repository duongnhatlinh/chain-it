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

 package blockchainfx

 import (
	 "context"
	 "os"
 
	 "it-chain/api_gateway"
	 "it-chain/blockchain/api"
	 "it-chain/blockchain/infra/adapter"
	 "it-chain/blockchain/infra/mem"
	 "it-chain/blockchain/infra/repo"
	 "it-chain/common"
	 "it-chain/common/rabbitmq/pubsub"
	 "it-chain/conf"
	 "github.com/DE-labtory/iLogger"
	 "go.uber.org/fx"
 )
 
 const publisherID = "publisher.1"
 const BbPath = "./db"
 
 var Module = fx.Options(
	 fx.Provide(
		 NewBlockRepository,
		 NewSyncStateRepository,
		 mem.NewBlockPool,
		 NewBlockAdapter,
		 NewQueryService,
		 NewBlockApi,
		 NewSyncApi,
		 NewConnectionEventHandler,
		 NewBlockProposeHandler,
		 NewConsensusEventHandler,
	 ),
	 fx.Invoke(
		 RegisterPubsubHandlers,
		 RegisterTearDown,
		 CreateGenesisBlock,
	 ),
 )
 
 func NewBlockAdapter() *adapter.HttpBlockAdapter {
	 return adapter.NewHttpBlockAdapter()
 }
 
 func NewQueryService(blockAdapter *adapter.HttpBlockAdapter, peerQueryApi *api_gateway.PeerQueryApi) *adapter.QuerySerivce {
	 return adapter.NewQueryService(blockAdapter, peerQueryApi)
 }
 
 func NewBlockRepository() (*repo.BlockRepository, error) {
 
	 return repo.NewBlockRepository(BbPath)
 }
 
 func NewSyncStateRepository() *mem.SyncStateRepository {
	 return mem.NewSyncStateRepository()
 }
 
 func NewBlockApi(config *conf.Configuration, blockRepository *repo.BlockRepository, blockPool *mem.BlockPool, service common.EventService) (*api.BlockApi, error) {
 
	 NodeId := common.GetNodeID(config.Engine.KeyPath, "ECDSA256")
	 return api.NewBlockApi(NodeId, blockRepository, service, blockPool)
 }
 
 func NewSyncApi(config *conf.Configuration, blockRepository *repo.BlockRepository, syncStateRepository *mem.SyncStateRepository, eventService common.EventService, queryService *adapter.QuerySerivce, blockPool *mem.BlockPool) (*api.SyncApi, error) {
	 NodeId := common.GetNodeID(config.Engine.KeyPath, "ECDSA256")
	 api, err := api.NewSyncApi(NodeId, blockRepository, syncStateRepository, eventService, queryService, blockPool)
	 return &api, err
 }
 
 func NewBlockProposeHandler(blockApi *api.BlockApi, config *conf.Configuration) *adapter.BlockProposeCommandHandler {
	 return adapter.NewBlockProposeCommandHandler(blockApi, config.Engine.Mode)
 }
 
 func NewConnectionEventHandler(syncApi *api.SyncApi) *adapter.NetworkEventHandler {
	 return adapter.NewNetworkEventHandler(syncApi)
 }
 
 func NewConsensusEventHandler(syncStateRepository *mem.SyncStateRepository, blockApi *api.BlockApi) *adapter.ConsensusEventHandler {
	 return adapter.NewConsensusEventHandler(syncStateRepository, blockApi)
 
 }
 
 func CreateGenesisBlock(blockApi *api.BlockApi, config *conf.Configuration) {
	 if err := blockApi.CommitGenesisBlock(config.Blockchain.GenesisConfPath); err != nil {
		 panic(err)
	 }
 }
 
 func RegisterPubsubHandlers(subscriber *pubsub.TopicSubscriber, networkEventHandler *adapter.NetworkEventHandler, blockCommandHandler *adapter.BlockProposeCommandHandler, consensusEventHandler *adapter.ConsensusEventHandler) {
	 iLogger.Infof(nil, "[Main] Blockchain is starting")
 
	 if err := subscriber.SubscribeTopic("network.joined", networkEventHandler); err != nil {
		 panic(err)
	 }
 
	 if err := subscriber.SubscribeTopic("block.propose", blockCommandHandler); err != nil {
		 panic(err)
	 }
 
	 if err := subscriber.SubscribeTopic("block.confirm", consensusEventHandler); err != nil {
		 panic(err)
	 }
 
 }
 
 func RegisterTearDown(lifecycle fx.Lifecycle) {
 
	 lifecycle.Append(fx.Hook{
		 OnStart: func(context context.Context) error {
			 return nil
		 },
		 OnStop: func(context context.Context) error {
			 return os.RemoveAll(BbPath)
		 },
	 })
 }
 