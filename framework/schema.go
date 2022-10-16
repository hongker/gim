package framework

import (
	"fmt"
	"gim/framework/protocol"
	"github.com/ebar-go/ego/utils/runtime"
	"net"
)

const (
	TCP       = "tcp"
	WEBSOCKET = "ws"
	HTTP      = "http"
)

// Schema represents a protocol specification
type Schema struct {
	Protocol string
	Addr     string
}

// Listen run acceptor with handler
func (schema Schema) Listen(stopCh <-chan struct{}, handler func(conn net.Conn)) error {
	var acceptor protocol.Acceptor
	switch schema.Protocol {
	case TCP:
		acceptor = protocol.NewTCPTCPAcceptor(handler)
	case WEBSOCKET:
		acceptor = protocol.NewWSAcceptor(handler)
	default:
		return fmt.Errorf("unsupported protocol: %v", schema.Protocol)
	}

	go func() {
		defer runtime.HandleCrash()
		runtime.WaitClose(stopCh, acceptor.Shutdown)
	}()
	return acceptor.Run(schema.Addr)
}

func NewSchema(protocol string, addr string) Schema {
	return Schema{
		Protocol: protocol,
		Addr:     addr,
	}
}
