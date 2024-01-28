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

 package adapter

 import (
	 "sync"
 
	 "it-chain/common/event"
	 "it-chain/ivm"
	 "it-chain/ivm/api"
 )
 
 type BlockCommittedEventHandler struct {
	 icodeApi api.ICodeApi
	 mutex    *sync.Mutex
 }
 
 func NewBlockCommittedEventHandler(icodeApi api.ICodeApi) *BlockCommittedEventHandler {
	 return &BlockCommittedEventHandler{
		 icodeApi: icodeApi,
		 mutex:    &sync.Mutex{},
	 }
 }
 
 func (b *BlockCommittedEventHandler) HandleBlockCommittedEventHandler(blockCommittedEvent event.BlockCommitted) {
	 b.mutex.Lock()
	 defer b.mutex.Unlock()
 
	 b.icodeApi.ExecuteRequestList(createRequestList(blockCommittedEvent.TxList))
 }
 
 func createRequestList(transactionList []event.Tx) []ivm.Request {
 
	 requestList := make([]ivm.Request, 0)
 
	 for _, transaction := range transactionList {
		 requestList = append(requestList, ivm.Request{
			 Function: transaction.Function,
			 Args:     transaction.Args,
			 ICodeID:  transaction.ICodeID,
			 Type:     "invoke",
		 })
	 }
 
	 return requestList
 }
 