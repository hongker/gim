package framework

// Options represents app options
type Options struct {
	OnConnect         ConnectionHandler
	OnDisconnect      ConnectionHandler
	MaxReadBufferSize int
}

type Option func(options *Options)

func defaultOptions() *Options {
	return &Options{
		OnConnect:         func(conn *Connection) {},
		OnDisconnect:      func(conn *Connection) {},
		MaxReadBufferSize: 512,
	}
}

// WithConnectCallback set OnConnect callback
func WithConnectCallback(onConnect ConnectionHandler) Option {
	return func(options *Options) {
		if onConnect == nil {
			return
		}
		options.OnConnect = onConnect
	}
}

// WithDisconnectCallback set OnDisconnect callback
func WithDisconnectCallback(onDisconnect ConnectionHandler) Option {
	return func(options *Options) {
		if onDisconnect == nil {
			return
		}
		options.OnDisconnect = onDisconnect
	}
}

// WithMaxReadBufferSize set MaxReadBufferSize
func WithMaxReadBufferSize(size int) Option {
	return func(options *Options) {
		options.MaxReadBufferSize = size
	}
}
