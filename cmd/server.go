package main

import (
	"gim/cmd/options"
	"gim/internal"
	"gim/internal/infrastructure"
	"gim/internal/infrastructure/config"
	"gim/pkg/system"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	cmd := NewServerCommand()
	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal("failed to run server command: ", err)
	}
}

func NewServerCommand() *cli.App {
	app := &cli.App{
		Name:    "gim",
		Version: internal.Version,
		Usage:   "simple and fast im service",
		Flags:   []cli.Flag{configFlag, portFlag, limitFlag, storageFlag, debugFlag, pushCountFlag},
		Action:  action,
	}
	return app
}
func action(ctx *cli.Context) error {
	s := options.NewServerRunOptions()

	completedOptions, err := Complete(s, ctx)
	if err != nil {
		return err
	}

	if err = run(completedOptions); err != nil {
		return err
	}

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
	return nil
}

func run(completedOptions completedServerRunOptions) error {
	conf := createServerConfig(completedOptions)
	server := internal.NewServer(conf).WithDebug(completedOptions.Debug)
	return server.Run()
}

func createServerConfig(completedOptions completedServerRunOptions) *config.Config {
	conf := config.New()
	completedOptions.ApplyTo(conf)
	return conf
}

type completedServerRunOptions struct {
	*options.ServerRunOptions
}

func (options completedServerRunOptions) Validate() []error {
	return nil
}

func Complete(s *options.ServerRunOptions, ctx *cli.Context) (completedServerRunOptions, error) {
	opts := completedServerRunOptions{}
	opts.ServerRunOptions = s
	opts.Debug = ctx.Bool("debug")
	opts.Port = ctx.Int("port")
	opts.Protocol = ctx.String("protocol")
	opts.MessageMaxStoreSize = ctx.Int("max-store-size")
	opts.MessagePushCount = ctx.Int("push-count")
	opts.MessageStorage = ctx.String("storage")
	return opts, nil
}

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
