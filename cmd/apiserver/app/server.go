package app

import (
	"gim/cmd/apiserver/app/options"
	"gim/internal"
	"gim/internal/infrastructure/config"
	"gim/pkg/runtime/signal"
	"gim/pkg/system"
	"github.com/urfave/cli/v2"
	"log"
)

var (
	appName = "gim"
)

func NewServerCommand() *cli.App {
	app := &cli.App{
		Name:    appName,
		Version: internal.Version,
		Usage:   "simple and fast im service",
		Flags:   appFlags(),
		Action: func(ctx *cli.Context) error {
			// new options
			s := options.NewServerRunOptions()

			// parse command line arguments
			completedOptions, err := ParseFlagAndCompleteOptions(s, ctx)
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
	stopCh := signal.SetupSignalHandler()

	conf := createServerConfig(completedOptions)
	server := createServer(conf)
	go server.Run(stopCh)

	return nil
}

func createServerConfig(completedOptions *completedServerRunOptions) *config.Config {
	conf := config.New()
	completedOptions.ApplyTo(conf)
	return conf
}

func createServer(config *config.Config) *internal.Server {
	return &internal.Server{}
}

type completedServerRunOptions struct {
	*options.ServerRunOptions
}
