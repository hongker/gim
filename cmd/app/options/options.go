package options

import (
	"gim/internal/aggregator"
	"github.com/urfave/cli/v2"
	"time"
)

// ServerRunOptions run a server.
type ServerRunOptions struct {
	gatewayAddress    string
	apiAddress        string
	traceHeader       string
	enableProfiling   bool
	heartbeatInterval time.Duration
}

const (
	flagGatewayAddress   = "gateway-address"
	flagProfilingEnabled = "profiling-enabled" //
	flagTraceHeader      = "trace-header"
	flagApiAddress       = "api-address"
)

// Flags returns the command-line flags.
func (ServerRunOptions) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: flagGatewayAddress, Aliases: []string{"ws"}, Value: ":8080", Usage: "Set websocket server bind address"},
		&cli.StringFlag{Name: flagTraceHeader, Aliases: []string{"trace"}, Value: "trace", Usage: "Set trace header"},
		&cli.BoolFlag{Name: flagProfilingEnabled, Aliases: []string{"profiling"}, Value: false, Usage: "Set pprof switch"},
		&cli.StringFlag{Name: flagApiAddress, Aliases: []string{"http"}, Value: ":8081", Usage: "Set http server bind address"},
		//&cli.IntFlag{Name: "message", Aliases: []string{"l"}, Value: 10000, Usage: "Set max number of session history messages"},
		//&cli.IntFlag{Name: "push-count", Value: 5, Usage: "Set count of message push event"},
		//&cli.BoolFlag{Name: "debug", Value: false, Usage: "Set debug mode"},
		//&cli.StringFlag{Name: "storage", Aliases: []string{"s"}, Value: infrastructure.MemoryStore, Usage: "Set storage, like memory/redis"},
		//&cli.DurationFlag{Name: "heartbeat", Value: time.Minute, Usage: "Set connection heartbeat interval"},
	}
}

func (o *ServerRunOptions) ParseArgsFromContext(ctx *cli.Context) {
	o.gatewayAddress = ctx.String(flagGatewayAddress)
	o.traceHeader = ctx.String(flagTraceHeader)
	o.enableProfiling = ctx.Bool(flagProfilingEnabled)
}

func NewServerRunOptions() *ServerRunOptions {
	o := &ServerRunOptions{}
	return o
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

func (o completedServerRunOptions) applyTo(config *aggregator.Config) {
	config.GatewayControllerConfig.Address = o.gatewayAddress
	config.GatewayControllerConfig.HeartbeatInterval = o.heartbeatInterval

	config.ApiControllerConfig.Address = o.apiAddress
	config.ApiControllerConfig.TraceHeader = o.traceHeader
	config.ApiControllerConfig.EnableProfiling = o.enableProfiling

}

func (o *completedServerRunOptions) NewServer() *aggregator.Aggregator {
	c := aggregator.NewConfig()
	o.applyTo(c)

	return c.New()
}
