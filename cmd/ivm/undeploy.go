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

 package ivm

 import (
	 "log"
 
	 "it-chain/common/command"
	 "it-chain/common/rabbitmq/rpc"
	 "it-chain/conf"
	 "github.com/urfave/cli"
 )
 
 func UnDeployCmd() cli.Command {
	 return cli.Command{
		 Name:  "undeploy",
		 Usage: "it-chain ivm undeploy [icode-id] ",
		 Action: func(c *cli.Context) error {
 
			 icodeId := c.Args().Get(0)
			 unDeploy(icodeId)
			 return nil
		 },
	 }
 }
 func unDeploy(icodeId string) {
 
	 config := conf.GetConfiguration()
	 client := rpc.NewClient(config.Engine.Amqp)
 
	 defer client.Close()
 
	 undeployCommand := command.UnDeploy{
		 ICodeId: icodeId,
	 }
 
	 err := client.Call("ivm.undeploy", undeployCommand, func(empty struct{}, err rpc.Error) {
 
		 if !err.IsNil() {
			 log.Printf("fail to undeploy icode err: [%s]", err.Message)
			 return
		 }
 
		 log.Printf("[%s] icode has undeployed", icodeId)
	 })
 
	 if err != nil {
		 log.Fatal(err.Error())
	 }
 }
 