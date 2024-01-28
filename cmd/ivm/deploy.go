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
 
	 "it-chain/common"
	 "it-chain/common/command"
	 "it-chain/common/rabbitmq/rpc"
	 "it-chain/conf"
	 "it-chain/ivm"
	 "github.com/DE-labtory/iLogger"
	 "github.com/urfave/cli"
 )
 
 func DeployCmd() cli.Command {
	 return cli.Command{
		 Name:  "deploy",
		 Usage: "it-chain ivm deploy [icode-git-url] [ssh-path] [password]",
		 Action: func(c *cli.Context) error {
 
			 gitUrl := c.Args().Get(0)
			 sshPath := c.Args().Get(1)
			 password := c.Args().Get(2)
			 deploy(gitUrl, sshPath, password)
 
			 return nil
		 },
	 }
 }
 
 func deploy(gitUrl string, sshPath string, password string) {
 
	 config := conf.GetConfiguration()
	 client := rpc.NewClient(config.Engine.Amqp)
 
	 defer client.Close()
 
	 absPath, err := common.RelativeToAbsolutePath(sshPath)
	 if err != nil {
		 log.Fatal(err.Error())
	 }
 
	 deployCommand := command.Deploy{
		 Url:      gitUrl,
		 SshPath:  absPath,
		 Password: password,
	 }
 
	 iLogger.Infof(nil, "[Cmd] deploying icode...")
	 iLogger.Infof(nil, "[Cmd] This may take a few minutes")
 
	 err = client.Call("ivm.deploy", deployCommand, func(icode ivm.ICode, err rpc.Error) {
 
		 if !err.IsNil() {
			 iLogger.Infof(nil, "fail to deploy icode err: [%s]", err.Message)
			 return
		 }
 
		 iLogger.Infof(nil, "[Cmd] icode has deployed - icodeID: [%s]", icode.ID)
	 })
 
	 if err != nil {
		 log.Fatal(err.Error())
	 }
 }
 