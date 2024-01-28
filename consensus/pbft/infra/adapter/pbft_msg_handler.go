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
	"errors"

	"it-chain/common"
	"it-chain/common/command"
	"it-chain/consensus/pbft"

	"github.com/DE-labtory/iLogger"
)

var DeserializingError = errors.New("Message deserializing is failed.")
var UndefinedProtocolError = errors.New("Received Undefined protocol message")

type StateMsgApi interface {
	AcceptProposal(msg pbft.ProposeMsg) error
	ReceivePrevote(msg pbft.PrevoteMsg) error
	ReceivePreCommit(msg pbft.PreCommitMsg) error
}

type PbftMsgHandler struct {
	sApi StateMsgApi
}

func NewPbftMsgHandler(sApi StateMsgApi) *PbftMsgHandler {
	return &PbftMsgHandler{
		sApi: sApi,
	}
}

func (p *PbftMsgHandler) HandleGrpcMsgCommand(command command.ReceiveGrpc) error {
	protocol := command.Protocol
	body := command.Body

	switch protocol {

	case "ProposeMsgProtocol":
		iLogger.Infof(nil, "[PBFT] Received protocol - Protocol: [%s]", protocol)

		msg := pbft.ProposeMsg{}
		if err := common.Deserialize(body, &msg); err != nil {
			iLogger.Debugf(nil, "[PBFT] %s", DeserializingError.Error())
		}

		if err := p.sApi.AcceptProposal(msg); err != nil {
			iLogger.Debugf(nil, "[PBFT] %s", err.Error())
		}

	case "PrevoteMsgProtocol":
		iLogger.Infof(nil, "[PBFT] Received protocol - Protocol: [%s]", protocol)

		msg := pbft.PrevoteMsg{}
		if err := common.Deserialize(body, &msg); err != nil {
			iLogger.Debugf(nil, "[PBFT] %s", DeserializingError.Error())
		}

		if err := p.sApi.ReceivePrevote(msg); err != nil {
			iLogger.Debugf(nil, "[PBFT] %s", err.Error())
		}

	case "PreCommitMsgProtocol":
		iLogger.Infof(nil, "[PBFT] Received protocol - Protocol: [%s]", protocol)

		msg := pbft.PreCommitMsg{}
		if err := common.Deserialize(body, &msg); err != nil {
			iLogger.Debugf(nil, "[PBFT] %s", DeserializingError.Error())
		}

		if err := p.sApi.ReceivePreCommit(msg); err != nil {
			iLogger.Debugf(nil, "[PBFT] %s", err.Error())
		}

	default:
		iLogger.Debugf(nil, "[PBFT} %s - Protocol: [%s]", UndefinedProtocolError.Error(), protocol)
	}

	return nil
}
