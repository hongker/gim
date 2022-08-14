package main

import (
	"gim/internal"
	"gim/internal/infrastructure"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	run()
}


func run()  {
	app := &cli.App{
		Name:  "gim",
		Version: internal.Version,
		Usage: "simple and fast im service",
		Commands: []*cli.Command{
			{
				Name:    "run",
				Usage:   "run service",
				Action: func(ctx *cli.Context) error {
					internal.NewApp().WithConfigFile(ctx.String("config")).
						WithLimit(ctx.Int("limit")).
						WithPort(ctx.Int("port")).
						WithPushCount(ctx.Int("push-count")).
						WithStorage(ctx.String("storage")).
						WithDebug(ctx.Bool("debug")).
						Run()
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value: "./app.yaml",
						Usage:   "Load configuration from `FILE`",
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value: 8080,
						Usage:   "Set tcp port",
					},
					&cli.IntFlag{
						Name:    "limit",
						Aliases: []string{"l"},
						Value: 10000,
						Usage:   "Set max number of session history messages",
					},
					&cli.StringFlag{
						Name:    "storage",
						Aliases: []string{"s"},
						Value: infrastructure.MemoryStore,
						Usage:   "Set storage",
					},
					&cli.BoolFlag{
						Name:    "debug",
						Value: false,
						Usage:   "Set debug mode",
					},
					&cli.IntFlag{
						Name:    "push-count",
						Value: 5,
						Usage:   "Set count of message push event",
					},

				},
			},
		},
		Action: func(ctx *cli.Context) error {
			return nil
		},

	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
