package protocol

import (
	"net"
	"sync"
)

type Acceptor interface {
	Run(bind string) error
	Shutdown()
}
type Options struct {
	core            int
	readBufferSize  int
	writeBufferSize int
	keepalive       bool
}

type Property struct {
	once    sync.Once
	done    chan struct{}
	handler func(conn net.Conn)
}

func (p *Property) Signal() <-chan struct{} {
	return p.done
}

func (p *Property) Done() {
	p.once.Do(func() {
		close(p.done)
	})
}

func NewProperty(handler func(conn net.Conn)) *Property {
	return &Property{
		once:    sync.Once{},
		done:    make(chan struct{}),
		handler: handler,
	}
}
