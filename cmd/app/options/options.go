package options

import (
	"gim/internal"
	"github.com/ebar-go/ego/errors"
	"github.com/urfave/cli/v2"
	"runtime"
	"time"
)

// ServerRunOptions run a server.
type ServerRunOptions struct {
	gatewayAddress    string
	apiAddress        string
	traceHeader       string
	enableProfiling   bool
	heartbeatInterval time.Duration
	workerNumber      int
	gatewayCodec      string
	queuePollInterval time.Duration
	queuePollCount    int
}

const (
	flagGatewayAddress           = "gateway-address"
	flagProfilingEnabled         = "profiling-enabled" //
	flagTraceHeader              = "trace-header"
	flagApiAddress               = "api-address"
	flagGatewayWorker            = "gateway-worker"
	flagGatewayCodec             = "gateway-codec"
	flagGatewayHeartbeatInterval = "gateway-heartbeat-interval"
	flagJobQueuePollInterval     = "job-queue-poll-interval"
	flagJobQueuePollCount        = "job-queue-poll-count"
)

// Flags returns the command-line flags.
func (o *ServerRunOptions) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{Name: flagJobQueuePollCount, Aliases: []string{"queue-poll-count"}, Value: 10, Usage: "Set job queue poll count"},
		&cli.IntFlag{Name: flagGatewayWorker, Aliases: []string{"ws-worker"}, Value: runtime.NumCPU(), Usage: "Set websocket server worker number"},
		&cli.BoolFlag{Name: flagProfilingEnabled, Aliases: []string{"profiling"}, Value: false, Usage: "Set pprof switch"},
		&cli.StringFlag{Name: flagGatewayAddress, Aliases: []string{"ws"}, Value: ":8080", Usage: "Set websocket server bind address"},
		&cli.StringFlag{Name: flagTraceHeader, Aliases: []string{"trace"}, Value: "trace", Usage: "Set trace header"},
		&cli.StringFlag{Name: flagApiAddress, Aliases: []string{"http"}, Value: ":8081", Usage: "Set http server bind address"},
		&cli.StringFlag{Name: flagGatewayCodec, Aliases: []string{"codec"}, Value: "json", Usage: "Set packet codec type(json/protobuf)"},
		&cli.DurationFlag{Name: flagGatewayHeartbeatInterval, Aliases: []string{"heartbeat-interval"}, Value: time.Second * 10, Usage: "Set connection heartbeat interval"},
		&cli.DurationFlag{Name: flagJobQueuePollInterval, Aliases: []string{"queue-poll-interval"}, Value: time.Second, Usage: "Set job queue poll interval"},
	}
}

func (o *ServerRunOptions) ParseArgsFromContext(ctx *cli.Context) {
	o.gatewayAddress = ctx.String(flagGatewayAddress)
	o.traceHeader = ctx.String(flagTraceHeader)
	o.enableProfiling = ctx.Bool(flagProfilingEnabled)
	o.workerNumber = ctx.Int(flagGatewayWorker)
	o.gatewayCodec = ctx.String(flagGatewayCodec)
	o.heartbeatInterval = ctx.Duration(flagGatewayHeartbeatInterval)
	o.queuePollInterval = ctx.Duration(flagJobQueuePollCount)
	o.queuePollCount = ctx.Int(flagJobQueuePollCount)

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

func (o *completedServerRunOptions) Validate() error {
	if o.workerNumber <= 0 {
		return errors.InvalidParam("worker number must be greater than zero")
	}
	return nil
}

func (o *completedServerRunOptions) applyTo(config *internal.Config) {
	config.GatewayControllerConfig.Address = o.gatewayAddress
	config.GatewayControllerConfig.HeartbeatInterval = o.heartbeatInterval

	config.ApiControllerConfig.Address = o.apiAddress
	config.ApiControllerConfig.TraceHeader = o.traceHeader
	config.ApiControllerConfig.EnableProfiling = o.enableProfiling

	config.JobControllerConfig.QueuePollCount = o.queuePollCount
	config.JobControllerConfig.QueuePollInterval = o.queuePollInterval

}

func (o *completedServerRunOptions) NewServer() *internal.Server {
	c := internal.NewConfig()
	o.applyTo(c)

	return c.BuildInstance()
}
