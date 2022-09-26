package app

import (
	"gim/cmd/app/options"
	"gim/internal"
	"github.com/urfave/cli/v2"
)

// NewCommand returns a new instance of *cli.App
func NewCommand(name string) *cli.App {
	// new options
	s := options.NewServerRunOptions()

	app := &cli.App{
		Name:    name,
		Version: internal.Version,
		Usage:   "simple and fast im service",
		Flags:   s.Flags(),
		Action: func(ctx *cli.Context) error {
			// parse command line arguments
			s.ParseArgsFromContext(ctx)

			return run(s)
		},
	}
	return app
}

// run executes command.
func run(opts *options.ServerRunOptions) error {
	// use completedOptions to initialize server instance.
	completedOptions := opts.Complete()
	if err := completedOptions.Validate(); err != nil {
		return err
	}

	// run server with signal.
	server := completedOptions.NewServer()
	server.Run()

	return nil
}
