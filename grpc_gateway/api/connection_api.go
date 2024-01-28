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

 package api

 import (
	 "encoding/json"
 
	 "it-chain/common"
	 "it-chain/common/event"
	 "it-chain/grpc_gateway"
	 "github.com/DE-labtory/iLogger"
 )
 
 type ConnectionApi struct {
	 grpcService  grpc_gateway.GrpcService
	 eventService common.EventService
 }
 
 func NewConnectionApi(grpcService grpc_gateway.GrpcService, eventService common.EventService) *ConnectionApi {
	 return &ConnectionApi{
		 grpcService:  grpcService,
		 eventService: eventService,
	 }
 }
 
 // create all connections
 func (c ConnectionApi) JoinNetwork(address string) error {
	 iLogger.Infof(nil, "[gRPC-Gateway] Joining it-chain network - Address: [%s]", address)
 
	 connection, err := c.Dial(address)
	 if err != nil {
		 iLogger.Errorf(nil, "[gRPC-Gateway] Fail to join - Err: [%s]", err)
		 return err
	 }
 
	 c.grpcService.SendMessages([]byte(""), grpc_gateway.RequestPeerList, connection.ConnectionID)
	 return nil
 }
 
 // create connection only for the address
 func (c ConnectionApi) Dial(address string) (grpc_gateway.Connection, error) {
 
	 iLogger.Infof(nil, "[gRPC-Gateway] Dialing - Address: [%s]", address)
 
	 connection, err := c.grpcService.Dial(address)
	 if err != nil {
		 iLogger.Errorf(nil, "[gRPC-Gateway] Fail to dial - Err: [%s]", err)
		 return grpc_gateway.Connection{}, err
	 }
 
	 err = c.eventService.Publish("connection.created", createConnectionCreatedEvent(connection))
	 if err != nil {
		 return connection, err
	 }
 
	 iLogger.Infof(nil, "[gRPC-Gateway] Connection created - gRPC-Address [%s], ConnectionID [%s]", connection.GrpcGatewayAddress, connection.ConnectionID)
 
	 return connection, nil
 }
 
 func createConnectionCreatedEvent(connection grpc_gateway.Connection) event.ConnectionCreated {
	 return event.ConnectionCreated{
		 ConnectionID:       connection.ConnectionID,
		 GrpcGatewayAddress: connection.GrpcGatewayAddress,
		 ApiGatewayAddress:  connection.ApiGatewayAddress,
	 }
 }
 
 func (c ConnectionApi) CloseConnection(connectionID string) error {
	 iLogger.Infof(nil, "[gRPC-Gateway] Disconnect connection - ConnectionID [%s]", connectionID)
 
	 err := c.grpcService.CloseConnection(connectionID)
	 if err != nil {
		 return err
	 }
 
	 return c.eventService.Publish("connection.closed", createConnectionClosedEvent(connectionID))
 }
 
 func createConnectionClosedEvent(connectionID string) event.ConnectionClosed {
	 return event.ConnectionClosed{
		 ConnectionID: connectionID,
	 }
 }
 
 func (c ConnectionApi) OnConnection(connection grpc_gateway.Connection) {
	 iLogger.Infof(nil, "[gRPC-Gateway] Connection created - gRPC-Address [%s], ConnectionID [%s]", connection.GrpcGatewayAddress, connection.ConnectionID)
 
	 if err := c.eventService.Publish("connection.created", createConnectionCreatedEvent(connection)); err != nil {
		 iLogger.Infof(nil, "[gRPC-Gateway] Fail to publish connection createdEvent - ConnectionID: [%s]", connection.ConnectionID)
	 }
 }
 
 func (c ConnectionApi) OnDisconnection(connection grpc_gateway.Connection) {
	 iLogger.Infof(nil, "[gRPC-Gateway] Connection closed - ConnectionID [%s]", connection.ConnectionID)
 
	 if err := c.eventService.Publish("connection.closed", createConnectionClosedEvent(connection.ConnectionID)); err != nil {
		 iLogger.Infof(nil, "[gRPC-Gateway] Fail to publish connection createdEvent - ConnectionID: [%s]", connection.ConnectionID)
	 }
 }
 
 func (c ConnectionApi) GetAllConnections() ([]grpc_gateway.Connection, error) {
	 return c.grpcService.GetAllConnections()
 }
 
 func (c ConnectionApi) HandleRequestPeerList(connectionId string) {
	 iLogger.Infof(nil, "[gRPC-Gateway] Handling peerList reqeust - ConnectionId: [%s]", connectionId)
 
	 connectionList, _ := c.grpcService.GetAllConnections()
	 response, err := json.Marshal(connectionList)
	 if err != nil {
		 iLogger.Errorf(nil, "[gRPC-Gateway] Fail to handle request peer list - Err: [%s]", err.Error())
		 return
	 }
 
	 c.grpcService.SendMessages(response, grpc_gateway.ResponsePeerList, connectionId)
 }
 
 // 자기자신 or 연결되어 있는 connection 제외하고 연결!!
 func (c ConnectionApi) DialConnectionList(connectionList []grpc_gateway.Connection) {
	 iLogger.Infof(nil, "[gRPC-Gateway] Dialing all peers in it-chain network - Total peer: [%d]", len(connectionList))
	 for _, connection := range connectionList {
		 //remove already connected connection
		 if c.grpcService.IsConnectionExist(connection.ConnectionID) {
			 continue
		 }
		 //자기 자신
		 if c.grpcService.GetHostID() == connection.ConnectionID {
			 continue
		 }
		 _, err := c.Dial(connection.GrpcGatewayAddress)
		 if err != nil {
			 return
		 }
	 }
 
	 connections, err := c.GetAllConnections()
	 if err != nil {
		 iLogger.Errorf(nil, "[gRPC-Gateway] Fail to get connections - Err: [%s]", err.Error())
		 return
	 }
 
	 if err := c.eventService.Publish("network.joined", createNetworkjoinedEvent(connections)); err != nil {
		 iLogger.Errorf(nil, "[gRPC-Gateway] Fail to publish network joined event - Err: [%s]", err.Error())
	 }
 }
 
 func createNetworkjoinedEvent(connectionList []grpc_gateway.Connection) event.NetworkJoined {
	 connectionCreatedEvents := make([]event.ConnectionCreated, 0)
	 for _, c := range connectionList {
		 connectionCreatedEvents = append(connectionCreatedEvents, event.ConnectionCreated{
			 ConnectionID:       c.ConnectionID,
			 ApiGatewayAddress:  c.ApiGatewayAddress,
			 GrpcGatewayAddress: c.GrpcGatewayAddress,
		 })
	 }
 
	 return event.NetworkJoined{
		 Connections: connectionCreatedEvents,
	 }
 }
 