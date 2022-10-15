package framework

import (
	"fmt"
	"gim/framework/protocol"
	"github.com/ebar-go/ego/utils/runtime"
)

type Protocol string

const (
	TCP       Protocol = "tcp"
	WEBSOCKET Protocol = "ws"
	HTTP      Protocol = "http"
)

type Schema struct {
	Protocol Protocol
	Addr     string
}

func (schema Schema) Listen(stopCh <-chan struct{}) error {
	var acceptor protocol.Acceptor
	switch schema.Protocol {
	case TCP:
		acceptor = protocol.NewTCPTCPAcceptor(schema.Addr)
	case WEBSOCKET:
		acceptor = protocol.NewWSAcceptor(schema.Addr)
	default:
		return fmt.Errorf("unsupported protocol: %v", schema.Protocol)
	}

	go runtime.WaitClose(stopCh, acceptor.Shutdown)
	return acceptor.Run()
}

func NewSchema(protocol Protocol, addr string) Schema {
	return Schema{
		Protocol: protocol,
		Addr:     addr,
	}
}
