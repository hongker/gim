package protocol

import (
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/pkg/errors"
	"log"
	"net"
)

type TCPAcceptor struct {
	options  *Options
	property *Property
}

func (server *TCPAcceptor) Run() (err error) {
	addr, err := net.ResolveTCPAddr("tcp", server.property.bind)
	if err != nil {
		return errors.WithMessage(err, "resolve tcp addr")
	}

	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}

	for i := 0; i < server.options.core; i++ {
		go func() {
			defer runtime.HandleCrash()
			server.accept(lis)
		}()
	}

	return
}

func (acceptor *TCPAcceptor) Shutdown() {
	acceptor.property.Done()
}

func (acceptor *TCPAcceptor) accept(lis *net.TCPListener) {
	var (
		conn *net.TCPConn
		err  error
	)

	for {
		select {
		case <-acceptor.property.Signal():
			return
		default:
			if conn, err = lis.AcceptTCP(); err != nil {
				// if listener close then return
				log.Printf("listener.Accept(\"%s\") error(%v)", lis.Addr().String(), err)
				continue
			}
			if err = conn.SetKeepAlive(acceptor.options.keepalive); err != nil {
				log.Printf("conn.SetKeepAlive() error(%v)", err)
				continue
			}
			if err = conn.SetReadBuffer(acceptor.options.readBufferSize); err != nil {
				log.Printf("conn.SetReadBuffer() error(%v)", err)
				continue
			}
			if err = conn.SetWriteBuffer(acceptor.options.writeBufferSize); err != nil {
				log.Printf("conn.SetWriteBuffer() error(%v)", err)
				continue
			}

			acceptor.property.handler(conn)
		}
	}

}

func NewTCPTCPAcceptor(bind string) *TCPAcceptor {
	return &TCPAcceptor{
		property: NewProperty(bind),
		options: &Options{
			core:            4,
			readBufferSize:  0,
			writeBufferSize: 0,
			keepalive:       false,
		}}
}
