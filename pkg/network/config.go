package network

import "runtime"

type Config struct {
	Debug           bool
	Bind            string // 服务地址
	Accept          int    // 线程数
	QueueSize       int    // 队列长度
	DataLength      int    // 协议的数据长度
	DataMaxLength   int    // 数据包最大长度
	ContextPoolSize int    // Context对象池大小
	ContentEncoding string // 压缩算法,gzip,zlib

	Sndbuf    int
	Rcvbuf    int
	KeepAlive bool
}

func defaultConfig() *Config {
	return &Config{
		Accept:          runtime.NumCPU(),
		QueueSize:       10,
		DataLength:      0,
		DataMaxLength:   512,
		Sndbuf:          1024,
		Rcvbuf:          1024,
		ContextPoolSize: 32,
	}
}
