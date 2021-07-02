// Author: Steve Zhang
// Date: 2020/9/17 2:27 下午

package main

import (
	"fmt"

	"github.com/urfave/cli"

	"go-server/common"
	"go-server/component"
	"go-server/library/log"
)

func setupCmd() {
	cmd = cli.NewApp()
	cmd.Name = name
	cmd.Version = build
	cmd.Commands = []cli.Command{
		{
			Name:     "start",
			HideHelp: true,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "c", Value: ".env", Usage: "config file"},
				cli.StringFlag{Name: "p", Value: "3000", Usage: "http listen port"},
			},
			Before: func(ctx *cli.Context) (err error) {
				err = setupComponent(ctx.String("c"), ctx.Int("p"))
				return
			},
			Action: func(ctx *cli.Context) (err error) {
				router := setupRouter()
				component.InfLogger.Info(log.F{
					"log_type":   common.LogTypeForAppStart,
					"name":       name,
					"version":    build,
					"build_time": btime,
				})
				if err = component.HttpServer.Run(router); err != nil {
					err = fmt.Errorf("common.HttpServer.Run: %w", err)
				}
				return
			},
		},
	}
}
