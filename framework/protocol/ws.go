package protocol

import (
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/gobwas/ws"
	"log"
	"net"
	"sync"
)

type WebsocketAcceptor struct {
	bind    string
	options *Options
	once    sync.Once
	done    chan struct{}
	handler func(conn net.Conn)
}

func (acceptor *WebsocketAcceptor) Run() (err error) {
	ln, err := net.Listen("tcp", acceptor.bind)
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
	acceptor.once.Do(func() {
		close(acceptor.done)
	})
}

func (acceptor *WebsocketAcceptor) accept(ln net.Listener, u ws.Upgrader) {
	for {
		select {
		case <-acceptor.done:
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
			acceptor.handler(conn)
		}

	}
}

func NewWSAcceptor(bind string) *WebsocketAcceptor {
	return &WebsocketAcceptor{bind: bind, options: &Options{
		core:            4,
		readBufferSize:  0,
		writeBufferSize: 0,
		keepalive:       false,
	}}
}
