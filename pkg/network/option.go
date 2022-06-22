package network

type Option func(conf *Config)

func WithAcceptCount(count int) Option {
	return func(conf *Config) {
		conf.Accept = count
	}
}

func WithQueueSize(size int) Option {
	return func(conf *Config) {
		conf.QueueSize = size
	}
}

func WithPacketLength(size int) Option {
	return func(conf *Config) {
		conf.DataLength = size
	}
}

func WithDebug() Option {
	return func(conf *Config) {
		conf.Debug = true
	}
}
