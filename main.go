
package main

import (
	"log"
	"os"
	"time"

	"it-chain/cmd/connection"
	"it-chain/cmd/ivm"
	"it-chain/cmd/on"
	"it-chain/common"
	"it-chain/conf"
	"github.com/DE-labtory/iLogger"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "it-chain"
	app.Version = "0.1.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "it-chain",
			Email: "it-chain@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "",
			Usage: "name for config",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "set debug mode",
		},
	}
	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, ivm.IcodeCmd())
	app.Commands = append(app.Commands, connection.Cmd())
	app.Before = func(c *cli.Context) error {
		if configPath := c.String("config"); configPath != "" {
			absPath, err := common.RelativeToAbsolutePath(configPath)
			if err != nil {
				return err
			}
			conf.SetConfigPath(absPath)
		}

		if c.Bool("debug") {
			iLogger.SetToDebug()
		}
		return nil
	}
	app.Action = on.Action
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
