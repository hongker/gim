package app

import (
	"gim/cmd/apiserver/app/options"
	"gim/internal"
	"gim/internal/infrastructure/config"
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
