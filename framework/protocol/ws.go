package protocol

import (
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/gobwas/ws"
	"log"
	"net"
)

type WebsocketAcceptor struct {
	options  *Options
	property *Property
}

func (acceptor *WebsocketAcceptor) Run() (err error) {
	ln, err := net.Listen("tcp", acceptor.property.bind)
	if err != nil {
		return err
	}
	u := ws.Upgrader{
		OnHeader: func(key, value []byte) (err error) {
			log.Printf("non-websocket header: %q=%q", key, value)
			return
		},
	}

	for i := 0; i < acceptor.options.core; i++ {
		go func() {
			defer runtime.HandleCrash()
			acceptor.accept(ln, u)
		}()
	}
	return nil
}

func (acceptor *WebsocketAcceptor) Shutdown() {
	acceptor.property.Done()
}

func (acceptor *WebsocketAcceptor) accept(ln net.Listener, u ws.Upgrader) {
	for {
		select {
		case <-acceptor.property.Signal():
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("listener.Accept(\"%s\") error(%v)", ln.Addr().String(), err)
				continue
			}

			_, err = u.Upgrade(conn)
			if err != nil {
				log.Printf("upgrade(\"%s\") error(%v)", conn.RemoteAddr().String(), err)
				continue
			}
			acceptor.property.handler(conn)
		}

	}
}

func NewWSAcceptor(bind string, handler func(conn net.Conn)) *WebsocketAcceptor {
	return &WebsocketAcceptor{
		property: NewProperty(bind, handler),
		options: &Options{
			core:            4,
			readBufferSize:  0,
			writeBufferSize: 0,
			keepalive:       false,
		}}
}
