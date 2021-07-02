package main

import (
	"os"
	"runtime"

	"github.com/urfave/cli"

	"go-server/library/clean"
)

var (
	cmd   *cli.App
	name  string
	btime string
	build string
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setupCmd()
}

func main() {
	if err := cmd.Run(os.Args); err != nil {
		clean.ExitErr(err)
	}
	clean.Exit()
}
