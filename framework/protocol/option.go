package protocol

import (
	"net"
	"sync"
)

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

type Property struct {
	bind    string
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

func NewProperty(bind string, handler func(conn net.Conn)) *Property {
	return &Property{
		bind:    bind,
		once:    sync.Once{},
		done:    make(chan struct{}),
		handler: handler,
	}
}
