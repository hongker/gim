package main

import (
	"gim/internal/infrastructure"
	"github.com/urfave/cli/v2"
)

var (
	configFlag = &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Load configuration from `FILE`",
	}
	portFlag = &cli.IntFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Value:   8080,
		Usage:   "Set tcp port",
	}
	limitFlag = &cli.IntFlag{
		Name:    "limit",
		Aliases: []string{"l"},
		Value:   10000,
		Usage:   "Set max number of session history messages",
	}
	storageFlag = &cli.StringFlag{
		Name:    "storage",
		Aliases: []string{"s"},
		Value:   infrastructure.MemoryStore,
		Usage:   "Set storage, like memory/redis",
	}
	debugFlag = &cli.BoolFlag{
		Name:  "debug",
		Value: false,
		Usage: "Set debug mode",
	}
	pushCountFlag = &cli.IntFlag{
		Name:  "push-count",
		Value: 5,
		Usage: "Set count of message push event",
	}
)
