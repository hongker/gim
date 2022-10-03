package app

import (
	"gim/cmd/app/options"
	"gim/internal"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/urfave/cli/v2"
)

// NewCommand returns a new instance of *cli.App
func NewCommand(name, usage string) *cli.App {
	// new options
	opts := options.NewServerRunOptions()

	app := &cli.App{
		Name:    name,
		Version: internal.Version,
		Usage:   usage,
		Flags:   opts.Flags(),
		Action: func(ctx *cli.Context) error {
			// parse command line arguments
			opts.ParseArgsFromContext(ctx)

			// run with options
			return run(opts)
		},
	}
	return app
}

// run executes command.
func run(opts *options.ServerRunOptions) error {
	// use completedOptions to initialize server instance.
	completedOptions := opts.Complete()
	return runtime.Call(
		// validate options
		completedOptions.Validate,

		// run server
		completedOptions.NewServer().Run,
	)

}
