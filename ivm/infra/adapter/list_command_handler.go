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
	 "it-chain/common/command"
	 "it-chain/common/rabbitmq/rpc"
	 "it-chain/ivm/api"
 )
 
 type ListCommandHandler struct {
	 icodeApi api.ICodeApi
 }
 
 func NewListCommandHandler(icodeApi api.ICodeApi) *ListCommandHandler {
	 return &ListCommandHandler{
		 icodeApi: icodeApi,
	 }
 }
 
 func (l *ListCommandHandler) HandleListCommand(getICodeListCommand command.GetICodeList) (command.ICodeList, rpc.Error) {
	 iCodes := l.icodeApi.GetRunningICodeList()
 
	 return command.ICodeList{ICodes: iCodes}, rpc.Error{}
 }
 