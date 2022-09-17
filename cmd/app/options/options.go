package options

import (
	"gim/internal"
	"gim/internal/infrastructure"
	"gim/internal/infrastructure/config"
	"gim/internal/module/gateway"
	"gim/internal/module/message"
	"github.com/urfave/cli/v2"
	"time"
)

// ServerRunOptions run a server.
type ServerRunOptions struct {
	GatewayOptions *gateway.Options
	MessageOptions *message.Options
}

const (
	flagGatewayServerAddress = "gateway-server-address"
)

// Flags returns the command-line flags.
func (ServerRunOptions) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: flagGatewayServerAddress, Aliases: []string{"address"}, Value: ":8080", Usage: "Set server bind address"},
		&cli.IntFlag{Name: "message", Aliases: []string{"l"}, Value: 10000, Usage: "Set max number of session history messages"},
		&cli.IntFlag{Name: "push-count", Value: 5, Usage: "Set count of message push event"},
		&cli.BoolFlag{Name: "debug", Value: false, Usage: "Set debug mode"},
		&cli.StringFlag{Name: "storage", Aliases: []string{"s"}, Value: infrastructure.MemoryStore, Usage: "Set storage, like memory/redis"},
		&cli.DurationFlag{Name: "heartbeat", Value: time.Minute, Usage: "Set connection heartbeat interval"},
	}
}

func (o *ServerRunOptions) ParseArgsFromContext(ctx *cli.Context) {
	o.GatewayOptions.HttpServerAddress = ctx.String(flagGatewayServerAddress)
}

func NewServerRunOptions() *ServerRunOptions {
	o := &ServerRunOptions{
		GatewayOptions: gateway.NewOptions(),
		MessageOptions: message.NewOptions(),
	}
	return o
}

func (o ServerRunOptions) ApplyTo(conf *config.Config) {

}

func (o *ServerRunOptions) Complete() *completedServerRunOptions {
	return &completedServerRunOptions{ServerRunOptions: o}
}

type completedServerRunOptions struct {
	*ServerRunOptions
}

func (o completedServerRunOptions) Validate() error {
	return nil
}

func (o *completedServerRunOptions) NewServer() *internal.Server {
	return internal.NewServer(o.GatewayOptions.BuildInstance())
}