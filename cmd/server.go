package main

import (
	"gim/cmd/options"
	"gim/internal"
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
		Flags:   []cli.Flag{configFlag, portFlag, limitFlag, storageFlag, debugFlag, pushCountFlag, heartbeatFlag},
		Action: func(ctx *cli.Context) error {
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
		},
	}
	return app
}

func run(completedOptions *completedServerRunOptions) error {
	conf := createServerConfig(completedOptions)
	return internal.Run(conf)
}

func createServerConfig(completedOptions *completedServerRunOptions) *config.Config {
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

func Complete(s *options.ServerRunOptions, ctx *cli.Context) (*completedServerRunOptions, error) {
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
