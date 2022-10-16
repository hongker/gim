package framework

// Options represents app options
type Options struct {
	OnConnect         ConnectionHandler
	OnDisconnect      ConnectionHandler
	MaxReadBufferSize int

	Reactor ReactorOptions
}

func (options *Options) NewReactor() *Reactor {
	reactor, err := NewReactor(options.Reactor)
	if err != nil {
		panic(err)
	}
	return reactor
}

type Option func(options *Options)

func defaultOptions() *Options {
	return &Options{
		OnConnect:         func(conn *Connection) {},
		OnDisconnect:      func(conn *Connection) {},
		MaxReadBufferSize: 512,
		Reactor: ReactorOptions{
			EpollBufferSize:  100,
			WorkerPoolSize:   1000,
			PacketLengthSize: 4,
		},
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

func WithEpollBufferSize(size int) Option {
	return func(options *Options) {
		options.Reactor.EpollBufferSize = size
	}
}

func WithWorkerPoolSize(size int) Option {
	return func(options *Options) {
		options.Reactor.WorkerPoolSize = size
	}
}

func WithPacketLengthSize(size int) Option {
	return func(options *Options) {
		options.Reactor.PacketLengthSize = size
	}
}
