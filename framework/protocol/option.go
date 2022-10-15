package protocol

type Acceptor interface {
	Run() error
	Shutdown()
}
type Options struct {
	core            int
	readBufferSize  int
	writeBufferSize int
	keepalive       bool
}
