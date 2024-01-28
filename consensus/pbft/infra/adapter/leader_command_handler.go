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
	"it-chain/common"
	"it-chain/common/command"
	"it-chain/consensus/pbft"
)

type LeaderDeliveryParliamentApi interface {
	DeliverLeader(connectionId string)
	UpdateLeader(nodeId string) error
}

type LeaderCommandHandler struct {
	parliamentApi LeaderDeliveryParliamentApi
}

func NewLeaderCommandHandler(parliamentApi LeaderDeliveryParliamentApi) *LeaderCommandHandler {
	return &LeaderCommandHandler{
		parliamentApi: parliamentApi,
	}
}

func (l *LeaderCommandHandler) HandleMessageReceive(command command.ReceiveGrpc) error {

	switch command.Protocol {

	case "RequestLeaderProtocol":
		l.parliamentApi.DeliverLeader(command.ConnectionID)

	case "LeaderDeliveryProtocol":
		message := &pbft.LeaderDeliveryMessage{}
		deserializeErr := common.Deserialize(command.Body, message)
		if deserializeErr != nil {
			return deserializeErr
		}

		l.parliamentApi.UpdateLeader(message.Leader.LeaderId)
	}

	return nil
}
