package gateway

type Config struct {
	HttpServerAddress string
	GrpcServerAddress string
	SockServerAddress string
	EnablePprof       bool
	TraceHeader       string
}
