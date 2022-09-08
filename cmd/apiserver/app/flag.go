package app

import (
	"gim/cmd/apiserver/app/options"
	"gim/internal/infrastructure"
	"github.com/urfave/cli/v2"
	"time"
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
	heartbeatFlag = &cli.DurationFlag{
		Name:  "heartbeat",
		Value: time.Minute,
		Usage: "Set connection heartbeat interval",
	}
)

func appFlags() []cli.Flag {
	return []cli.Flag{configFlag, portFlag, limitFlag, storageFlag, debugFlag, pushCountFlag, heartbeatFlag}
}

// ParseFlagAndCompleteOptions parses the command line arguments and returns *completedServerRunOptions
func ParseFlagAndCompleteOptions(s *options.ServerRunOptions, ctx *cli.Context) (*completedServerRunOptions, error) {
	opts := &completedServerRunOptions{}
	opts.ServerRunOptions = s
	opts.Debug = ctx.Bool("debug")
	opts.Port = ctx.Int("port")
	opts.Protocol = ctx.String("protocol")
	opts.MessageMaxStoreSize = ctx.Int("max-store-size")
	opts.MessagePushCount = ctx.Int("push-count")
	opts.MessageStorage = ctx.String("storage")
	opts.HeartbeatInterval = ctx.Duration("heartbeat")
	return opts, nil
}
